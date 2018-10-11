# Docker image Go builder
GO_BUILDER = go_builder:1.10-alpine3.7
# Local source path
ifeq ($(OS),Windows_NT)
	SRC_LOCAL = "%CD%"
	REMOVE = del
else
	SRC_LOCAL = $(shell pwd)
	REMOVE = rm
endif
# image Go builder source path
ifeq ($(@shell uname),MINGW64_NT-10.0)
# double // to avoid path translate in windows git bash
	SRC_IMAGE = //go/app
else
	SRC_IMAGE = /go/app
endif
# Application config
APP_NAME = api-klit
# OpenShift config
OC_PROJECT = frlm-web-team-quality-01-dev
REGISTRY = registry.infra.op.acp.adeo.com
OC_REPO = $(REGISTRY)/$(OC_PROJECT)/api-klit:develop
# Develop config
DEV_IMAGE = klit_api
DEV_CONTAINER = klit_api_run

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
	
image_go_build: ## Create the go build image containing so git
ifeq ($(shell docker inspect $(GO_BUILDER)),[])
	docker build -t $(GO_BUILDER) -f ./docker/Dockerfile.buildImage-go1.10-alpine3.7 ./docker
endif

build: image_go_build ## Build the binary
	docker run --rm -v $(SRC_LOCAL):$(SRC_IMAGE) -w $(SRC_IMAGE) $(GO_BUILDER) /bin/sh -c "go get -d -v ./... && go build -v -ldflags \"-X main.buildDate=`date -Iseconds` -X main.gitHash=`git rev-parse HEAD`\" -o $(APP_NAME)"

local_deploy: build ## Build the local image for developpment
# Clean container/image
ifneq ($(shell docker inspect $(DEV_CONTAINER)),[])
	docker stop $(DEV_CONTAINER)
	docker rm $(DEV_CONTAINER)
endif
ifneq ($(shell docker inspect $(DEV_IMAGE)),[])
	docker rmi $(DEV_IMAGE)
endif
# Create image
	docker build -t $(DEV_IMAGE) -f ./docker/Dockerfile.openshift .
# Create container, copy database config file, and start
	docker create -p 7000:8080 -p 7001:8081 --env APP_SPRING_CONFIG_URI='http://config-server.qa.apps.op.acp.adeo.com' --env APP_PROFILES=local_sd --name $(DEV_CONTAINER) -ti $(DEV_IMAGE)
	docker cp ./Configuration.json $(DEV_CONTAINER):/app
	docker start $(DEV_CONTAINER)
	
oc_deploy: build ## Build the image, and push on OpenShift for deploying
# Build the image
	docker build -t $(OC_REPO)  -f ./docker/Dockerfile.openshift .
# Login on OpenShift
	oc login https://openshift.op.acp.adeo.com
# Switch on project
	oc project $(OC_PROJECT)
# Push l'image sur la registry
	oc whoami -t | docker login --password-stdin -u unused $(REGISTRY)
	docker push $(OC_REPO)
	
clean: ## Clean build files
	$(REMOVE) $(APP_NAME)
