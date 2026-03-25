# Auth & Error Codes

## JWT (Vernon App only)

```json
{ "sub": "<uuid>", "name": "John", "role": "sales", "iat": 1711234567, "exp": 1711320967 }
```

- Algorithm: HS256, secret: `JWT_SECRET`
- Lifetime: 24h, no refresh token
- Stored: `localStorage` (go-app)
- Header: `Authorization: Bearer <token>`

## Auth Flow

```
App load → cek localStorage token
  valid   → lanjut ke halaman
  expired → Clear() → navigate /login
  missing → navigate /login

Login  → POST /api/internal/auth/login { email, password }
       → response: { token } → simpan di localStorage

API 401 → authStore.Clear() → navigate /login (lihat ui/api/client.go ErrUnauthorized)
```

## Setup Flow (first run)

```
GET /api/internal/setup/status → { is_setup: false }
POST /api/internal/setup/install { name, email, password } → creates superuser
```

---

## Error Codes

### Public API
| Code | HTTP |
|---|---|
| `INVALID_CLIENT_CODE` | 403 |
| `PRODUCT_NOT_FOUND` | 403 |
| `ALREADY_REGISTERED` | 409 |
| `LICENSE_NOT_FOUND` | 404 |
| `RATE_LIMIT_EXCEEDED` | 429 |
| `INTERNAL_ERROR` | 500 |

### Internal API
| Code | HTTP |
|---|---|
| `AUTH_INVALID_CREDENTIALS` | 401 |
| `UNAUTHORIZED` | 401 |
| `FORBIDDEN` | 403 |
| `LICENSE_NOT_FOUND` | 404 |
| `LICENSE_INVALID_TRANSITION` | 422 |
| `COMPANY_NOT_FOUND` | 404 |
| `PROJECT_NOT_FOUND` | 404 |
| `PROPOSAL_NOT_FOUND` | 404 |
| `PROPOSAL_NOT_SUBMITTED` | 422 |
| `PROPOSAL_ACTIVE_EXISTS` | 409 |
| `PRODUCT_SLUG_EXISTS` | 409 |
| `USER_EMAIL_EXISTS` | 409 |
| `VALIDATION_FAILED` | 400 |
| `INTERNAL_ERROR` | 500 |
