.PHONY: up
up:
	docker compose up -d

.PHONY: up-db
up-db:
	goose -dir migrations postgres "postgres://support-assistant-api:very-secure@localhost:5432/support-assistant-api" up

.PHONY: build
build:
	go build ./cmd/support-assistant-api-server && go build ./cmd/support-assistant-load-history

.PHONY: migrate-history
migrate-history:
	./support-assistant-load-history

.PHONY: run
run:
	./support-assistant-api-server

.PHONY: psql
psql:
	docker compose exec postgres-pgvector psql -U support-assistant-api
