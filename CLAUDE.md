# Vernon License

License registry + kill switch untuk produk klien Vernon.

## Stack
Go 1.23 · Chi v5 · Uber FX · PostgreSQL 17 (port 5433) · sqlx · go-app/v10 (WASM) · JWT HS256

## Structure
```
api-developer/
├── cmd/api/main.go          # server entrypoint (!wasm)
├── cmd/api/app.go           # server-side route reg (!wasm)
├── cmd/api/main_wasm.go     # WASM entrypoint + client routes
├── infrastructure/database/ # sqlx repos
├── internal/
│   ├── domain/              # structs + repository interfaces
│   ├── handler/             # internal API handlers (JWT)
│   ├── publicapi/           # register + validate
│   ├── service/             # business logic
│   └── ui/pages/            # WASM pages
├── migrations/              # 001–016, sql-migrate format
├── pkg/licenseutil/         # GenerateOTP, GenerateLicenseKey
└── web/app.wasm
```

## Commands
```bash
make infra-up       # docker postgres port 5433
make migrate-up     # apply migrations
make dev            # go run ./cmd/api → :8081
make test           # go test ./... -v -race
make build-wasm     # GOARCH=wasm GOOS=js go build -o web/app.wasm ./cmd/api
make build          # go build -o bin/api ./cmd/api
```

## Public API (no auth)
```
POST /api/v1/register   { otp, instance_url, instance_name, product_slug }
GET  /api/v1/validate   ?key=FL-XXXXXXXX
```
Rate limit: 60 req/min per IP.

## Internal API (Bearer JWT)
```
POST /api/internal/auth/login
GET  /api/internal/companies               POST /api/internal/companies
GET  /api/internal/licenses                POST /api/internal/licenses
GET  /api/internal/licenses/{id}
PUT  /api/internal/licenses/{id}/activate|suspend|renew|status|constraints
GET  /api/internal/licenses/{id}/otp       (superuser)
GET  /api/internal/proposals               POST /api/internal/proposals
PUT  /api/internal/proposals/{id}/submit|approve|reject
GET  /api/internal/products                POST /api/internal/products
GET  /api/internal/dashboard
GET  /api/internal/dashboard/otp              ← lightweight OTP-only (AJAX refresh)
GET  /api/internal/notifications
```

## Roles
| Role | Akses |
|---|---|
| `superuser` | Full + products + users + audit |
| `project_owner` | Create license, approve proposals, suspend/activate/renew |
| `sales` | Create & submit proposal only |

## Domain
`Company` → `Project` → `License` + `Proposal`

License status: `pending → active ↔ suspended`, `expired`, `trial`

OTP: per-license (registration auth). Global `otp` table: rotating dashboard display code.

## WASM Routing — Critical Rule
Parameterized routes HARUS `app.RouteWithRegexp`, bukan `app.Route` (exact match only).
```go
// WRONG: app.Route("/licenses/{id}", ...)  ← only matches literal "/licenses/{id}"
// RIGHT:
app.RouteWithRegexp(`^/licenses/[^/]+$`, func() app.Composer { return &pages.LicenseDetailPage{} })
```
Daftarkan di **dua file**: `cmd/api/app.go` (server) + `cmd/api/main_wasm.go` (client).

## Env
```
DATABASE_URL=postgres://vernon:secret@localhost:5433/vernon_license?sslmode=disable
JWT_SECRET=<32+ chars>
PORT=8081
LOG_LEVEL=info
```

## Docs
| Topic | File |
|---|---|
| Public API | `docs/PUBLIC_API.md` |
| Domain models + schema | `docs/DOMAIN_MODELS.md` |
| App pages + features | `docs/APP_FEATURES.md` |
| Auth + error codes | `docs/AUTH_AND_FLOWS.md` |
