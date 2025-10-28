# Copyright (c) Unikraft GmbH
# SPDX-License-Identifier: MPL-2.0

.DEFAULT_GOAL := help

# Binary name
BINARY := terraform-provider-ukc
GOBIN := $(shell go env GOBIN)
ifeq ($(GOBIN),)
GOBIN := $(shell go env GOPATH)/bin
endif

##@ Development

.PHONY: build
build: ## Build the provider binary
	@go build -v -o $(BINARY) .

.PHONY: install
install: build ## Build and install provider to $GOBIN for local development
	@mkdir -p $(GOBIN)
	@cp $(BINARY) $(GOBIN)/
	@echo "Provider installed to $(GOBIN)/$(BINARY)"
	@echo "Make sure your ~/.terraformrc has dev_overrides configured"

.PHONY: fmt
fmt: ## Format Go code with gofumpt
	@command -v gofumpt >/dev/null 2>&1 || { echo "Installing gofumpt..."; go install mvdan.cc/gofumpt@latest; }
	@gofumpt -l -w .

.PHONY: tidy
tidy: ## Run go mod tidy
	@go mod tidy

##@ Testing

.PHONY: test
test: ## Run unit tests
	@go test -v -cover ./internal/provider/...

##@ Quality

.PHONY: lint
lint: ## Run golangci-lint
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not found."; exit 1; }
	@golangci-lint run

##@ Cleanup

.PHONY: clean
clean: ## Remove build artifacts
	@rm -f $(BINARY)
	@echo "Build artifacts cleaned"

##@ Help

.PHONY: help
help: ## Show this help menu
	@awk 'BEGIN { \
		FS = ":.*##"; \
		printf "Terraform Provider for Unikraft Cloud\n\n"; \
		printf "\033[1mUSAGE\033[0m\n"; \
		printf "  make [VAR=... [VAR=...]] \033[36mTARGET\033[0m\n"; \
	} \
	/^##@/ { \
		printf "\n\033[1m%s\033[0m\n", substr($$0, 5); \
	} \
	/^[a-zA-Z0-9_-]+:.*?##/ { \
		printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2; \
	}' $(MAKEFILE_LIST)
	@echo ""
