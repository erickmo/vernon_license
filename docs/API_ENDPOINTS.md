# API Endpoints

## Public (no auth)

```
GET  /api/v1/setup/status                    ← Cek apakah system sudah di-setup
POST /api/v1/setup/install                   ← First-run setup (buat superuser)
POST /api/v1/auth/login                      ← Login, dapat JWT
GET  /api/v1/client/license?key=FL-XXXXX     ← Klien ambil data lisensinya sendiri *
```

> \* Validasi: `Origin`/`Referer` header **harus cocok** dengan `flasherp_url` di lisensi. Tidak cocok → 403.

## Protected (Bearer JWT)

```
GET  /api/v1/auth/me
GET  /api/v1/dashboard                       ← Analytics dashboard (lihat Dashboard section)
```

## Protected + Role: `sales` | `project_owner` | `superuser`

### Products (read)
```
GET    /api/v1/products                      ← List semua produk aktif (untuk dropdown di form)
GET    /api/v1/products/{id}                 ← Detail produk + available modules/apps
```

### Products — Role: `project_owner` | `superuser` only
```
POST   /api/v1/products                      ← Buat produk baru
PUT    /api/v1/products/{id}                 ← Update produk
DELETE /api/v1/products/{id}                 ← Soft delete produk
```

### Licenses
```
GET    /api/v1/licenses                      ← List (paginated, filterable, sortable)
POST   /api/v1/licenses                      ← Buat lisensi baru + auto-generate proposal PDF
GET    /api/v1/licenses/{id}                 ← Detail (by UUID atau license key)
PUT    /api/v1/licenses/{id}/constraints     ← Update batasan (max_users, dll)
GET    /api/v1/licenses/{id}/audit           ← Audit trail untuk lisensi ini
GET    /api/v1/licenses/export               ← Export CSV
```

### Proposals
```
GET    /api/v1/licenses/{id}/proposal        ← Download proposal PDF
POST   /api/v1/licenses/{id}/proposal/regenerate  ← Regenerate proposal (jika data berubah)
```

### Licenses — Role: `project_owner` | `superuser` only
```
PUT    /api/v1/licenses/{id}/approve         ← Approve client (trial → active)
PUT    /api/v1/licenses/{id}/status          ← Suspend / activate
POST   /api/v1/licenses/{id}/renew           ← Perpanjang lisensi
POST   /api/v1/licenses/{id}/provision       ← Push ke FlashERP instance
POST   /api/v1/licenses/bulk/status          ← Bulk update status
```

### Notifications
```
GET    /api/v1/notifications                 ← List (paginated, unread first)
PUT    /api/v1/notifications/{id}/read       ← Tandai sudah dibaca
PUT    /api/v1/notifications/read-all        ← Tandai semua sudah dibaca
GET    /api/v1/notifications/unread-count    ← Badge count
POST   /api/v1/devices                       ← Register FCM token
DELETE /api/v1/devices/{token}               ← Unregister device
```

### Audit (superuser only)
```
GET    /api/v1/audit                         ← Global audit log
         ?entity_type=client_license
         &action=status_changed
         &actor_id=<uuid>
         &from=2026-01-01&to=2026-03-24
         &page=1&per_page=50
```

---

## Pagination (standar semua list endpoint)

### Query Parameters

| Param | Type | Default | Deskripsi |
|---|---|---|---|
| `page` | int | `1` | Halaman (1-indexed) |
| `per_page` | int | `20` | Item per halaman (max: 100) |
| `sort_by` | string | `created_at` | Field sort |
| `order` | string | `desc` | `asc` \| `desc` |
| `search` | string | — | ILIKE on `client_name` |
| `status` | string | — | `active` \| `trial` \| `suspended` \| `expired` |
| `product` | string | — | Filter by product slug (dynamic, from `GET /api/v1/products`) |
| `plan` | string | — | `saas` \| `dedicated` |

### Response Envelope

```json
{
  "data": [ ... ],
  "meta": {
    "page": 2,
    "per_page": 20,
    "total": 142,
    "total_pages": 8
  }
}
```

---

## Endpoint Details

### GET /api/v1/products — List produk

```json
// Response
{
  "data": [
    {
      "id": "uuid",
      "name": "FlashERP",
      "slug": "flasherp",
      "description": "Sistem ERP terintegrasi",
      "available_modules": [
        { "key": "accounting", "name": "Accounting", "description": "Modul akuntansi" },
        { "key": "sales", "name": "Sales", "description": "Modul penjualan" }
      ],
      "available_apps": [
        { "key": "web-ui", "name": "Web UI", "description": "Antarmuka web utama" },
        { "key": "app-pos", "name": "App POS", "description": "Aplikasi point of sale" }
      ],
      "available_plans": ["saas", "dedicated"],
      "base_pricing": {
        "saas": { "base_price": 5000000, "per_user_price": 100000, "currency": "IDR" },
        "dedicated": { "base_price": 25000000, "per_user_price": 75000, "currency": "IDR" }
      },
      "is_active": true
    }
  ]
}
```

### POST /api/v1/products — Buat produk baru (`project_owner`/`superuser`)

