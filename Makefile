PHONY:
SILENT:
include .env
export
MIGRATION_NAME ?= new_migration

POSTGRES_USER := postgres
POSTGRES_PASSWORD ?= avito
POSTGRES_HOST := localhost
POSTGRES_PORT := 5435
POSTGRES_DB := postgres

POSTGRES_DSN := "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

DIR_PG = ./migrations
lint:
	golangci-lint run --config=golangci.yaml
swag:
	swag init -g cmd/main/main.go
build:
	go build -o ./.bin/main ./cmd/main/main.go

run: build	swag
	./.bin/main


migrations-new:
	goose -dir $(DIR_PG) create $(MIGRATION_NAME) sql

migrations-up-pg:
	goose -dir $(DIR_PG) postgres  $(POSTGRES_DSN) up

migrations-down-pg:
	goose -dir $(DIR_PG) postgres  $(POSTGRES_DSN) down

migrations-status-pg:
	goose -dir $(DIR_PG) postgres  $(POSTGRES_DSN) status

compose:
	docker-compose up -d
docker-build:
	docker build -t avito-httpserver .

test:
	go test -v ./tests/...
