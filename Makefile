VERSION := $(shell git describe --tags --dirty)
BUILD := $(shell git rev-parse --short HEAD)
GOFILES := $(shell find . -not -path './vendor*' -type f -name '*.go')
LDFLAGS := -ldflags "-X=go.smartmachine.io/cumulus/cmd.Version=$(VERSION) -X=go.smartmachine.io/cumulus/cmd.Build=$(BUILD)"
TEST_STAMP := .test.stamp

.phony: all
all: dep test build ## Generate and build everything

.phony: test
test: $(TEST_STAMP) ## Run unit tests

.phony: dep
dep: ## Make sure all dependencies are up to date
	@go mod tidy
	@go mod vendor

$(TEST_STAMP): $(GOFILES)
	$(info Running unit tests)
	@go test ./pkg/...
	@touch $@

cumulus: $(GOFILES)
	$(info Compiling project)
	@go build -v $(LDFLAGS)

.phony: build
build: cumulus ## Build all binary artifacts

.phony: clean
clean: ## Clean all build artifacts
	$(info Cleaning all build artifacts)
	@rm -rf cumulus .test.stamp
	@go clean

.phony: veryclean
veryclean: clean ## Clean all caches and generated objects
	@go clean -cache -testcache -modcache

.phony: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
