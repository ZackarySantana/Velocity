# Binaries
GO=go

BUILD_DIR=$(HOME)/go/bin

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build-cli       Build the CLI"
	@echo "  cli             Run the CLI"
	@echo "  build-api       Build the API"
	@echo "  api             Run the API"
	@echo "  build-agent     Build the Agent"
	@echo "  agent           Run the Agent"
	@echo "  dev             Run the development environment"
	@echo "  dev-%           Run a specific command in the development environment"
	@echo "  prod            Run the production environment"
	@echo "  prod-%          Run a specific command in the production environment"

build-cli:
	$(GO) build -o $(BUILD_DIR)/velocity cmd/cli/main.go

cli:
	$(GO) run cmd/cli/main.go

build-api:
	$(GO) build -o $(BUILD_DIR)/velocity-api cmd/api/main.go

api:
	$(GO) run cmd/api/main.go

build-agent:
	$(GO) build -o $(BUILD_DIR)/velocity-agent cmd/agent/main.go

agent:
	$(GO) run cmd/agent/main.go

dev:
	docker compose -f compose.yml -f compose.dev.yml --profile dev up -d

dev-%:
	docker compose -f compose.yml -f compose.dev.yml --profile dev $* $(ARGS)

prod:
	docker compose -f compose.yml -f compose.prod.yml --profile prod up -d

prod-%:
	docker compose -f compose.yml -f compose.prod.yml --profile prod $* $(ARGS)
