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

## License Registration & Superuser Flow

```
1. Client App registers:
   POST /api/v1/register { app_name, otp, client_name, instance_url }
   → License app captures client_app_ip, validates, creates license (pending)
   → Returns license_key
   → Client stores: license_key, license_app_ip, instance_url in .env

2. License App admin approves (activate/trial):
   Admin inputs username + password in UI
   → License app fetches active OTP from otp table
   → License app calls POST {instance_url}/api/v1/create-superuser
     Body: { otp, license_key, username, password }
   → Client app validates:
     - license_key matches .env
     - sender IP = license_app_ip in .env
     - OTP via POST {license_app_url}/api/v1/validate_otp
   → Client creates superuser, returns { username }
   → License app stores superuser_username, updates status

3. Reset Superuser (anytime):
   PUT /api/internal/licenses/{id}/reset-superuser { username, password }
   → Same callback flow as step 2, no status change
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
| `SUPERUSER_CREATION_FAILED` | 502 |
| `NO_ACTIVE_OTP` | 400 |
| `NO_INSTANCE_URL` | 400 |
| `VALIDATION_FAILED` | 400 |
| `INTERNAL_ERROR` | 500 |
