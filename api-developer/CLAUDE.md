# FlashERP Developer API

Server manajemen lisensi FlashLab. Mengelola lisensi klien dan provisioning ke FlashERP instances.

**Digunakan oleh:** app-developer Flutter app (role: developer_sales)

## Stack
- Go 1.23 + Chi v5 + Uber FX v1.22
- PostgreSQL 17 + sqlx
- JWT HS256

## Architecture Rules
- Domain layer: ZERO external dependency
- Command handler: hanya WriteRepository
- Query handler: hanya ReadRepository
- HTTP handler: hanya CommandBus + QueryBus (kecuali login dan create yang return result)

## API Endpoints
- POST /api/v1/auth/login
- GET  /api/v1/auth/me
- GET  /api/v1/licenses
- POST /api/v1/licenses
- GET  /api/v1/licenses/{id}
- PUT  /api/v1/licenses/{id}/status
- PUT  /api/v1/licenses/{id}/constraints
- POST /api/v1/licenses/{id}/provision

## Key Commands
```bash
make infra-up    # start PostgreSQL di port 5433
make dev         # run server (port 8081)
make migrate-up  # apply migrations
make tidy        # go mod tidy
```

## Provisioning Flow
1. developer_sales buat license baru → sistem generate FL-XXXXXXXX key + initial password
2. Set flasherp_url + provision_api_key untuk dedicated deployment
3. POST /api/v1/licenses/{id}/provision → push ke FlashERP /internal/provision
4. Kirim credentials (email + initial_password) ke klien

## Notes
- Port 5433 (bukan 5432) untuk tidak konflik dengan FlashERP utama
- Initial password hanya tampil sekali di response POST /api/v1/licenses
- CommandBus key format: "package_name.TypeName"
