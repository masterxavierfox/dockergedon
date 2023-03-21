DOCKER_IMAGE = docker.io/masterxavierfox/dockergedon
BINARY_NAME = dockerdon

build:
ifeq ($(shell uname -m), arm64)
	GOOS=darwin GOARCH=arm64 go build -o ./builds/$(BINARY_NAME)
else
	GOOS=linux GOARCH=amd64 go build -o ./builds/$(BINARY_NAME)
endif

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-push:
	docker push $(DOCKER_IMAGE)

run-binary:
ifeq ($(shell uname -m), arm64)
	./builds/$(BINARY_NAME) $(ARGS)
else
	./builds/$(BINARY_NAME) $(ARGS)
endif

run-docker:
	docker run --rm $(DOCKER_IMAGE) ./builds/$(BINARY_NAME) $(ARGS)