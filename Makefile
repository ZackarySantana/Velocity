.SILENT:

ENV_FILE=.env
ENV=$(shell awk '{gsub(/#.*/, ""); printf "%s ", $$0}' $(ENV_FILE))

# Experimental pkl generation
pkl-gen:
	pkl-gen-go pkl/velocity.pkl

build-cli:
	go build -o bin/velocity cmd/cli/main.go

agent:
	$(ENV) go run cmd/agent/main.go

build-agent:
	go build -o bin/agent cmd/agent/main.go

api:
	$(ENV) go run cmd/api/main.go

build-api:
	go build -o bin/api cmd/api/main.go

ui:
	$(ENV) go run cmd/ui/main.go

build-ui:
	go build -o bin/ui cmd/ui/main.go

clean:
	rm -r bin