.PHONY: fmt lint test up reup deps

fmt:
	gofumpt -w ./cmd ./internal
	goimports -w -local github.com/jossbnd/trainwatch ./cmd ./internal

lint:
	golangci-lint run ./...

test:
	go test ./...

up:
	docker compose -f docker/docker-compose.yml up -d traefik && docker compose -f docker/docker-compose.yml up backend

reup:
	docker compose -f docker/docker-compose.yml up -d traefik && docker compose -f docker/docker-compose.yml up --force-recreate --build backend

deps:
	go mod tidy
	go mod vendor
