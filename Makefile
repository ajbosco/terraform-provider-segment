.EXPORT_ALL_VARIABLES:
NAME := terraform-provider-segment
BUILD_DIR := $(shell pwd)/build
TARGET := ${BUILD_DIR}/${NAME}
LDFLAGS ?=

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: build-linux ## run go build for linux

.PHONY: build-native
build-native: ## run go build for current OS
	@go build --mod=vendor -ldflags "$(LDFLAGS)" -o "${TARGET}"

.PHONY: build-linux
build-linux:
	@GOOS=linux GOARCH=amd64 go build --mod=vendor -ldflags "$(LDFLAGS)" -o "${TARGET}"

.PHONY: fmt
fmt: ## Verifies all files have been `gofmt`ed.
	@gofmt -s -l . | grep -v vendor | tee /dev/stderr

.PHONY: lint
lint: ## Verifies `golint` passes.
	@golint ./... | grep -v vendor | tee /dev/stderr

.PHONY: test
test: ## Runs the go tests.
	@go test -cover -race $(shell go list ./... | grep -v vendor)

.PHONY: vet
vet: ## Verifies `go vet` passes.
	@go vet $(shell go list ./... | grep -v vendor) | tee /dev/stderr