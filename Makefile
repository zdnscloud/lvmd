REGISTRY_NAME = zdnscloud
IMAGE_Name = lvmd
BRANCH=`git branch | sed -n '/\* /s///p'`
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`

.PHONY: all container

all: build

build:
	CGO_ENABLED=0 GOOS=linux go build

image: 
	docker build -t $(REGISTRY_NAME)/$(IMAGE_Name):$(BRANCH) ./ --build-arg version=${VERSION} --build-arg buildtime=${BUILD} --no-cache
	docker image prune -f

docker:image
	docker push $(REGISTRY_NAME)/$(IMAGE_Name):$(BRANCH)
