# Error Codes

## Public API Errors (register + validate)

| Code | HTTP | Message |
|---|---|---|
| `INVALID_API_KEY` | 403 | Invalid provision API key |
| `ALREADY_REGISTERED` | 409 | Instance URL already registered |
| `PRODUCT_NOT_FOUND` | 404 | Unknown product slug |
| `PRODUCT_INACTIVE` | 422 | Product is not active |
| `LICENSE_NOT_FOUND` | 404 | Unknown license key |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |

## Internal App Errors (Vernon App — not exposed publicly)

| Code | Message (ID) |
|---|---|
| `AUTH_INVALID_CREDENTIALS` | Email atau password salah |
| `AUTH_TOKEN_EXPIRED` | Sesi habis |
| `AUTH_INSUFFICIENT_ROLE` | Tidak memiliki akses |
| `COMPANY_NOT_FOUND` | Company tidak ditemukan |
| `PROJECT_NOT_FOUND` | Project tidak ditemukan |
| `LICENSE_INVALID_TRANSITION` | Perubahan status tidak valid |
| `PROPOSAL_NOT_FOUND` | Proposal tidak ditemukan |
| `PROPOSAL_NOT_APPROVED` | Proposal belum disetujui |
| `PROPOSAL_NOT_DRAFT` | Hanya draft yang bisa diedit sales |
| `PROPOSAL_NOT_SUBMITTED` | Hanya submitted yang bisa di-review |
| `PROPOSAL_ACTIVE_EXISTS` | Masih ada proposal aktif |
| `PROPOSAL_GENERATION_FAILED` | Gagal membuat PDF |
| `PRODUCT_SLUG_EXISTS` | Slug sudah digunakan |
| `PRODUCT_INVALID_MODULE` | Modul tidak tersedia |
| `USER_EMAIL_EXISTS` | Email sudah digunakan |
| `VALIDATION_FAILED` | Data tidak valid |
