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

build:
	go build -o ./.bin/main ./cmd/main/main.go

run: build
	./.bin/main


migrations-new:
	goose -dir $(DIR_PG) create $(MIGRATION_NAME) sql

migrations-up-pg:
	goose -dir $(DIR_PG) postgres  $(POSTGRES_DSN) up

migrations-down-pg:
	goose -dir $(DIR_PG) postgres  $(POSTGRES_DSN) down

migrations-status-pg:
	goose -dir $(DIR_PG) postgres  $(POSTGRES_DSN) status

docker-build:
	docker build -t avito-httpserver .