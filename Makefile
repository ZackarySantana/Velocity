.SILENT:
PACKAGES=$(shell go list ./... | grep -v /vendor/)
ENV_FILE=.env
ENV=$(shell awk '{gsub(/#.*/, ""); printf "%s ", $$0}' $(ENV_FILE))

test:
	$(foreach package,$(PACKAGES), \
		go test -v $(package); \
	)

build-cli:
	go build -o $(GOPATH)/bin/velocity cmd/cli/main.go

agent:
	$(ENV) go run cmd/agent/main.go

server:
	$(ENV) go run cmd/server/main.go

example:
	$(ENV) go run cmd/example/main.go

packages:
	echo $(PACKAGES)

clean:
	rm $(GOPATH)/bin/velocity






