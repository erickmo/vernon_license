# Vernon License — API Documentation

**Base URL:** `http://localhost:8081` (default, sesuaikan dengan `PORT` env)

---

## Daftar Isi

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Error Format](#error-format)
4. [Public API](#public-api)
   - [POST /api/v1/register](#post-apiv1register)
   - [GET /api/v1/validate](#get-apiv1validate)
5. [Internal API — Setup](#internal-api--setup)
   - [GET /api/internal/setup/status](#get-apiinternalsetupstatus)
   - [POST /api/internal/setup/install](#post-apiinternalsetupinstall)
6. [Internal API — Auth](#internal-api--auth)
   - [POST /api/internal/auth/login](#post-apiinternalauthlogin)
   - [GET /api/internal/auth/me](#get-apiinternalauthme)
7. [Internal API — Dashboard](#internal-api--dashboard)
8. [Internal API — Companies](#internal-api--companies)
9. [Internal API — Projects](#internal-api--projects)
10. [Internal API — Licenses](#internal-api--licenses)
11. [Internal API — Proposals](#internal-api--proposals)
12. [Internal API — Products](#internal-api--products)
13. [Internal API — Users](#internal-api--users)
14. [Internal API — Notifications](#internal-api--notifications)
15. [Role Matrix](#role-matrix)

---

## Overview

Vernon License menyediakan dua jenis API:

| Jenis | Base Path | Auth | Tujuan |
|---|---|---|---|
| **Public API** | `/api/v1/` | Tidak ada | Dipanggil oleh client apps (register + validate license) |
| **Internal API** | `/api/internal/` | JWT Bearer | Digunakan oleh Vernon App (WASM) untuk manajemen |

Rate limit Public API: **60 req/menit per IP**.

---

## Authentication

Internal API menggunakan **JWT HS256**. Token didapat dari `POST /api/internal/auth/login`.

**Request Header:**
```
Authorization: Bearer <token>
```

**JWT Claims:**
```json
{
  "sub":  "uuid-user-id",
  "name": "John Doe",
  "role": "superuser",
  "iat":  1711234567,
  "exp":  1711320967
}
```

Token berlaku selama **24 jam**. Setelah expired, login ulang diperlukan.

---

## Error Format

Semua error menggunakan format JSON yang konsisten:

```json
{
  "error": {
    "code":    "ERROR_CODE",
    "message": "Human-readable message"
  }
}
```

**HTTP Status Codes yang Digunakan:**

| Status | Arti |
|---|---|
| `200 OK` | Sukses |
| `201 Created` | Resource berhasil dibuat |
| `400 Bad Request` | Validasi gagal |
| `401 Unauthorized` | Token tidak ada atau tidak valid |
| `403 Forbidden` | Tidak punya akses |
| `404 Not Found` | Resource tidak ditemukan |
| `409 Conflict` | Duplikasi (misal: email sudah terdaftar) |
| `422 Unprocessable Entity` | Transisi status tidak valid |
| `500 Internal Server Error` | Kesalahan server |

**Error Codes:**

| Code | Deskripsi |
|---|---|
| `VALIDATION_FAILED` | Field wajib kosong atau format salah |
| `UNAUTHORIZED` | Token tidak ada, kadaluarsa, atau tidak valid |
| `FORBIDDEN` | Role tidak memiliki akses |
| `NOT_FOUND` | Resource tidak ditemukan |
| `LICENSE_NOT_FOUND` | License key tidak ditemukan |
| `INVALID_API_KEY` | `client_registration_code` tidak valid |
| `ALREADY_REGISTERED` | License sudah pernah di-register |
| `INVALID_CREDENTIALS` | Email atau password salah |
| `USER_NOT_FOUND` | User tidak ditemukan |
| `USER_EMAIL_EXISTS` | Email sudah digunakan |
| `PRODUCT_NOT_FOUND` | Product tidak ditemukan |
| `PRODUCT_INACTIVE` | Product tidak aktif |
| `PROPOSAL_NOT_FOUND` | Proposal tidak ditemukan |
| `PROPOSAL_NOT_DRAFT` | Proposal bukan dalam status draft |
| `PROPOSAL_NOT_SUBMITTED` | Proposal bukan dalam status submitted |
| `INVALID_TRANSITION` | Transisi status license tidak diizinkan |
| `INTERNAL_ERROR` | Kesalahan internal server |

---

## Public API

### POST /api/v1/register

Mendaftarkan client app (instance) ke Vernon License. Dipanggil **sekali** saat client app pertama kali di-install.

**Auth:** Tidak diperlukan
**Rate limit:** 60 req/menit per IP

**Request Body:**
```json
{
  "product_slug":       "flasherp",
  "instance_url":       "https://mycompany.flasherp.id",
  "instance_name":      "PT Maju Bersama",
  "client_registration_code":  "prov_xxxxxxxxxxxx"
}
```

| Field | Tipe | Wajib | Deskripsi |
|---|---|---|---|
| `product_slug` | string | Ya | Slug produk (contoh: `flasherp`) |
| `instance_url` | string | Ya | URL instance client |
| `instance_name` | string | Ya | Nama instance / perusahaan |
| `client_registration_code` | string | Ya | Kunci provisioning dari Vernon App |

**Response 201 Created:**
```json
{
  "license_key":     "FL-XXXXXXXX",
  "product":         "FlashERP",
  "check_interval":  "6h",
  "valid":           true,
  "message":         "License registered successfully"
}
```

**Error Responses:**

| Status | Code | Kondisi |
|---|---|---|
| `400` | `VALIDATION_FAILED` | Field kosong |
| `403` | `INVALID_API_KEY` | `client_registration_code` tidak cocok dengan `product_slug` |
| `409` | `ALREADY_REGISTERED` | Instance sudah pernah register |

---

### GET /api/v1/validate

Memvalidasi status license. Dipanggil secara periodik oleh client app (interval sesuai `check_interval`).

**Auth:** Tidak diperlukan
**Rate limit:** 60 req/menit per IP

**Query Parameters:**

| Param | Tipe | Wajib | Deskripsi |
|---|---|---|---|
| `key` | string | Ya | License key (contoh: `FL-XXXXXXXX`) |

**Contoh Request:**
```
GET /api/v1/validate?key=FL-XXXXXXXX
```

**Response 200 — License Valid:**
```json
{
  "valid":          true,
  "license_key":    "FL-XXXXXXXX",
  "check_interval": "6h"
}
```

**Response 200 — License Tidak Valid:**
```json
{
  "valid":          false,
  "license_key":    "FL-XXXXXXXX",
  "reason":         "license suspended",
  "check_interval": "6h"
}
```

> Endpoint ini selalu mengembalikan `200`. Client app harus memeriksa field `valid`.

**Error Responses:**

| Status | Code | Kondisi |
|---|---|---|
| `400` | `VALIDATION_FAILED` | Query param `key` tidak ada |
| `404` | `LICENSE_NOT_FOUND` | License key tidak dikenali |

---

## Internal API — Setup

### GET /api/internal/setup/status

Mengecek apakah system sudah di-setup (sudah ada superuser). Digunakan untuk redirect ke `/setup` jika belum ada user.

**Auth:** Tidak diperlukan

**Response 200:**
```json
{
  "is_setup": true
}
```

---

### POST /api/internal/setup/install

First-run setup: membuat akun superuser pertama. Hanya bisa dipanggil sekali (sebelum ada user di database).

**Auth:** Tidak diperlukan

**Request Body:**
```json
{
  "name":     "Admin",
  "email":    "admin@company.com",
  "password": "securepassword"
}
```

**Response 200:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id":    "uuid",
    "name":  "Admin",
    "role":  "superuser",
    "email": "admin@company.com"
  }
}
```

---

## Internal API — Auth

### POST /api/internal/auth/login

Login ke Vernon App. Mengembalikan JWT token.

**Auth:** Tidak diperlukan

**Request Body:**
```json
{
  "email":    "admin@company.com",
  "password": "securepassword"
}
```

**Response 200:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id":    "uuid",
    "name":  "Admin",
    "role":  "superuser",
    "email": "admin@company.com"
  }
}
```

**Error Responses:**

| Status | Code | Kondisi |
|---|---|---|
| `400` | `VALIDATION_FAILED` | Email atau password kosong |
| `401` | `INVALID_CREDENTIALS` | Email atau password salah |

---

### GET /api/internal/auth/me

Mendapatkan data user yang sedang login berdasarkan JWT token.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
{
  "id":    "uuid",
  "name":  "Admin",
  "role":  "superuser",
  "email": "admin@company.com"
}
```

---

## Internal API — Dashboard

### GET /api/internal/dashboard

Mendapatkan statistik agregasi untuk Vernon App dashboard.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
{
  "total_licenses":     142,
  "active_licenses":    98,
  "pending_licenses":   23,
  "suspended_licenses": 8,
  "expired_licenses":   13,
  "total_companies":    45,
  "total_proposals":    67,
  "pending_proposals":  12,
  "total_revenue":      285000000,
  "expiring_licenses": [
    {
      "id":          "uuid",
      "license_key": "FL-XXXXXXXX",
      "company":     "PT Maju Bersama",
      "expires_at":  "2026-04-15",
      "days_left":   21
    }
  ],
  "recent_activity": [
    {
      "entity_type": "license",
      "action":      "status_changed",
      "actor_name":  "John Doe",
      "created_at":  "2026-03-25T10:30:00Z"
    }
  ]
}
```

> Semua query bersifat non-fatal — jika ada query yang gagal, nilai field tersebut akan `0` dan endpoint tetap mengembalikan `200`.

---

## Internal API — Companies

### GET /api/internal/companies

List semua companies.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
[
  {
    "id":      "uuid",
    "name":    "PT Maju Bersama",
    "email":   "info@maju.com",
    "phone":   "+62-21-1234567",
    "address": "Jl. Teknologi No. 1, Jakarta"
  }
]
```

---

### POST /api/internal/companies

Membuat company baru.

**Auth:** Bearer JWT (semua role)

**Request Body:**
```json
{
  "name":    "PT Maju Bersama",
  "email":   "info@maju.com",
  "phone":   "+62-21-1234567",
  "address": "Jl. Teknologi No. 1, Jakarta"
}
```

**Response 201:** Company object yang dibuat.

---

### GET /api/internal/companies/{id}

Mendapatkan detail company.

**Auth:** Bearer JWT (semua role)

**Response 200:** Company object lengkap.

---

### PUT /api/internal/companies/{id}

Update data company.

**Auth:** Bearer JWT (semua role)

**Request Body:** Sama dengan POST, semua field opsional (partial update).

**Response 200:** Company object yang diupdate.

---

### DELETE /api/internal/companies/{id}

Soft delete company.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
{ "result": "ok" }
```

---

## Internal API — Projects

Projects selalu berada di bawah Company.

### GET /api/internal/companies/{companyID}/projects

List semua projects milik company tertentu.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
[
  {
    "id":         "uuid",
    "company_id": "uuid",
    "name":       "FlashERP Implementation",
    "description": "Implementasi FlashERP untuk PT Maju Bersama"
  }
]
```

---

### POST /api/internal/companies/{companyID}/projects

Membuat project baru di bawah company.

**Auth:** Bearer JWT (semua role)

**Request Body:**
```json
{
  "name":        "FlashERP Implementation",
  "description": "Implementasi FlashERP untuk PT Maju Bersama"
}
```

**Response 201:** Project object yang dibuat.

---

### GET /api/internal/projects/{id}

Mendapatkan detail project.

**Auth:** Bearer JWT (semua role)

---

### PUT /api/internal/projects/{id}

Update data project.

**Auth:** Bearer JWT (semua role)

---

### DELETE /api/internal/projects/{id}

Soft delete project.

**Auth:** Bearer JWT (semua role)

---

## Internal API — Licenses

### GET /api/internal/licenses

List semua licenses dengan info company, project, dan product.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
[
  {
    "id":            "uuid",
    "license_key":   "FL-XXXXXXXX",
    "company_name":  "PT Maju Bersama",
    "project_name":  "FlashERP Implementation",
    "product_name":  "FlashERP",
    "plan":          "saas",
    "status":        "active",
    "is_registered": true,
    "expires_at":    "2027-01-01T00:00:00Z"
  }
]
```

---

### GET /api/internal/projects/{projectID}/licenses

List licenses milik project tertentu.

**Auth:** Bearer JWT (semua role)

**Response 200:** Array license list item (sama dengan GET /licenses).

---

### POST /api/internal/licenses

Membuat license baru secara langsung (direct create). Menghasilkan `license_key` dan `client_registration_code` otomatis.

**Auth:** Bearer JWT — `project_owner` atau `superuser` saja

**Request Body:**
```json
{
  "project_id":           "uuid",
  "company_id":           "uuid",
  "product_id":           "uuid",
  "plan":                 "saas",
  "modules":              ["accounting", "sales"],
  "apps":                 ["web-ui", "app-sales-person"],
  "contract_amount":      15000000,
  "description":          "Lisensi FlashERP untuk PT Maju",
  "max_users":            50,
  "max_trans_per_month":  10000,
  "max_trans_per_day":    500,
  "max_items":            5000,
  "max_customers":        1000,
  "max_branches":         3,
  "max_storage":          10240,
  "expires_at":           "2027-01-01T00:00:00Z",
  "check_interval":       "6h"
}
```

| Field | Tipe | Wajib | Deskripsi |
|---|---|---|---|
| `project_id` | uuid string | Ya | ID project |
| `company_id` | uuid string | Ya | ID company |
| `product_id` | uuid string | Ya | ID product |
| `plan` | string | Ya | `saas` atau `dedicated` |
| `modules` | []string | Tidak | Modul yang diaktifkan |
| `apps` | []string | Tidak | App yang diaktifkan |
| `contract_amount` | float64 | Tidak | Nilai kontrak (IDR) |
| `description` | string | Tidak | Deskripsi lisensi |
| `max_users` | int | Tidak | Batas jumlah user |
| `max_trans_per_month` | int | Tidak | Batas transaksi per bulan |
| `max_trans_per_day` | int | Tidak | Batas transaksi per hari |
| `max_items` | int | Tidak | Batas jumlah item |
| `max_customers` | int | Tidak | Batas jumlah pelanggan |
| `max_branches` | int | Tidak | Batas jumlah cabang |
| `max_storage` | int | Tidak | Batas storage (MB) |
| `expires_at` | string | Tidak | RFC3339 atau `YYYY-MM-DD` |
| `check_interval` | string | Tidak | Contoh: `6h`, `24h`, `168h` |

**Response 201:** License detail object (lihat di bawah).

---

### GET /api/internal/licenses/{id}

Mendapatkan detail license lengkap termasuk `client_registration_code`.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
{
  "id":                  "uuid",
  "license_key":         "FL-XXXXXXXX",
  "company_id":          "uuid",
  "company_name":        "PT Maju Bersama",
  "project_id":          "uuid",
  "project_name":        "FlashERP Implementation",
  "product_name":        "FlashERP",
  "plan":                "saas",
  "status":              "active",
  "modules":             ["accounting", "sales"],
  "apps":                ["web-ui"],
  "contract_amount":     15000000,
  "description":         "Lisensi FlashERP untuk PT Maju",
  "max_users":           50,
  "max_trans_per_month": 10000,
  "max_trans_per_day":   500,
  "max_items":           5000,
  "max_customers":       1000,
  "max_branches":        3,
  "max_storage":         10240,
  "expires_at":          "2027-01-01T00:00:00Z",
  "is_registered":       true,
  "instance_url":        "https://mycompany.flasherp.id",
  "instance_name":       "PT Maju Bersama",
  "client_registration_code":   "prov_xxxxxxxxxxxx",
  "check_interval":      "6h",
  "last_pull_at":        "2026-03-25T10:00:00Z"
}
```

---

### PUT /api/internal/licenses/{id}/activate

Mengubah status license ke `active`.

**Auth:** Bearer JWT — `project_owner` atau `superuser` saja

**Transisi yang diizinkan:** `pending` → `active`, `suspended` → `active`

**Response 200:**
```json
{ "status": "active" }
```

---

### PUT /api/internal/licenses/{id}/suspend

Mengubah status license ke `suspended`.

**Auth:** Bearer JWT — `project_owner` atau `superuser` saja

**Transisi yang diizinkan:** `active` → `suspended`

**Response 200:**
```json
{ "status": "suspended" }
```

---

### PUT /api/internal/licenses/{id}/renew

Memperbarui tanggal expired license. Jika license sedang `expired`, otomatis diaktifkan kembali.

**Auth:** Bearer JWT — `project_owner` atau `superuser` saja

**Request Body:**
```json
{
  "new_expires_at": "2028-01-01T00:00:00Z"
}
```

Format `new_expires_at`: RFC3339 atau `YYYY-MM-DD`. Jika kosong, hanya status yang diubah (jika sedang expired).

**Response 200:**
```json
{ "status": "active" }
```

---

### PUT /api/internal/licenses/{id}/constraints

Memperbarui constraints dan konfigurasi license.

**Auth:** Bearer JWT — `project_owner` atau `superuser` saja

**Request Body:**
```json
{
  "modules":             ["accounting", "sales", "inventory"],
  "apps":                ["web-ui", "app-mobile"],
  "max_users":           100,
  "max_trans_per_month": 20000,
  "max_trans_per_day":   1000,
  "max_items":           10000,
  "max_customers":       2000,
  "max_branches":        5,
  "max_storage":         20480,
  "expires_at":          "2028-01-01T00:00:00Z",
  "check_interval":      "12h"
}
```

Semua field opsional. Hanya field yang dikirim yang akan diupdate.

**Response 200:**
```json
{ "result": "ok" }
```

---

### GET /api/internal/licenses/{id}/audit

Mendapatkan audit trail untuk license tertentu.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
[
  {
    "id":         "uuid",
    "action":     "license_created",
    "actor_id":   "uuid",
    "actor_name": "John Doe",
    "changes":    { "status": { "from": null, "to": "active" } },
    "metadata":   {},
    "created_at": "2026-03-25T10:00:00Z"
  }
]
```

**Action Values:**
- `license_created` — License dibuat
- `status_changed` — Status diubah (activate/suspend/renew)
- `constraints_updated` — Constraints diupdate
- `client_registered` — Client app melakukan register

---

## Internal API — Proposals

Proposals adalah dokumen penawaran yang dibuat oleh Sales, lalu disetujui/ditolak oleh Project Owner. Saat disetujui, license dibuat otomatis.

**Status Flow:**
```
draft → submitted → approved (license dibuat)
                 → rejected → draft (bisa diedit ulang)
```

---

### GET /api/internal/proposals

List semua proposals.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
{
  "data": [
    {
      "id":         "uuid",
      "project_id": "uuid",
      "company_id": "uuid",
      "product_id": "uuid",
      "version":    1,
      "status":     "submitted",
      "plan":       "saas",
      "created_at": "2026-03-20T08:00:00Z",
      "updated_at": "2026-03-22T10:00:00Z"
    }
  ]
}
```

---

### GET /api/internal/projects/{projectID}/proposals

List proposals milik project tertentu.

**Auth:** Bearer JWT (semua role)

**Response 200:** Sama dengan GET /proposals.

---

### POST /api/internal/proposals

Membuat proposal baru dengan status `draft`.

**Auth:** Bearer JWT (semua role)

**Request Body:**
```json
{
  "project_id":          "uuid",
  "company_id":          "uuid",
  "product_id":          "uuid",
  "plan":                "saas",
  "modules":             ["accounting", "sales"],
  "apps":                ["web-ui"],
  "max_users":           50,
  "max_trans_per_month": 10000,
  "max_trans_per_day":   500,
  "max_items":           5000,
  "max_customers":       1000,
  "max_branches":        3,
  "max_storage":         10240,
  "contract_amount":     15000000,
  "expires_at":          "2027-01-01T00:00:00Z",
  "notes":               "Proposal untuk implementasi tahap 1"
}
```

| Field | Tipe | Wajib |
|---|---|---|
| `project_id` | uuid string | Ya |
| `company_id` | uuid string | Ya |
| `product_id` | uuid string | Ya |
| `plan` | string | Ya |
| Semua lainnya | — | Tidak |

**Response 201:** Proposal list item object.

---

### GET /api/internal/proposals/{id}

Mendapatkan detail proposal lengkap beserta changelog dan nama entitas terkait.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
{
  "id":                  "uuid",
  "project_id":          "uuid",
  "project_name":        "FlashERP Implementation",
  "company_id":          "uuid",
  "company_name":        "PT Maju Bersama",
  "product_id":          "uuid",
  "product_name":        "FlashERP",
  "version":             2,
  "status":              "submitted",
  "plan":                "saas",
  "modules":             ["accounting", "sales"],
  "apps":                ["web-ui"],
  "max_users":           50,
  "max_trans_per_month": 10000,
  "max_trans_per_day":   500,
  "max_items":           5000,
  "max_customers":       1000,
  "max_branches":        3,
  "max_storage":         10240,
  "contract_amount":     15000000,
  "expires_at":          "2027-01-01T00:00:00Z",
  "notes":               "Proposal untuk implementasi tahap 1",
  "owner_notes":         "Disetujui dengan syarat...",
  "rejection_reason":    "",
  "changelog": {
    "compared_to_version": 1,
    "summary":             "2 perubahan dari versi sebelumnya",
    "changes": [
      {
        "field":     "max_users",
        "old_value": 30,
        "new_value": 50
      }
    ],
    "unchanged": ["plan", "modules", "contract_amount"]
  },
  "pdf_path":          "proposals/proposal-FL-XXXXXXXX-v2.pdf",
  "submitted_by_name": "Jane Sales",
  "reviewed_by_name":  "Bob Owner",
  "reviewed_at":       "2026-03-23T14:00:00Z",
  "created_at":        "2026-03-20T08:00:00Z",
  "updated_at":        "2026-03-23T14:00:00Z"
}
```

---

### PUT /api/internal/proposals/{id}

Update proposal. Sales hanya bisa edit `draft`. Project Owner / Superuser bisa edit `draft` atau `submitted`.

**Auth:** Bearer JWT (semua role, dengan batasan di atas)

**Request Body:**
```json
{
  "plan":                "dedicated",
  "modules":             ["accounting", "sales", "inventory"],
  "apps":                ["web-ui", "app-mobile"],
  "max_users":           100,
  "max_trans_per_month": 20000,
  "max_trans_per_day":   1000,
  "max_items":           10000,
  "max_customers":       2000,
  "max_branches":        5,
  "max_storage":         20480,
  "contract_amount":     25000000,
  "expires_at":          "2028-01-01T00:00:00Z",
  "notes":               "Updated notes",
  "owner_notes":         "Catatan dari PO"
}
```

Semua field opsional (partial update). `owner_notes` hanya bisa diisi oleh `project_owner` / `superuser`.

**Response 200:** Proposal list item object.

**Error:**

| Status | Code | Kondisi |
|---|---|---|
| `422` | `PROPOSAL_NOT_DRAFT` | Sales mencoba edit non-draft proposal |

---

### PUT /api/internal/proposals/{id}/submit

Mengubah status proposal dari `draft` → `submitted`.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
{ "message": "Proposal berhasil di-submit" }
```

**Error:**

| Status | Code | Kondisi |
|---|---|---|
| `422` | `PROPOSAL_NOT_DRAFT` | Proposal bukan dalam status draft |

---

### PUT /api/internal/proposals/{id}/approve

Menyetujui proposal. Otomatis membuat license baru.

**Auth:** Bearer JWT — `project_owner` atau `superuser` saja

**Request Body (opsional):**
```json
{
  "owner_notes": "Disetujui. Pastikan onboarding selesai sebelum go-live."
}
```

**Response 200:**
```json
{
  "message":    "Proposal disetujui dan lisensi telah dibuat",
  "license_id": "uuid"
}
```

**Error:**

| Status | Code | Kondisi |
|---|---|---|
| `403` | `FORBIDDEN` | Role bukan project_owner atau superuser |
| `422` | `PROPOSAL_NOT_SUBMITTED` | Proposal bukan dalam status submitted |

---

### PUT /api/internal/proposals/{id}/reject

Menolak proposal.

**Auth:** Bearer JWT — `project_owner` atau `superuser` saja

**Request Body:**
```json
{
  "reason": "Budget tidak sesuai, perlu revisi contract amount"
}
```

`reason` wajib diisi.

**Response 200:**
```json
{ "message": "Proposal ditolak" }
```

**Error:**

| Status | Code | Kondisi |
|---|---|---|
| `403` | `FORBIDDEN` | Role bukan project_owner atau superuser |
| `422` | `PROPOSAL_NOT_SUBMITTED` | Proposal bukan dalam status submitted |

---

### GET /api/internal/proposals/{id}/pdf

Download PDF proposal.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```
Content-Type: application/pdf
Content-Disposition: attachment; filename="proposal-FL-XXXXXXXX-v1.pdf"

(binary PDF content)
```

**Error:**

| Status | Code | Kondisi |
|---|---|---|
| `404` | `PDF_NOT_FOUND` | PDF belum tersedia atau file tidak ada |

---

## Internal API — Products

### GET /api/internal/products

List semua products (termasuk yang tidak aktif).

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
[
  {
    "id":          "uuid",
    "name":        "FlashERP",
    "slug":        "flasherp",
    "description": "Sistem ERP terintegrasi",
    "is_active":   true
  }
]
```

---

### POST /api/internal/products

Membuat product baru.

**Auth:** Bearer JWT — `project_owner` atau `superuser` saja

**Request Body:**
```json
{
  "name":        "FlashHRD",
  "slug":        "flashhrd",
  "description": "Sistem manajemen SDM",
  "is_active":   true
}
```

**Response 201:** Product object yang dibuat.

---

### GET /api/internal/products/{id}

Detail product.

**Auth:** Bearer JWT (semua role)

---

### PUT /api/internal/products/{id}

Update product.

**Auth:** Bearer JWT — `project_owner` atau `superuser` saja

---

### DELETE /api/internal/products/{id}

Soft delete product.

**Auth:** Bearer JWT — `project_owner` atau `superuser` saja

---

## Internal API — Users

### GET /api/internal/users

List semua users.

**Auth:** Bearer JWT — `superuser` saja

**Response 200:**
```json
{
  "data": [
    {
      "id":        "uuid",
      "name":      "John Doe",
      "email":     "john@company.com",
      "role":      "sales",
      "is_active": true
    }
  ]
}
```

---

### POST /api/internal/users

Membuat user baru.

**Auth:** Bearer JWT — `superuser` saja

**Request Body:**
```json
{
  "name":     "Jane Sales",
  "email":    "jane@company.com",
  "password": "securepassword",
  "role":     "sales"
}
```

| Field | Tipe | Wajib | Nilai |
|---|---|---|---|
| `name` | string | Ya | — |
| `email` | string | Ya | — |
| `password` | string | Ya | — |
| `role` | string | Ya | `project_owner` atau `sales` |

**Response 201:** User object.

**Error:**

| Status | Code | Kondisi |
|---|---|---|
| `409` | `USER_EMAIL_EXISTS` | Email sudah terdaftar |

---

### PUT /api/internal/users/{id}/active

Mengaktifkan atau menonaktifkan user. **Superuser tidak bisa dinonaktifkan.**

**Auth:** Bearer JWT — `superuser` saja

**Request Body:**
```json
{
  "is_active": false
}
```

**Response 200:**
```json
{ "is_active": false }
```

**Error:**

| Status | Code | Kondisi |
|---|---|---|
| `403` | `FORBIDDEN` | Target user adalah superuser |
| `404` | `USER_NOT_FOUND` | User tidak ditemukan |

---

## Internal API — Notifications

### GET /api/internal/notifications

List notifikasi untuk user yang sedang login.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
[
  {
    "id":         "uuid",
    "title":      "Proposal Disetujui",
    "body":       "Proposal untuk PT Maju Bersama telah disetujui",
    "is_read":    false,
    "created_at": "2026-03-25T10:00:00Z"
  }
]
```

---

### GET /api/internal/notifications/unread-count

Mendapatkan jumlah notifikasi yang belum dibaca (untuk badge UI).

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
{ "count": 5 }
```

---

### PUT /api/internal/notifications/{id}/read

Menandai notifikasi sebagai sudah dibaca.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
{ "result": "ok" }
```

---

### PUT /api/internal/notifications/read-all

Menandai semua notifikasi sebagai sudah dibaca.

**Auth:** Bearer JWT (semua role)

**Response 200:**
```json
{ "result": "ok" }
```

---

## Role Matrix

| Endpoint | `sales` | `project_owner` | `superuser` |
|---|:---:|:---:|:---:|
| Login, Me, Dashboard | ✓ | ✓ | ✓ |
| Companies CRUD | ✓ | ✓ | ✓ |
| Projects CRUD | ✓ | ✓ | ✓ |
| Licenses — List, Detail, Audit | ✓ | ✓ | ✓ |
| Licenses — Create | — | ✓ | ✓ |
| Licenses — Activate, Suspend, Renew | — | ✓ | ✓ |
| Licenses — Update Constraints | — | ✓ | ✓ |
| Proposals — List, Create, View | ✓ | ✓ | ✓ |
| Proposals — Edit (draft only) | ✓ | ✓ | ✓ |
| Proposals — Edit (submitted) | — | ✓ | ✓ |
| Proposals — Submit | ✓ | ✓ | ✓ |
| Proposals — Approve, Reject | — | ✓ | ✓ |
| Products — List, Detail | ✓ | ✓ | ✓ |
| Products — Create, Edit, Delete | — | ✓ | ✓ |
| Users — List, Create, SetActive | — | — | ✓ |
| Notifications | ✓ | ✓ | ✓ |
