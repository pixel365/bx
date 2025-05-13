.PHONY: all fa fmt lint test build cover

all: fa fmt lint test

fa:
	@fieldalignment -fix ./...

fmt:
	@goimports -w -local github.com/pixel365/bx .
	@gofmt -w .
	@golines -w .

lint:
	@golangci-lint run

test:
	@go $@ ./...

build:
	@go $@ -o ./bin/bx -ldflags="-s -w"

cover:
	go test -coverprofile=coverage.out ./... && go tool $@ -html=coverage.out

coverfn:
	go test -coverprofile=coverage.out ./... && \
	go tool cover -func=coverage.out

doc:
	docsify serve docs
