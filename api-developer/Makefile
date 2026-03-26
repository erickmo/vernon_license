.PHONY: infra-up infra-down dev dev-wasm migrate-up migrate-down tidy test lint build build-wasm build-all reset

# ─── Infrastructure ────────────────────────────────────────────────
infra-up:
	docker-compose up -d postgres

infra-down:
	docker-compose down

# ─── Development ───────────────────────────────────────────────────
dev:
	@lsof -ti :8081 | xargs -r kill -9 2>/dev/null || true
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

# ─── Reset (drop DB + delete .env) ────────────────────────────────
reset:
	@if [ -f .env ]; then \
		DB_URL=$$(grep '^DATABASE_URL=' .env | cut -d'=' -f2-); \
		if [ -n "$$DB_URL" ]; then \
			DB_NAME=$$(echo "$$DB_URL" | sed -n 's|.*/\([^?]*\).*|\1|p'); \
			DB_HOST=$$(echo "$$DB_URL" | sed -n 's|.*@\([^:]*\).*|\1|p'); \
			DB_PORT=$$(echo "$$DB_URL" | sed -n 's|.*:\([0-9]*\)/.*|\1|p'); \
			DB_USER=$$(echo "$$DB_URL" | sed -n 's|.*://\([^:]*\).*|\1|p'); \
			echo "Dropping database '$$DB_NAME' on $$DB_HOST:$$DB_PORT..."; \
			PGPASSWORD=$$(echo "$$DB_URL" | sed -n 's|.*://[^:]*:\([^@]*\)@.*|\1|p') \
				dropdb -h "$$DB_HOST" -p "$$DB_PORT" -U "$$DB_USER" --if-exists "$$DB_NAME"; \
			echo "Database '$$DB_NAME' dropped."; \
		fi; \
		rm -f .env; \
		echo ".env deleted."; \
	else \
		echo "No .env file found, nothing to reset."; \
	fi

# ─── Lint ──────────────────────────────────────────────────────────
lint:
	golangci-lint run ./...
