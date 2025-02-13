.PHONY: fa fmt

fa:
	@fieldalignment -fix ./...

fmt:
	@goimports -w -local github.com/pixel365/bx .
	@gofmt -w .
	@golines -w .