SHELL := /bin/bash
.PHONY: help build test deps generate clean

# use the latest git tag as release-version
GIT_TAG_NAME=$(shell git tag --sort=-creatordate|head -n 1)
BUILD_TIME=$(shell date "+%Y-%m-%dT%H:%M:%S%z")
ifeq ($(GIT_TAG_NAME),)
GIT_TAG_NAME := "unknown"
endif
LD_FLAGS=-ldflags="-s -w -X github.com/lindb/lindb/config.Version=$(GIT_TAG_NAME) -X github.com/lindb/lindb/config.BuildTime=$(BUILD_TIME)"

# Ref: https://gist.github.com/prwhite/8168133
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

run: ## run local standalone cluster for demo/debug
	go run github.com/lindb/lindb/cmd/lind standalone run --pprof --doc

build-frontend: ## build frontend
	cd web/ && make web_build

GOARCH = amd64
build: clean-build build-frontend build-lind ## Build executable files.

build-all: clean-frontend-build build-frontend clean-build build-lind ## Build executable files with front-end files inside.

build-lind: ## build lindb binary
	env GOOS=darwin GOARCH=$(GOARCH) go build -o 'bin/lind-darwin' $(LD_FLAGS) ./cmd/lind
	env GOOS=linux GOARCH=$(GOARCH) go build -o 'bin/lind-linux' $(LD_FLAGS) ./cmd/lind
	env GOOS=windows GOARCH=$(GOARCH) go build -o 'bin/lind-windows.exe' $(LD_FLAGS) ./cmd/lind
	env GOOS=darwin GOARCH=$(GOARCH) go build -o 'bin/lindcli-darwin' $(LD_FLAGS) ./cmd/cli
	env GOOS=linux GOARCH=$(GOARCH) go build -o 'bin/lindcli-linux' $(LD_FLAGS) ./cmd/cli
	env GOOS=windows GOARCH=$(GOARCH) go build -o 'bin/lindcli-windows.exe' $(LD_FLAGS) ./cmd/cli

deploy: build-frontend ## deploy release packages
	sh deploy.sh

.PHONY: docker-build
docker-build: ## build docker image
	docker build -t lindata/lindb:$(GIT_TAG_NAME) --build-arg LD_FLAGS=${LD_FLAGS} .

.PHONY: docker-push
docker-push: ## push docker image
	docker push lindata/lindb:$(GIT_TAG_NAME)

gomock: ## go generate mock file.
	if [[ ! -x $$(command -v mockgen) ]]; then \
		go install go.uber.org/mock/mockgen@latest; \
	fi;
	find . -type f -name '*_mock.go' -exec rm -f {} +
	go generate -x ./...

header: ## check and add license header.
	sh scripts/addlicense.sh

import: ## opt go imports format.
	sh scripts/imports.sh

format: ## go format 
	go fmt ./...

lint: ## run lint
ifeq (, $(shell which golangci-lint))
	# binary will be $(go env GOPATH)/bin/golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2
else
	echo "Found golangci-lint"
endif
	golangci-lint run ./...

api-doc: ## generate api document
	go install "github.com/swaggo/swag/cmd/swag@v1.8.7"
	swag init -g pkg/http/doc.go

test-without-lint: ## Run test without lint
	go install "github.com/rakyll/gotest@v0.0.6"
	GIN_MODE=release
	LOG_LEVEL=fatal ## disable log for test
	gotest -v -race -coverprofile=coverage_tmp.out -covermode=atomic ./...
	cat coverage_tmp.out |grep -v "_mock.go" > coverage.out
	go tool cover -html=coverage.out -o coverage.html

test: header lint test-without-lint ## Run test cases.

e2e-test:
	go install "github.com/rakyll/gotest@v0.0.6"
	GIN_MODE=release
	LOG_LEVEL=fatal ## disable log for test
	gotest -v --tags=integration -covermode=atomic ./e2e/...

e2e: header e2e-test

deps:  ## Update vendor.
	go mod verify
	go mod tidy -v

generate:  ## generate pb file.
    # brew install flatbuffers
	sh ./proto/generate.sh

gen-sql-grammar: ## generate lin query language gen-sql-grammar
	# need install antrl4-tools
	# https://github.com/antlr/antlr4/blob/master/doc/getting-started.md
	antlr4 -Dlanguage=Go -listener -visitor -package grammar ./sql/grammar/SQLLexer.g4
	antlr4 -Dlanguage=Go -listener -visitor -package grammar ./sql/grammar/SQLParser.g4

key-words: ## print all key words for lin query language
	go run github.com/lindb/lindb/cmd/tools keywords 

clean-mock: ## remove all mock files
	find ./ -name "*_mock.go" | xargs rm

clean-build:
	rm -rf bin/*

clean-frontend-build:
	cd web/ && make web_clean

clean-tmp: ## clean up tmp and test out files
	find . -type f -name '*.out' -exec rm -f {} +
	find . -type f -name '.DS_Store' -exec rm -f {} +
	find . -type f -name '*.test' -exec rm -f {} +
	find . -type f -name '*.prof' -exec rm -f {} +
	find . -type s -name 'localhost:*' -exec rm -f {} +
	find . -type s -name '127.0.0.1:*' -exec rm -f {} +
	rm -rf data

clean: clean-mock clean-tmp clean-build clean-frontend-build ## Clean up useless files.


