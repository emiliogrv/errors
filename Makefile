SHELL := /bin/bash

ifeq ($(OS),Windows_NT)
    # Windows absolute path, compatible with Go
	GOBIN := $(shell cygpath -w $(shell pwd))/bin
else
    # Unix-like absolute path
	GOBIN := $(shell go env PWD)/bin
endif

.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@grep -Eh '^[A-Za-z0-9_.-]+:.*##' $(firstword $(MAKEFILE_LIST)) \
		| awk 'BEGIN {FS = ":.*##"} {printf "  %-20s %s\n", $$1, $$2}' \
		| sort

.PHONY: setup
setup: ## Setup environment
	@echo "Setting up environment"
	@pre-commit install
	@go mod tidy
	@$(MAKE) install-tools

.PHONY: generate
generate: install-tools ## Run generator
	@echo "Running generator"
	@"$(GOBIN)/errors_generator" -with-gen-header=false -output-dir pkg/core
	@"$(GOBIN)/errors_generator" -with-gen-header=false -output-dir pkg/logrus -formats logrus
	@"$(GOBIN)/errors_generator" -with-gen-header=false -output-dir pkg/slog -formats slog
	@"$(GOBIN)/errors_generator" -with-gen-header=false -output-dir pkg/zap -formats zap
	@"$(GOBIN)/errors_generator" -with-gen-header=false -output-dir pkg/zerolog -formats zerolog
	@"$(GOBIN)/errors_generator" -test-gen strict -with-gen-header=false -output-dir pkg/full -formats all

.PHONY: lint
lint: install-tools ## Run linter
	@echo "Running linter"
	@"$(GOBIN)/golangci-lint" run

.PHONY: lint-fix
lint-fix: install-tools ## Run linter fix
	@echo "Running linter fix"
	@"$(GOBIN)/golangci-lint" run --fix

.PHONY: test
test: ## Run tests
	@echo "Running tests"
	@go test -race -count=1 ./...

.PHONY: install-tools
install-tools: ## Install tools
	@mkdir -p "$(GOBIN)"
	@if [ ! -f "$(GOBIN)/golangci-lint" ]; then \
		echo "Installing golangci-lint"; \
		GOBIN="$(GOBIN)" go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0; \
	fi
	@go build -C cmd/errors_generator -o ../../bin/errors_generator
