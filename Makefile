# Binaries
GO=go

OUT=$(HOME)/go/bin

# Default services
ID_CREATOR?=mock
REPOSITORY_MANAGER?=mock
PROCESS_QUEUE?=mock
PRIORITY_QUEUE?=mock

export ID_CREATOR
export REPOSITORY_MANAGER
export PROCESS_QUEUE
export PRIORITY_QUEUE

# Default variables
PORT?=8080

export PORT

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
	$(GO) build -o $(OUT)/velocity cmd/cli/main.go

cli:
	$(GO) run cmd/cli/main.go

build-api:
	$(GO) build -o $(OUT)/velocity-api cmd/api/main.go

api:
	$(GO) run cmd/api/main.go

build-agent:
	$(GO) build -o $(OUT)/velocity-agent cmd/agent/main.go

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
