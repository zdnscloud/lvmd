REGISTRY_NAME = zdnscloud
IMAGE_NAME = lvmd
BRANCH=`git branch | sed -n '/\* /s///p'`
VERSION=`git describe --tags`
IMAGE_VERSION= latest
BUILD=`date +%FT%T%z`

.PHONY: all container

all: build

build:
	CGO_ENABLED=0 GOOS=linux go build

image: 
	docker build -t $(REGISTRY_NAME)/$(IMAGE_NAME):$(BRANCH) ./ --build-arg version=${VERSION} --build-arg buildtime=${BUILD} --no-cache
	docker image prune -f

docker:
	docker build -t $(REGISTRY_NAME)/$(IMAGE_NAME):$(IMAGE_VERSION) ./ --build-arg version=${VERSION} --build-arg buildtime=${BUILD} --no-cache
	docker image prune -f
	docker push $(REGISTRY_NAME)/$(IMAGE_NAME):$(IMAGE_VERSION)
