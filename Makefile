.SILENT:
PACKAGES=$(shell go list ./... | grep -v /vendor/)
CONFIG_FILE=velocity.yml
ENV_FILE=.env
ENV=$(shell awk '{gsub(/#.*/, ""); printf "%s ", $$0}' $(ENV_FILE))

test:
	$(foreach package,$(PACKAGES), \
		go test -v $(package); \
	)

build-cli:
	go build -o build/velocity cmd/cli/main.go
	mv build/velocity $(GOPATH)/bin

agent:
	$(ENV) go run cmd/agent/main.go

test-db:
	$(ENV) go run cmd/test-db/main.go

test-basicworkflow:
	$(ENV) go run cmd/test-basicworkflow/main.go

example:
	$(ENV) go run cmd/example/main.go

packages:
	echo $(PACKAGES)






