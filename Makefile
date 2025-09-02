.PHONY: $(shell sed -n -e '/^$$/ { n ; /^[^ .\#][^ ]*:/ { s/:.*$$// ; p ; } ; }' $(MAKEFILE_LIST))
.DEFAULT_GOAL := help

help: ## Display available make commands
	@echo $(BLUE)"Proxmox Homelab Management"$(RESET)
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/ghprs cmd/main.go

run: ## Run the application
	go run cmd/main.go

test: ## Run tests
	go test ./...

test-v: ## Run tests with verbose output
	go test -v ./...

test-race: ## Run tests with race detection
	go test -race ./...

test-coverage: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

test-coverage-summary: ## Run tests with coverage and show summary
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-bench: ## Run benchmark tests
	go test -bench=. ./...

test-bench-mem: ## Run benchmark tests with memory allocation info
	go test -bench=. -benchmem ./...

test-api: ## Run tests for specific packages
	go test -v ./internal/api/...

test-models:
	go test -v ./internal/models/...

test-ui:
	go test -v ./internal/ui/...

test-config:
	go test -v ./pkg/config/...

test-timeout: ## Run tests with timeout
	go test -timeout=30s ./...

test-badge: ## Run tests and generate coverage badge
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//' > coverage.txt

test-all: test-v test-race test-coverage-summary test-bench ## Run all tests (comprehensive)
	@echo "All tests completed!"

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out
	rm -f coverage.txt

install-deps: ## Install dependencies
	go mod tidy
	go mod download

lint: ## Run linter
	golangci-lint run

install-hooks: ## Install pre-commit hooks
	pre-commit install

pre-commit: ## Run pre-commit on all files
	pre-commit run --all-files

setup: install-deps install-hooks ## Development setup
	@echo "Development environment setup complete!"
