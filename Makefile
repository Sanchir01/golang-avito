PHONY:
SILENT:

MIGRATION_NAME ?= new_migration

PASSWORD ?= password
lint:
	golangci-lint run --config=golangci.yaml


build:
	go build -o ./.bin/main ./cmd/main/main.go

run: build	lint
	./.bin/main

docker-build:
	docker build -t avito-httpserver .