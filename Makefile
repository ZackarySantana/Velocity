.SILENT:
PACKAGES=$(shell go list ./... | grep -v /vendor/)

test:
	$(foreach package,$(PACKAGES), \
		go test -v $(package); \
	)

packages:
	echo $(PACKAGES)

# cover-xml:
#     @$(foreach package, $(packages), \
#         gocov convert $(package)/cover.out | gocov-xml > $(package)/coverage.xml;)







