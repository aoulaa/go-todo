include .env

create_migration:
	migrate create -ext=sql -dir=internal/db/migrations -seq init

migrate_up:
	migrate -path=internal/db/migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=internal/db/migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose down

.PHONY: create_migration migrate_up migrate_down

run:
	${HOME}/go/bin/air