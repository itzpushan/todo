#  Loading .env if available 
-include .env
export $(shell sed 's/=.*//' .env)

# Fallback defaults (if .env vars not set) 
DATABASE_URL ?= postgresql://postgres:pushan@localhost:5432/postgres?sslmode=disable
MIGRATIONS_DIR ?= file://ent/migrate/migrations

# Commands

## Run server
dev:
	go run cmd/server/main.go

## seedint the database
seed:
	go run cmd/seed/main.go

## Appling DB migrations with Atlas
migrate:
	atlas migrate apply --dir "$(MIGRATIONS_DIR)" --url "$(DATABASE_URL)"

## Generating Ent client
generate:
	go generate ./ent

.PHONY: dev seed migrate generate clean
