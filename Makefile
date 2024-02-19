.SILENT:

ENV_FILE=.env
-include $(ENV_FILE)
export $(shell sed 's/#.*//' $(ENV_FILE) | xargs)

# Binaries
GO=go
PKL=pkl

build-cli:
	$(GO) build -o bin/velocity cmd/cli/main.go

build-%:
	$(GO) build -o bin/$* cmd/$*/main.go

run-%:
	$(GO) run cmd/$*/main.go

pkl-gen: clean-pkl
	pkl-gen-go --generator-settings=pkl/generator-settings.pkl pkl/velocity.pkl
	for file in pkl/prebuilts/*.pkl; do \
		pkl-gen-go "$$file"; \
	done

pkl-test:
	$(PKL) test pkl/tests/sections/** pkl/tests/prebuilts/**

pkl-eval:
	$(PKL) eval self.pkl

clean: clean-pkl
	rm -rf bin

clean-pkl:
	rm -rf gen/pkl
