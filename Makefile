.PHONY: help build test deps generate clean

# use the latest git tag as release-version
GIT_TAG_NAME=$(shell git tag --sort=-creatordate|head -n 1)
BUILD_TIME=$(shell date "+%Y-%m-%dT%H:%M:%S%z")
LD_FLAGS=-ldflags="-X github.com/lindb/lindb/cmd/lind.version=$(GIT_TAG_NAME) -X github.com/lindb/lindb/cmd/lind.buildTime=$(BUILD_TIME)"

# Ref: https://gist.github.com/prwhite/8168133
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build-frontend: clean-frontend-build
	cd web/ && make web_build

GOARCH = amd64
build: clean-build build-lind ## Build executable files.

build-all: clean-frontend-build build-frontend clean-build build-lind ## Build executable files with front-end files inside.

build-lind: ## build lindb binary
	env GOOS=darwin GOARCH=$(GOARCH) go build -o 'bin/lind-darwin' $(LD_FLAGS) ./cmd/
	env GOOS=linux GOARCH=$(GOARCH) go build -o 'bin/lind-linux' $(LD_FLAGS) ./cmd/


GOLANGCI_LINT_VERSION ?= "v1.28.3"

GOMOCK_VERSION = "v1.5.0"

gomock: ## go generate mock file.
	go get "github.com/golang/mock/mockgen@$(GOMOCK_VERSION)"
	go install "github.com/golang/mock/mockgen"
	go list ./... |grep -v '/gomock' | xargs go generate -v

header: ## check and add license header.
	sh license.sh

lint: ## run lint
	if [ ! -e ./bin/golangci-lint ]; then \
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s $(GOLANGCI_LINT_VERSION); \
	fi
	./bin/golangci-lint run

test-without-lint: ## Run test without lint
	go get -u github.com/rakyll/gotest
	export LOG_LEVEL="fatal" ## disable log for test
	gotest -v -race -coverprofile=coverage.out -covermode=atomic ./...

test: header lint test-without-lint ## Run test cases.

deps:  ## Update vendor.
	go mod verify
	go mod tidy -v

generate:  ## generate pb/tmpl file.
	# go get github.com/benbjohnson/tmpl
	# go install github.com/benbjohnson/tmpl
	sh ./rpc/pb/generate_pb.sh
	cd tsdb/template && sh generate_tmpl.sh

clean-mock: ## remove all mock files
	find ./ -name "*_mock.go" | xargs rm

clean-build:
	rm -f bin/lind-darwin
	rm -f bin/lind-linux

clean-frontend-build:
	cd web/ && make web_clean

clean-tmp: ## clean up tmp and test out files
	find . -type f -name '*.out' -exec rm -f {} +
	find . -type f -name '.DS_Store' -exec rm -f {} +
	find . -type f -name '*.test' -exec rm -f {} +
	find . -type f -name '*.prof' -exec rm -f {} +
	find . -type s -name 'localhost:*' -exec rm -f {} +
	find . -type s -name '127.0.0.1:*' -exec rm -f {} +

clean: clean-mock clean-tmp clean-build clean-frontend-build ## Clean up useless files.
