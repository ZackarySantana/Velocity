.SILENT:
PACKAGES=$(shell go list ./... | grep -v /vendor/)
CONFIG_FILE=velocity.yml

test:
	$(foreach package,$(PACKAGES), \
		go test -v $(package); \
	)

workflows:
	go run cmd/workflows/main.go $(CONFIG_FILE)

packages:
	echo $(PACKAGES)

# cover-xml:
#     @$(foreach package, $(packages), \
#         gocov convert $(package)/cover.out | gocov-xml > $(package)/coverage.xml;)







