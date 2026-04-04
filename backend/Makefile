.PHONY: fmt lint test up reup deps

fmt:
	gofumpt -w ./cmd ./internal
	goimports -w -local github.com/jossbnd/trainwatch ./cmd ./internal

lint:
	golangci-lint run ./...

test:
	go test ./...

up:
	docker compose up -d traefik && docker compose up trainwatch-backend

reup:
	docker compose up -d traefik && docker compose up --force-recreate --build trainwatch-backend

deps:
	go mod tidy
	go mod vendor