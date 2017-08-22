
SHELL := /bin/bash
VERSION = $(shell git describe --tags --always --dirty)
BUILDDATE = $(shell date +%s)
EXECUTABLE = greeter
PKG = github.com/thingful/${EXECUTABLE}
BUILDFLAGS = -a -installsuffix cgo -ldflags "-X ${PKG}/pkg/version.Version=${VERSION} -X ${PKG}/pkg/version.BuildDate=${BUILDDATE}"
BUILD_DIR = build

.PHONY: help
help: ## Show this help message
	@echo 'usage: make [target] ...'
	@echo
	@echo 'targets:'
	@echo
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s : | sed -e 's/^/  /')"

.PHONY: protoc
protoc: ## Compile proto definitions to generate client/server code
	protoc --proto_path=pkg/${EXECUTABLE}/ ./pkg/${EXECUTABLE}/*.proto --go_out=plugins=grpc:pkg/${EXECUTABLE}

.PHONY: build-internal
build-internal: ## Build our Go executable. Note this is designed to be run inside the container
	export CGO_ENABLED=0
	export GOOS=linux
	mkdir -p $(BUILD_DIR)
	go build ${BUILDFLAGS} -o ${BUILD_DIR}/${EXECUTABLE} ${PKG}/pkg/server

.PHONY: build
build: ## Package our app inside a container using docker-compose
	docker-compose build
