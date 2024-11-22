include .env # to have access to the environment variables

migrate-create:
	migrate create -ext sql -dir migrations 'migrate_name'
.PHONY: migrate-create

migrate-up:
	migrate -path migrations -database '$(PG_DSN)?sslmode=disable' up
.PHONY: migrate-up

migrate-down:
	migrate -path migrations -database '$(PG_DSN)?sslmode=disable' down
.PHONY: migrate-down

run:
	go mod tidy && go mod download && CGO_ENABLED=0 go run ./cmd/main.go
.PHONY: run