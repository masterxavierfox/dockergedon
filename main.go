package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

type Repositories struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

type Tags struct {
	Results []struct {
		Name string `json:"name"`
	} `json:"results"`
}

var srcOrg string
var dstOrg string
var shouldMigrateRepos []string

func init() {
	flag.StringVar(&srcOrg, "srcOrg", "cellulantops", "The source Docker Hub organization")
	flag.StringVar(&dstOrg, "dstOrg", "cellulant", "The destination Docker Hub organization")
	var shouldMigrateReposStr string
	flag.StringVar(&shouldMigrateReposStr, "shouldMigrate", "", "A comma-separated list of repositories that should be migrated")
	flag.Parse()

	if shouldMigrateReposStr != "" {
		shouldMigrateRepos = strings.Split(shouldMigrateReposStr, ",")
	}
}

func shouldMigrate(repoName string) bool {
	if len(shouldMigrateRepos) == 0 {
		return true
	}

	for _, r := range shouldMigrateRepos {
		if r == repoName {
			return true
		}
	}

	return false
}

func main() {

	repoURL := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/?page_size=100", srcOrg)
	resp, err := http.Get(repoURL)
	if err != nil {
		log.Fatalf("Failed to get repositories: %v\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v\n", err)
	}

	var repos Repositories
	if err := json.Unmarshal(body, &repos); err != nil {
		log.Fatalf("Failed	to	unmarshal	response	body:%v\n", err)
	}

	for _, repo := range repos.Results {
		if !shouldMigrate(repo.Name) {
			continue
		}

		tagURL := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/%s/tags/?page_size=100", srcOrg, repo.Name)
		resp, err := http.Get(tagURL)
		if err != nil {
			log.Fatalf("Failed	to	get	tags	for	repository %s:%v\n", repo.Name, err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed	to	read	response	body:%v\n", err)
		}

		var tags Tags
		if err := json.Unmarshal(body, &tags); err != nil {
			log.Fatalf("Failed to unmarshal response body: %v\n", err)
		}

		bar := pb.StartNew(len(tags.Results))

		for _, tag := range tags.Results {
			srcImageWithTag := fmt.Sprintf("%s/%s:%s", srcOrg, repo.Name, tag.Name)
			dstImageWithTag := fmt.Sprintf("%s/%s:%s", dstOrg, repo.Name, tag.Name)

			cmd := exec.Command("docker", "pull", srcImageWithTag)
			if err := cmd.Run(); err != nil {
				log.Fatalf("Failed to pull %s: %v\n", srcImageWithTag, err)
			}

			cmd = exec.Command("docker", "tag", srcImageWithTag, dstImageWithTag)
			if err := cmd.Run(); err != nil {
				log.Fatalf("Failed to tag %s: %v\n", dstImageWithTag, err)
			}

			cmd = exec.Command("docker", "push", dstImageWithTag)
			if err := cmd.Run(); err != nil {
				log.Fatalf("Failed to push %s: %v\n", dstImageWithTag, err)
			}

			fmt.Printf("Migrated %s	to  %s successfully.\n", srcImageWithTag, dstImageWithTag)

			bar.Increment()
		}

		bar.Finish()
	}
}
