APP ?= explorer
VERSION ?= $(strip $(shell cat VERSION))
GOOS ?= linux
SRC = ./

COMMIT = $(shell git rev-parse --short HEAD)
BRANCH = $(strip $(shell git rev-parse --abbrev-ref HEAD))
CHANGES = $(shell git rev-list --count ${COMMIT})
BUILDED ?= $(shell date -u '+%Y-%m-%dT%H:%M:%S')
BUILD_FLAGS = "-X main.Version=$(VERSION) -X main.GitCommit=$(COMMIT) -X main.BuildedDate=$(BUILDED)"
BUILD_TAGS?=minter-explorer-api
DOCKER_TAG = latest
SERVER ?= explorer.minter.network
PACKAGES=$(shell go list ./... | grep -v '/vendor/')

GOTOOLS = \
    github.com/golang/dep/cmd/dep

check: check_tools ensure_deps

all: check test build

### Tools & dependencies ####
check_tools:
	@# https://stackoverflow.com/a/25668869
	@echo "Found tools: $(foreach tool,$(notdir $(GOTOOLS)),\
        $(if $(shell which $(tool)),$(tool),$(error "No $(tool) in PATH")))"

get_tools:
	@echo "--> Installing tools"
	./get_tools.sh

#Run this from CI
get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep"
	@dep ensure -vendor-only

#Run this locally.
ensure_deps:
	@rm -rf vendor/
	@echo "--> Running dep"
	@dep ensure

### Build ###################
build: clean
	GOOS=${GOOS} go build -ldflags $(BUILD_FLAGS) -o ./builds/$(APP)

install:
	GOOS=${GOOS} go install -ldflags $(BUILD_FLAGS)

clean:
	@rm -f $(BINARY)

### Test ####################
test:
	@echo "--> Running tests"
	go test -v ${SRC}

fmt:
	@go fmt ./...

.PHONY: check check_tools get_vendor_deps ensure_deps build clean fmt test
