.SILENT:
PACKAGES=$(shell go list ./... | grep -v /vendor/)
CONFIG_FILE=velocity.yml

test:
	$(foreach package,$(PACKAGES), \
		go test -v $(package); \
	)

cli:
	go build -o build/velocity cmd/cli/main.go
	mv build/velocity $(GOPATH)/bin

workflows:
	go run cmd/workflows/main.go $(CONFIG_FILE)

packages:
	echo $(PACKAGES)

# cover-xml:
#     @$(foreach package, $(packages), \
#         gocov convert $(package)/cover.out | gocov-xml > $(package)/coverage.xml;)