```json
// Request
{
  "name": "FlashHRD",
  "slug": "flashhrd",
  "description": "Sistem manajemen SDM",
  "available_modules": [
    { "key": "attendance", "name": "Attendance", "description": "Absensi karyawan" },
    { "key": "payroll", "name": "Payroll", "description": "Penggajian" }
  ],
  "available_apps": [
    { "key": "web-ui", "name": "Web UI", "description": "Antarmuka web" },
    { "key": "app-employee", "name": "App Employee", "description": "Aplikasi karyawan" }
  ],
  "available_plans": ["saas", "dedicated"],
  "base_pricing": {
    "saas": { "base_price": 3000000, "per_user_price": 50000, "currency": "IDR" },
    "dedicated": { "base_price": 15000000, "per_user_price": 40000, "currency": "IDR" }
  }
}
```

### POST /api/v1/licenses — Buat lisensi baru + auto-generate proposal

```json
// Request — product_id references products table
{
  "client_name": "PT Maju Bersama",
  "client_email": "admin@majubersama.com",
  "product_id": "uuid-of-product",
  "plan": "saas",
  "status": "trial",
  "modules": ["accounting", "sales"],
  "apps": ["web-ui", "app-sales-person"],
  "contract_amount": 15000000,
  "max_users": 50,
  "max_trans_per_month": 10000,
  "max_trans_per_day": 500,
  "max_items": 5000,
  "max_customers": 1000,
  "max_branches": 3,
  "expires_at": "2027-01-01T00:00:00Z"
}

// Response (201)
{
  "data": {
    "id": "...",
    "license_key": "FL-XXXXXXXX",
    "initial_password": "generated-password",
    "proposal_url": "/api/v1/licenses/{id}/proposal",
    ...
  }
}
```

**Side effects saat create:**
1. Generate `license_key` (FL-XXXXXXXX)
2. Generate `initial_password` (hanya muncul sekali di response ini)
3. **Auto-generate proposal PDF** (server-side, stored di filesystem/S3)
4. Insert audit log `license_created`

### GET /api/v1/licenses/{id}/proposal — Download proposal PDF

```
GET /api/v1/licenses/{id}/proposal
Response: application/pdf (binary stream)
Header: Content-Disposition: attachment; filename="proposal-FL-XXXXXXXX.pdf"
```

Returns 404 jika proposal belum di-generate.

### POST /api/v1/licenses/{id}/proposal/regenerate — Regenerate proposal

Gunakan jika data lisensi berubah (constraints, modules, etc.) dan proposal perlu di-update.

```json
// Response
{
  "data": {
    "proposal_url": "/api/v1/licenses/{id}/proposal",
    "generated_at": "2026-03-24T10:00:00Z"
  }
}
```

### POST /api/v1/licenses/{id}/renew — Perpanjang lisensi

```json
// Request
{
  "new_expires_at": "2028-01-01T00:00:00Z",
  "contract_amount": 50000000,
  "notes": "Renewal for 2028"
}
```

Validasi:
- `new_expires_at` harus di masa depan
- `new_expires_at` harus setelah `expires_at` saat ini
- Status harus: `active`, `expired`, atau `trial`
- Jika `expired` → otomatis transisi ke `active`

### POST /api/v1/licenses/bulk/status — Bulk update

```json
// Request
{
  "license_ids": ["uuid-1", "uuid-2", "uuid-3"],
  "new_status": "suspended",
  "reason": "Non-payment batch"
}

// Response (207 jika partial failure)
{
  "success_count": 2,
  "failed": [
    { "id": "uuid-3", "error": "Invalid transition: trial → suspended" }
  ]
}
```

### GET /api/v1/client/license — Public license check

```
GET /api/v1/client/license?key=FL-XXXXXXXX
Header: Origin: https://mycompany.flasherp.id
```

**Validasi:**
1. Ambil lisensi by `key`
2. Cek `Origin` (atau `Referer` jika kosong) vs `flasherp_url`
3. Bandingkan **host** setelah normalisasi (case-insensitive)
4. Tidak cocok → `403 Forbidden`

**Response (safe — tanpa `client_registration_code`):**
```json
{
  "data": {
    "id": "...",
    "license_key": "FL-XXXXXXXX",
    "client_name": "PT Maju Bersama",
    "product": "flasherp",
    "plan": "saas",
    "status": "active",
    "modules": ["accounting", "sales"],
    "apps": ["web-ui", "app-sales-person"],
    "max_users": 50,
    "expires_at": "2027-01-01T00:00:00Z",
    "is_provisioned": true
  }
}
```

### GET /api/v1/dashboard — Analytics

```json
{
  "summary": {
    "total_licenses": 142,
    "active": 98,
    "trial": 23,
    "suspended": 8,
    "expired": 13
  },
  "revenue": {
    "total_mrr": 285000000,
    "avg_contract": 2908163,
    "by_product": {
      "flasherp": 210000000,
      "flashpos": 45000000,
      "flashhrd": 20000000,
      "flashaccounting": 10000000
    }
  },
  "expiring_soon": [
    { "license_key": "FL-XXXXXXXX", "client_name": "...", "expires_at": "..." }
  ],
  "recent_activity": [
    { "action": "status_changed", "actor_name": "...", "created_at": "..." }
  ],
  "distribution": {
    "by_plan": { "saas": 120, "dedicated": 22 },
    "by_product": { "flasherp": 95, "flashpos": 28, "flashhrd": 12, "flashaccounting": 7 }
  }
}
```

Query handler (`internal/query/get_dashboard/`): single SQL dengan CTEs, target < 100ms.
