# Конфигурация
APP_NAME := main
CMD_PATH := ./cmd/app/main.go
DOCKER_COMPOSE := docker-compose

ifeq (, $(shell which $(DOCKER_COMPOSE)))
$(error "No $(DOCKER_COMPOSE) in $(PATH), consider installing Docker Compose")
endif

GOLANGCI_LINT := golangci-lint
ifeq (, $(shell which $(GOLANGCI_LINT)))
$(error "No $(GOLANGCI_LINT) in $(PATH), consider installing golangci-lint")
endif

.PHONY: build docker-build docker-start lint format clean

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(APP_NAME) $(CMD_PATH)

docker-build:
	@echo "Building Docker images..."
	@$(DOCKER_COMPOSE) build

docker-start:
	@echo "Starting Docker containers..."
	@$(DOCKER_COMPOSE) up -d

lint:
	@echo "Running linters..."
	@$(GOLANGCI_LINT) run ./...

format:
	@echo "Formatting code..."
	@go fmt ./...

clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)
	@$(DOCKER_COMPOSE) down --remove-orphans
