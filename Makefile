.SILENT:
PACKAGES=$(shell go list ./... | grep -v /vendor/)
ENV_FILE=.env
ENV=$(shell awk '{gsub(/#.*/, ""); printf "%s ", $$0}' $(ENV_FILE))

test:
	$(foreach package,$(PACKAGES), \
		go test -v $(package); \
	)

build-cli:
	go build -o ./bin/velocity cmd/cli/main.go

agent:
	$(ENV) go run cmd/agent/main.go

agent-mongodb:
	$(ENV) MONGODB_AGENT=true go run cmd/agent/main.go

server:
	$(ENV) go run cmd/server/main.go

indexes:
	$(ENV) go run cmd/indexes/main.go

packages:
	echo $(PACKAGES)

clean:
	rm ./bin/velocity
