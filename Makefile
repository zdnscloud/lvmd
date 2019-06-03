REGISTRY_NAME = zdnscloud
IMAGE_Name = lvmd
IMAGE_VERSION = v0.94

.PHONY: all container

all: container

container: 
	docker build -t $(REGISTRY_NAME)/$(IMAGE_Name):$(IMAGE_VERSION) ./ --no-cache
