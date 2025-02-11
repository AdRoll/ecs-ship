.PHONY: mod lint test build mockgen

mod:
	go mod download
	go mod tidy

lint:
	golangci-lint run

test:
	go test -cover -v `go list ./...`

build:
	go build .

mockgen:
	go generate ./...
