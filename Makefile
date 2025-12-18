# Makefile for idiogo

.PHONY: help build test lint clean run docker-up docker-down migrate

# Variables
BINARY_NAME=idiogo
GO=go
GOTEST=$(GO) test
GOVET=$(GO) vet
GOFMT=gofmt
DOCKER_COMPOSE=docker compose -f deployments/compose.yml

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m

help: ## Display this help screen
	@echo "$(COLOR_BOLD)idiogo - Makefile commands:$(COLOR_RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'

build: ## Build the application
	@echo "$(COLOR_BOLD)Building idiogo...$(COLOR_RESET)"
	@mkdir -p bin
	@$(GO) build -o bin/serve cmd/serve/main.go
	@$(GO) build -o bin/cron cmd/cron/main.go
	@echo "$(COLOR_GREEN)✓ Build complete$(COLOR_RESET)"

build-all: ## Build for all platforms
	@echo "$(COLOR_BOLD)Building for all platforms...$(COLOR_RESET)"
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 $(GO) build -o bin/serve-linux-amd64 cmd/serve/main.go
	@GOOS=darwin GOARCH=amd64 $(GO) build -o bin/serve-darwin-amd64 cmd/serve/main.go
	@GOOS=darwin GOARCH=arm64 $(GO) build -o bin/serve-darwin-arm64 cmd/serve/main.go
	@GOOS=windows GOARCH=amd64 $(GO) build -o bin/serve-windows-amd64.exe cmd/serve/main.go
	@echo "$(COLOR_GREEN)✓ Multi-platform build complete$(COLOR_RESET)"

run: ## Run the application
	@echo "$(COLOR_BOLD)Starting idiogo...$(COLOR_RESET)"
	@$(GO) run cmd/serve/main.go

test: ## Run tests
	@echo "$(COLOR_BOLD)Running tests...$(COLOR_RESET)"
	@$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@echo "$(COLOR_GREEN)✓ Tests complete$(COLOR_RESET)"

test-coverage: test ## Run tests with coverage report
	@echo "$(COLOR_BOLD)Generating coverage report...$(COLOR_RESET)"
	@$(GO) tool cover -html=coverage.txt -o coverage.html
	@echo "$(COLOR_GREEN)✓ Coverage report generated: coverage.html$(COLOR_RESET)"

test-integration: ## Run integration tests
	@echo "$(COLOR_BOLD)Running integration tests...$(COLOR_RESET)"
	@$(GOTEST) -v -tags=integration ./...

bench: ## Run benchmarks
	@echo "$(COLOR_BOLD)Running benchmarks...$(COLOR_RESET)"
	@$(GOTEST) -bench=. -benchmem ./...

lint: ## Run linter
	@echo "$(COLOR_BOLD)Running linter...$(COLOR_RESET)"
	@golangci-lint run --timeout=5m
	@echo "$(COLOR_GREEN)✓ Linting complete$(COLOR_RESET)"

fmt: ## Format code
	@echo "$(COLOR_BOLD)Formatting code...$(COLOR_RESET)"
	@$(GOFMT) -s -w .
	@$(GO) fmt ./...
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

vet: ## Run go vet
	@echo "$(COLOR_BOLD)Running go vet...$(COLOR_RESET)"
	@$(GOVET) ./...
	@echo "$(COLOR_GREEN)✓ Vet complete$(COLOR_RESET)"

sec: ## Run security scan
	@echo "$(COLOR_BOLD)Running security scan...$(COLOR_RESET)"
	@gosec -exclude-generated ./...
	@echo "$(COLOR_GREEN)✓ Security scan complete$(COLOR_RESET)"

tidy: ## Tidy go modules
	@echo "$(COLOR_BOLD)Tidying go modules...$(COLOR_RESET)"
	@$(GO) mod tidy
	@echo "$(COLOR_GREEN)✓ Modules tidied$(COLOR_RESET)"

clean: ## Clean build artifacts
	@echo "$(COLOR_BOLD)Cleaning...$(COLOR_RESET)"
	@rm -rf bin/
	@rm -f coverage.txt coverage.html
	@$(GO) clean
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

docker-build: ## Build Docker image
	@echo "$(COLOR_BOLD)Building Docker image...$(COLOR_RESET)"
	@docker build -f cmd/serve/Dockerfile -t idiogo:latest .
	@echo "$(COLOR_GREEN)✓ Docker image built$(COLOR_RESET)"

docker-up: ## Start Docker Compose services
	@echo "$(COLOR_BOLD)Starting Docker services...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(COLOR_GREEN)✓ Docker services started$(COLOR_RESET)"

docker-down: ## Stop Docker Compose services
	@echo "$(COLOR_BOLD)Stopping Docker services...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) down
	@echo "$(COLOR_GREEN)✓ Docker services stopped$(COLOR_RESET)"

docker-logs: ## View Docker Compose logs
	@$(DOCKER_COMPOSE) logs -f

docker-restart: docker-down docker-up ## Restart Docker services

install-tools: ## Install development tools
	@echo "$(COLOR_BOLD)Installing development tools...$(COLOR_RESET)"
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	@$(GO) install github.com/cosmtrek/air@latest
	@echo "$(COLOR_GREEN)✓ Tools installed$(COLOR_RESET)"

dev: ## Run with hot reload (requires air)
	@echo "$(COLOR_BOLD)Starting development server with hot reload...$(COLOR_RESET)"
	@air

db-up: ## Start only database
	@echo "$(COLOR_BOLD)Starting database...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) up -d idiogo-pg
	@echo "$(COLOR_GREEN)✓ Database started$(COLOR_RESET)"

db-down: ## Stop database
	@echo "$(COLOR_BOLD)Stopping database...$(COLOR_RESET)"
	@$(DOCKER_COMPOSE) stop idiogo-pg
	@echo "$(COLOR_GREEN)✓ Database stopped$(COLOR_RESET)"

db-reset: ## Reset database (WARNING: destroys data)
	@echo "$(COLOR_YELLOW)⚠ WARNING: This will destroy all data!$(COLOR_RESET)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(COLOR_BOLD)Resetting database...$(COLOR_RESET)"; \
		$(DOCKER_COMPOSE) down -v; \
		$(DOCKER_COMPOSE) up -d idiogo-pg; \
		echo "$(COLOR_GREEN)✓ Database reset complete$(COLOR_RESET)"; \
	fi

db-shell: ## Connect to database shell
	@docker exec -it idiogo-db psql -U idiogo -d idiogo

deps: ## Download dependencies
	@echo "$(COLOR_BOLD)Downloading dependencies...$(COLOR_RESET)"
	@$(GO) mod download
	@echo "$(COLOR_GREEN)✓ Dependencies downloaded$(COLOR_RESET)"

verify: fmt vet lint test ## Run all verification steps

release: ## Create a new release (requires goreleaser)
	@echo "$(COLOR_BOLD)Creating release...$(COLOR_RESET)"
	@goreleaser release --clean
	@echo "$(COLOR_GREEN)✓ Release complete$(COLOR_RESET)"

snapshot: ## Create a snapshot release
	@echo "$(COLOR_BOLD)Creating snapshot...$(COLOR_RESET)"
	@goreleaser release --snapshot --clean
	@echo "$(COLOR_GREEN)✓ Snapshot created$(COLOR_RESET)"

.DEFAULT_GOAL := help
