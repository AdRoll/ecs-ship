.PHONY: mod lint test build mockgen

mod:
	go mod download
	go mod vendor

lint:
	golangci-lint run

test:
	go test -cover -v `go list ./...`

build:
	go build -o artifacts/ecs-deploy ./main.go

mockgen:
	mockgen -source=./action/ecs-deploy.go -destination=./action/mock/ecs-deploy_mock.go
