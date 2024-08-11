# Binaries
GO=go

BUILD_DIR=$(HOME)/go/bin

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
	docker compose -f docker-compose.dev.yml up -d

dev-down:
	docker compose -f docker-compose.dev.yml down

prod:
	docker compose -f docker-compose.prod.yml up -d

prod-down:
	docker compose -f docker-compose.prod.yml down
