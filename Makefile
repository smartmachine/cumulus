VERSION := $(shell git describe --tags --dirty)
BUILD := $(shell git rev-parse --short HEAD)
GOFILES := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
TEST_STAMP := .test.stamp

.phony: all
all: check ensure api test build ## Generate and build everything

.phony: test
test: $(TEST_STAMP) ## Run unit tests

.phony: dep
dep: ## Make sure all dependencies are up to date
	@go mod vendor -v

$(TEST_STAMP): $(GOFILES)
	$(info Running unit tests)
	@go test ./pkg/...
	@touch $@

compile: $(GOFILES)
	$(info Compiling project)
	@go build -v $(LDFLAGS)

.phony: build
build: dep compile ## Build all binary artifacts

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
