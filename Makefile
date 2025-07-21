# Makefile for Item PDP Service

# Variables
BINARY_NAME=item-pdp-service
DOCKER_IMAGE=item-pdp-service
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Go commands
GO_BUILD=go build
GO_TEST=go test
GO_CLEAN=go clean
GO_MOD=go mod

# Docker commands
DOCKER_BUILD=docker build
DOCKER_RUN=docker run
DOCKER_COMPOSE=docker compose

.PHONY: help build test test-unit test-integration test-coverage test-verbose clean run docker-build docker-run docker-up docker-down mod-tidy mod-download lint format

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	$(GO_BUILD) -o bin/$(BINARY_NAME) cmd/api/main.go

test: ## Run all tests
	$(GO_TEST) -v ./...

test-unit: ## Run only unit tests (exclude integration tests)
	$(GO_TEST) -v ./internal/... -short

test-integration: ## Run only integration tests
	$(GO_TEST) -v ./test/integration/...

test-coverage: ## Run tests with coverage report
	$(GO_TEST) -v -race -covermode=atomic -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@go tool cover -func=$(COVERAGE_FILE) | tail -1

test-coverage-func: ## Show coverage by function
	$(GO_TEST) -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)

test-coverage-target: ## Check if coverage meets 80% target
	$(GO_TEST) -coverprofile=$(COVERAGE_FILE) ./...
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_FILE) | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$COVERAGE >= 80" | bc -l) -eq 1 ]; then \
		echo "✅ Coverage target met: $$COVERAGE%"; \
	else \
		echo "❌ Coverage below target: $$COVERAGE% (target: 80%)"; \
		exit 1; \
	fi

test-verbose: ## Run tests with verbose output
	$(GO_TEST) -v -race ./...

test-watch: ## Watch for changes and run tests
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air -c .air.toml

clean: ## Clean build artifacts and cache
	$(GO_CLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	go clean -testcache

run: ## Run the application locally
	go run cmd/api/main.go

run-dev: ## Run the application in development mode with auto-reload
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

docker-build: ## Build Docker image
	$(DOCKER_BUILD) -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	$(DOCKER_RUN) -p 8080:8080 $(DOCKER_IMAGE)

docker-up: ## Start services with docker-compose
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop services with docker-compose
	$(DOCKER_COMPOSE) down

docker-logs: ## View docker-compose logs
	$(DOCKER_COMPOSE) logs -f

mod-tidy: ## Tidy go modules
	$(GO_MOD) tidy

mod-download: ## Download go modules
	$(GO_MOD) download

mod-vendor: ## Vendor go modules
	$(GO_MOD) vendor

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

format: ## Format code
	go fmt ./...
	@which goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	goimports -w .

security: ## Run security scan
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	gosec ./...

deps: ## Install development dependencies
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

migrate-up: ## Run database migrations up
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/item_pdp_db?sslmode=disable" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/item_pdp_db?sslmode=disable" down

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	migrate create -ext sql -dir migrations $(NAME)

db-setup: ## Setup local database for testing
	$(DOCKER_COMPOSE) up -d postgres
	sleep 5
	make migrate-up

setup: deps mod-download ## Setup development environment
	@echo "Development environment setup complete!"

ci: mod-tidy test-coverage-target lint ## Run CI pipeline
	@echo "CI pipeline completed successfully!"

# Default target
.DEFAULT_GOAL := help 