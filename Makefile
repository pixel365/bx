.PHONY: all fa fmt lint test

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