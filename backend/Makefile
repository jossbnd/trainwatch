.PHONY: fmt lint test

fmt:
	gofumpt -w ./cmd ./internal
	goimports -w -local github.com/jossbnd/trainwatch ./cmd ./internal

lint:
	golangci-lint run ./...

test:
	go test ./...
