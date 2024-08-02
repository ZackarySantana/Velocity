build-cli:
	go build -o $(HOME)/go/bin/velocity cmd/cli/main.go

cli:
	go run cmd/cli/main.go

build-api:
	go build -o $(HOME)/go/bin/velocity-api cmd/api/main.go

api:
	go run cmd/api/main.go