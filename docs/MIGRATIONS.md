# Migrations

Format: sql-migrate. Files di `api-developer/migrations/`.

```bash
make migrate-up    # apply semua pending
make migrate-down  # rollback 1 step
```

## File List

| File | Isi |
|---|---|
| 001 | `users` table |
| 002 | `companies` table |
| 003 | `projects` table |
| 004 | `products` table |
| 005 | `client_licenses` table — includes `otp`, `otp_previous`, constraints |
| 006 | `proposals` table |
| 007 | `audit_logs` table |
| 008 | `notifications` table |
| 009 | no-op (merged ke 005) |
| 010 | `otp` table (global rotating code) |
| 011 | no-op (merged ke 005) |
| 012 | no-op (merged ke 005) |
| 013 | Drop FK + nullable `project_id`, `company_id` di client_licenses |
| 014 | Nullable `created_by` di client_licenses |
| 015 | Nullable `created_by` di companies |
| 016 | Unique constraint `company_id + product_id` di client_licenses |

## Notes
- Migrations 009, 011, 012 adalah no-op — kolom OTP sudah langsung dibuat di 005
- `client_licenses.project_id` dan `company_id` nullable sejak 013 (untuk register flow)
- Fresh DB: run `make migrate-up` dari 001–016
- Existing DB: jalankan `make migrate-up`, hanya pending migrations yang akan dijalankan
