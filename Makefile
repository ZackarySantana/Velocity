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

# add_service_compose_file returns the compose file for the service if
# the service is not a mock service.
define add_service_compose_file
	$(if $(filter-out mock,$(1)),-f compose/compose.dev.$(1).yml)
endef

# Dev compiles the compose files for the needed services and runs the dev profile.
dev:
	docker compose \
		-f compose.yml \
		-f compose/compose.dev.yml \
		-f compose/compose.dev.grafana.yml \
		$(call add_service_compose_file,$(ID_CREATOR)) \
		$(call add_service_compose_file,$(REPOSITORY_MANAGER)) \
		$(call add_service_compose_file,$(PROCESS_QUEUE)) \
		$(call add_service_compose_file,$(PRIORITY_QUEUE)) \
		--profile dev \
		-p velocity \
		up -d

dev-down:
	docker compose -p velocity down

mongo-dev:
	ID_CREATOR=mongodb \
	REPOSITORY_MANAGER=mongodb \
	PROCESS_QUEUE=mongodb \
	PRIORITY_QUEUE=mongodb \
	$(MAKE) dev

# Dev-% compiles the compose files for the needed services and runs the dev profile with the given command.
dev-%:
	docker compose \
		-f compose.yml \
		-f compose/compose.dev.yml \
		-f compose/compose.dev.grafana.yml \
		$(call add_service_compose_file,$(ID_CREATOR)) \
		$(call add_service_compose_file,$(REPOSITORY_MANAGER)) \
		$(call add_service_compose_file,$(PROCESS_QUEUE)) \
		$(call add_service_compose_file,$(PRIORITY_QUEUE)) \
		--profile dev \
		$* $(ARGS)

prod:
	docker compose \
		-f compose.yml \
		-f compose/compose.prod.yml \
		--profile prod \
		-p velocity-prod \
		up -d

prod-down:
	docker compose -p velocity-prod down

prod-%:
	docker compose \
		-f compose.yml \
		-f compose/compose.prod.yml \
		--profile prod \
		$* $(ARGS)
