SERVICE=flasherp-developer-api
MAIN=./cmd/api
BIN=./bin/$(SERVICE)

.PHONY: build dev test lint migrate-up migrate-down infra-up infra-down tidy

build:
	go build -o $(BIN) $(MAIN)

dev:
	go run $(MAIN)

test:
	go test ./... -v -timeout 30s

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

infra-up:
	docker compose up -d

infra-down:
	docker compose down

tidy:
	go mod tidy
