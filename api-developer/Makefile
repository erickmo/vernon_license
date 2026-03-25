.PHONY: infra-up infra-down dev dev-wasm migrate-up migrate-down tidy test lint build build-wasm build-all

# ─── Infrastructure ────────────────────────────────────────────────
infra-up:
	docker-compose up -d postgres

infra-down:
	docker-compose down

# ─── Development ───────────────────────────────────────────────────
dev:
	go run ./cmd/api

dev-wasm:
	@chmod +x scripts/watch-wasm.sh
	@bash scripts/watch-wasm.sh

build:
	go build -o bin/api ./cmd/api

# ─── WASM Build ────────────────────────────────────────────────────
build-wasm:
	GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/api

# ─── Full Build (server + WASM) ────────────────────────────────────
build-all: build build-wasm

# ─── Database Migrations ───────────────────────────────────────────
migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

# ─── Dependencies ──────────────────────────────────────────────────
tidy:
	go mod tidy

# ─── Testing ───────────────────────────────────────────────────────
test:
	go test ./... -v -race -coverprofile=coverage.out

test-cover:
	go tool cover -html=coverage.out

# ─── Lint ──────────────────────────────────────────────────────────
lint:
	golangci-lint run ./...
