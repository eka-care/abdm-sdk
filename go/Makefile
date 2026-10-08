# Eka Care SDK for Go

# Go parameters
GOCMD=go
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint
GOSEC=gosec

# Directories
TEST_DIR=./...
COVERAGE_DIR=coverage

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo "Eka Care SDK for Go - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Clean build artifacts and cache
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(COVERAGE_DIR)
	$(GOCMD) clean -cache -testcache -modcache

.PHONY: test
test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v $(TEST_DIR)

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -coverprofile=$(COVERAGE_DIR)/coverage.out $(TEST_DIR)
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	$(GOTEST) -race $(TEST_DIR)

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	$(GOLINT) run

.PHONY: lint-fix
lint-fix: ## Run linter with auto-fix
	@echo "Running linter with auto-fix..."
	$(GOLINT) run --fix

.PHONY: security
security: ## Run security scan
	@echo "Running security scan..."
	$(GOSEC) ./...

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

.PHONY: mod-tidy
mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	$(GOMOD) tidy

.PHONY: mod-verify
mod-verify: ## Verify go modules
	@echo "Verifying go modules..."
	$(GOMOD) verify

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

.PHONY: example
example: ## Run quickstart example
	@echo "Running quickstart example..."
	cd examples/quickstart && $(GOCMD) run main.go

.PHONY: check
check: fmt vet lint security ## Run all code quality checks
	@echo "All checks completed successfully!"

.PHONY: test-all
test-all: test test-race test-coverage ## Run all tests
	@echo "All tests completed successfully!"

.PHONY: ci
ci: mod-tidy check test-all ## Run full CI pipeline
	@echo "CI pipeline completed successfully!"

.PHONY: dev-setup
dev-setup: install-tools mod-tidy ## Setup development environment
	@echo "Development environment setup completed!"