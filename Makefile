.PHONY: help build test deps pb clean

# use the latest git tag as release-version
GIT_TAG_NAME=$(shell git tag --sort=-creatordate|head -n 1)
BUILD_TIME=$(shell date "+%Y-%m-%dT%H:%M:%S%z")
LD_FLAGS=-ldflags="-X github.com/lindb/lindb/cmd/lind.version=$(GIT_TAG_NAME) -X github.com/lindb/lindb/cmd/lind.buildTime=$(BUILD_TIME)"

# Ref: https://gist.github.com/prwhite/8168133
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build-frontend:  clean-build
	cd web/ && make web_build

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
build: clean-build ## Build executable files. (Args: GOOS=$(go env GOOS) GOARCH=$(go env GOARCH))
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o 'bin/lind' $(LD_FLAGS) ./cmd/

build-linux: clean-build ## Build executable files. (Args: GOOS=$(go env GOOS) GOARCH=$(go env GOARCH))
	env GOOS=linux GOARCH=amd64 go build -o 'bin/lind' $(LD_FLAGS) ./cmd/

build-all: build-frontend build  ## Build executable files with front-end files inside.

GOLANGCI_LINT_VERSION ?= "v1.18.0"

pre-test: ## go generate mock file.
	go install "./ci/mockgen"
	go list ./... | grep -v '/vendor/' | xargs go generate
	if [[ "$$(uname)" == "Darwin" ]]; then \
		find . -path vendor -prune -o -type f \( -name '*_mock.go' -o -name '*_mock.pb.go' \) -exec \
		sed -i '' -e 's#x\.##g; s#x "\."##g' {} +; \
	else \
		find . -path vendor -prune -o -type f \( -name '*_mock.go' -o -name '*_mock.pb.go' \) -exec \
		sed -i 's#x\.##g; s#x "\."##g' {} +; \
	fi

	if [ ! -e ./bin/golangci-lint ]; then \
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s $(GOLANGCI_LINT_VERSION); \
	fi
	./bin/golangci-lint run

test:  pre-test ## Run test cases. (Args: GOLANGCI_LINT_VERSION=latest)
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

deps:  ## Update vendor.
	go mod verify
	go mod tidy -v
	rm -rf vendor
	go mod vendor -v

pb:  ## generate pb file.
	./ci/generate_pb.sh

clean-build:
	rm -f bin/lind
	cd web/ && make web_clean

clean-tmp: ## clean up tmp and test out files
	find . -type f -name '*.out' -exec rm -f {} +
	find . -type f -name '.DS_Store' -exec rm -f {} +
	find . -type f -name '*.test' -exec rm -f {} +
	find . -type f -name '*.prof' -exec rm -f {} +
	find . -type s -name 'localhost:*' -exec rm -f {} +
	find . -type s -name '127.0.0.1:*' -exec rm -f {} +

clean:  ## Clean up useless files.
	$(clean-build)
	$(clean-tmp)
