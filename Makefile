DBG         ?= 0
#REGISTRY    ?= quay.io/openshift/
VERSION     ?= $(shell git describe --always --abbrev=7)
MUTABLE_TAG ?= latest
IMAGE        = $(REGISTRY)go-rag-bot
BUILD_IMAGE ?= docker.io/library/golang:1.23.7
GOLANGCI_LINT = go run ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint
ENVTEST_K8S_VERSION = 1.32.1

# Enable go modules and vendoring
# https://github.com/golang/go/wiki/Modules#how-to-install-and-activate-module-support
# https://github.com/golang/go/wiki/Modules#how-do-i-use-vendoring-with-modules-is-vendoring-going-away
GO111MODULE = on
export GO111MODULE
GOFLAGS ?= -mod=vendor
export GOFLAGS

ifeq ($(DBG),1)
GOGCFLAGS ?= -gcflags=all="-N -l"
endif

#
# Directories.
#
TOOLS_DIR=$(abspath ./tools)
BIN_DIR=bin
TOOLS_BIN_DIR := $(TOOLS_DIR)/$(BIN_DIR)

.PHONY: all
all: check build test

NO_DOCKER ?= 1

ifeq ($(shell command -v podman > /dev/null 2>&1 ; echo $$? ), 0)
	ENGINE=podman
else ifeq ($(shell command -v docker > /dev/null 2>&1 ; echo $$? ), 0)
	ENGINE=docker
else
	NO_DOCKER=1
endif

USE_DOCKER ?= 0
ifeq ($(USE_DOCKER), 1)
	ENGINE=docker
endif

ifeq ($(NO_DOCKER), 1)
  DOCKER_CMD =
  IMAGE_BUILD_CMD = imagebuilder
else
  DOCKER_CMD := $(ENGINE) run --env GO111MODULE=$(GO111MODULE) --env GOFLAGS=$(GOFLAGS) --rm -v "$(PWD)":/go/src/github.com/aguidirh/go-rag-chatbot:Z -w /go/src/github.com/aguidirh/go-rag-chatbot $(BUILD_IMAGE)
  # The command below is for building/testing with the actual image that Openshift uses. Uncomment/comment out to use instead of above command. CI registry pull secret is required to use this image.
  IMAGE_BUILD_CMD = $(ENGINE) build
endif

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: build
build: go-rag-bot ## Build binaries

.PHONY: go-rag-bot
go-rag-bot:
	$(DOCKER_CMD) ./hack/go-build.sh go-rag-bot

.PHONY: lint
lint: ## Run golangci-lint over the codebase.
	 $(call ensure-home, ${GOLANGCI_LINT} run ./pkg/... ./cmd/... --timeout=10m)

.PHONY: help
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## --------------------------------------
## Cleanup / Verification
## --------------------------------------

##@ clean:

.PHONY: clean
clean: ## Remove generated binaries, GitBook files, Helm charts, and Tilt build files
	$(MAKE) clean-bin

.PHONY: clean-bin
clean-bin: ## Remove all generated binaries
	rm -rf $(BIN_DIR)

define ensure-home
	@ export HOME=$${HOME:=/tmp/kubebuilder-testing}; \
	if [ $${HOME} == "/" ]; then \
	  export HOME=/tmp/kubebuilder-testing; \
	fi; \
	$(1)
endef