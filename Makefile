.PHONY: all setup services stop local docker watch dev clean help run tidy vendor build

# Variables
GO_CMD = go
GORELEASER_CMD = goreleaser
DOCKER_COMPOSE = docker-compose -f debug.docker-compose.yml
AIR_CMD = air

# Colors
COLOR_RESET = \033[0m
COLOR_CYAN = \033[36m
COLOR_GREEN = \033[32m

all: help

setup: ## Install dependencies and prepare development environment
	@echo "$(COLOR_CYAN)Setting up development environment...$(COLOR_RESET)"
	$(GO_CMD) mod tidy
	$(GO_CMD) install github.com/air-verse/air@latest

services: ## Start MongoDB and Redis in containers
	@echo "$(COLOR_CYAN)Starting MongoDB and Redis services...$(COLOR_RESET)"
	$(DOCKER_COMPOSE) up -d mongodb redis

stop: ## Stop all containers
	@echo "$(COLOR_CYAN)Stopping all services...$(COLOR_RESET)"
	$(DOCKER_COMPOSE) down

local: setup services ## Run app locally with containerized services
	@echo "$(COLOR_CYAN)Running app locally...$(COLOR_RESET)"
	@export $$(cat .env.local | xargs) && \
	export DEBUG=true && \
	$(GO_CMD) run main.go

docker: ## Run everything in containers
	@echo "$(COLOR_CYAN)Running in Docker...$(COLOR_RESET)"
	$(DOCKER_COMPOSE) up --build

watch: setup services ## Run locally with hot reload
	@echo "$(COLOR_CYAN)Running with hot reload...$(COLOR_RESET)"
	@export DEBUG=true && \
	$(AIR_CMD) -c .air.toml

dev: setup services watch ## Setup, start services and run with hot reload

run: ## Run the application
	$(GO_CMD) run main.go

tidy: ## Tidy Go modules
	$(GO_CMD) mod tidy

vendor: ## Vendor Go dependencies
	$(GO_CMD) mod vendor

build: ## Build release with goreleaser
	$(GORELEASER_CMD) release --snapshot --skip=publish --clean --skip=sign

clean: stop ## Clean up generated files and stopped containers
	@echo "$(COLOR_CYAN)Cleaning up...$(COLOR_RESET)"
	rm -rf tmp/
	rm -f alita_robot

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_GREEN)%-15s$(COLOR_RESET) %s\n", $$1, $$2}'
