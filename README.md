## Docker Hub Image Repo Migration Script
This Golang script scans a source Docker Hub organization for all image repositories and their tags and migrates them to a destination paid Dockerhub organization. It takes command-line parameters for `srcOrg`, `dstOrg`, and `shouldMigrate` where shouldMigrate can be a comma-separated list of repositories that should be migrated. If no parameters are provided, it uses placeholder default values.

### How it works
The script uses the Docker Hub API to retrieve a list of all image repositories in the source organization and then checks each one against the `shouldMigrate` function to determine if it should be migrated. For each repository that should be migrated, it retrieves a list of all tags using the Docker Hub API and then uses the Docker CLI to pull each image with its tag from the source organization, tag it with the destination organization, and push it to the destination organization.

A progress bar is displayed using the `pb` package to show how many tags have been migrated for each repository.

### Building and running
To build this script, you need to have Golang installed on your machine. You can then use the following command in your terminal from within the directory containing this script:
```go
go build -o migration
```

This will produce an executable binary named `migration`. You can then run this binary with any desired command-line parameters like so:
```shell
./migration -srcOrg=source_org -dstOrg=destination_org -shouldMigrate=repo1,repo2
```
or run the docker image:

```shell
docker run -it --rm masterxavierfox/dockergedon -srcOrg=source_org -dstOrg=destination_org -shouldMigrate=repo1,repo2
```
Make sure you have Docker installed on your machine as well since this script uses the Docker CLI. Also make sure you are logged into Docker Hub on your machine. and the migration user hass access to the source and destination organization.

### Use makefile

You can then use the following commands in your terminal from within this directory:

- `make build` : Builds a binary for either Mac M1 or Linux depending on your machine architecture.
- `make docker-build`: Builds a Docker image using the Dockerfile in this directory.
- `make docker-push`: Pushes the built Docker image to a repository on Docker Hub.
- `make run-docker ARGS="-srcOrg=source_org -dstOrg=destination_org -shouldMigrate=repo1,repo2"`: Runs the built Docker image with some example command-line parameters. You can modify these parameters as needed.
- `make run-binary ARGS="-srcOrg=source_org -dstOrg=destination_org -shouldMigrate=repo1,repo2"`: Runs the built binary with some example command-line parameters. You can modify these parameters as needed.

#### Motivation
[Docker is sunsetting Free Team organizations [pdf]](https://web.docker.com/rs/790-SSB-375/images/privatereposfaq.pdf)