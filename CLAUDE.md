# Vernon License — Centralized Licensing System

## Overview

Vernon License = **license registry + kill switch** untuk semua produk klien Vernon.

- **Public API**: hanya **2 endpoint** — register + validate. Dipanggil oleh client apps.
- **App** (`app-developer/`): Go WASM PWA — mengelola companies, projects, products, proposals, licenses, users. Bicara langsung ke DB (tidak lewat public API).

### How It Works

1. Client app (e.g. FlashERP instance) **registers itself** → Vernon returns license key
2. Client app **periodically validates** license → Vernon returns `true/false`
3. If `false` → client app denies user access
4. Vernon App (internal) mengelola semua data: companies, projects, proposals, products, users
5. PO approve proposal → auto-create license. Atau PO direct create license
6. PO suspend license → next validate call returns `false` → client app blocks users

## Stack

| Layer | Tech |
|---|---|
| Backend | Go 1.23 · Chi v5 · Uber FX v1.22 · PostgreSQL 17 · sqlx |
| Frontend | Go 1.23 · go-app/v10 (WASM) · PWA |
| Auth (App) | JWT HS256 (24h) — hanya untuk Vernon App, bukan public API |

## Public API (2 endpoints only)

```
POST /api/v1/register    ← Client app registers itself → returns license
GET  /api/v1/validate    ← Client app checks license → returns { valid: true/false }
```

Semua management lainnya (companies, projects, proposals, products, users, dashboard) **BUKAN** API — dikelola di App yang bicara langsung ke DB.

## Roles (Vernon App)

| Role | Deskripsi |
|---|---|
| `superuser` | Full access + manage users + manage products + global audit |
| `project_owner` | Create license directly, edit/approve proposals, suspend/activate/renew |
| `sales` | Create proposal → submit for review. Tidak bisa approve atau buat lisensi langsung |

## Architecture Rules (API — 2 public endpoints)

- Stateless — no session, no JWT pada public API
- Register: validasi `provision_api_key`, return license data
- Validate: cek `license_key` + status, return boolean
- Rate limit: 60 req/min per IP

## Architecture Rules (App — internal, talks to DB)

- **Framework**: go-app/v10 (`app.Compo`)
- **Data access**: langsung ke PostgreSQL via repository layer (bukan lewat API)
- **Auth**: JWT HS256 untuk login App. Stored di `localStorage`
- **Routing**: `app.Route()` di `main_wasm.go`

## Key Commands

```bash
# API (api-developer/)
make infra-up       # Start PostgreSQL (port 5433)
make dev            # Run server (port 8081) — serves public API + App WASM
make migrate-up     # Apply migrations
make tidy           # go mod tidy
make test           # go test ./... -v -race
```

## Environment Variables

```bash
DATABASE_URL=postgres://vernon:secret@localhost:5433/vernon_license?sslmode=disable
JWT_SECRET=your-256-bit-secret            # Untuk Vernon App login saja
PORT=8081
LOG_LEVEL=info
STORAGE_PATH=./storage
LICENSE_CHECK_INTERVAL=6h
COMPANY_NAME=FlashLab
COMPANY_ADDRESS=Jl. Teknologi No. 1, Jakarta
COMPANY_PHONE=+62-21-1234567
COMPANY_EMAIL=hello@flashlab.id
COMPANY_LOGO_PATH=assets/flashlab-logo.png
```

## Design Tokens

```
Primary:   #4D2975     Success: #22C55E
Accent:    #26B8B0     Error:   #EF4444
Secondary: #E9A800     Warning: #F59E0B
```

## Detailed Documentation

| Saat mengerjakan... | Baca file |
|---|---|
| Public API (register + validate) | `docs/PUBLIC_API.md` |
| Domain models | `docs/DOMAIN_MODELS.md` |
| App features, pages | `docs/APP_FEATURES.md` |
| Database migrations | `docs/MIGRATIONS.md` |
| Error codes | `docs/ERROR_CODES.md` |
| Visual reference | `docs/visualization.jsx` |
