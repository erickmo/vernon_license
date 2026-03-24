# Auth & Flows

## JWT Payload

```json
{
  "sub": "<user-uuid>",
  "name": "John Doe",
  "role": "sales",
  "iat": 1711234567,
  "exp": 1711320967
}
```

| Claim | Type | Deskripsi |
|---|---|---|
| `sub` | UUID string | User ID |
| `name` | string | Nama lengkap user |
| `role` | string | `superuser` \| `project_owner` \| `sales` |
| `iat` | unix timestamp | Issued at |
| `exp` | unix timestamp | Expiration (iat + 24h) |

- **Algorithm**: HS256
- **Secret**: `JWT_SECRET` env var (minimum 32 chars)
- **Lifetime**: 24 jam dari `iat`
- **Refresh**: tidak ada refresh token — expired = login ulang
- **Signing**: `pkg/jwt/` package

---

## Auth Flow (App)

```
App start
  → AuthNotifier.init()
  → Cek token di FlutterSecureStorage
      ↓ ada token   → decode, cek exp
          ↓ valid   → redirect ke /clients
          ↓ expired → hapus token, redirect ke /login
      ↓ tidak ada   → /login

Login
  → POST /api/v1/auth/login { email, password }
  → Response: { token: "eyJ..." }
  → Validasi: decode JWT, cek role in [sales, project_owner, superuser]
  → Simpan token di SecureStorage key: "auth_token"
  → Register FCM token: POST /api/v1/devices
  → Redirect ke /clients

Logout
  → Unregister FCM: DELETE /api/v1/devices/{token}
  → Hapus token dari SecureStorage
  → Redirect ke /login

Token Expired (mid-session)
  → Dio interceptor catch 401
  → Hapus token dari SecureStorage
  → AuthNotifier.logout()
  → Redirect ke /login
```

### Dio Interceptor

```dart
// core/network/api_client.dart
_dio.interceptors.add(InterceptorsWrapper(
  onRequest: (options, handler) async {
    final token = await _storage.read(key: 'auth_token');
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  },
  onError: (error, handler) async {
    if (error.response?.statusCode == 401) {
      await _storage.delete(key: 'auth_token');
      _authNotifier.logout();
    }
    handler.next(error);
  },
));
```

---

## Setup Flow (First Run)

```
GET /api/v1/setup/status
  → { "is_setup": false }

POST /api/v1/setup/install
  → Body: { name, email, password }
  → Creates superuser account
  → Response: { user: {...}, token: "eyJ..." }
  → Redirect ke /clients
```

Di Flutter: `SetupNotifier` cek status saat app start. Jika belum setup → tampilkan Setup Wizard sebelum login.

---

## Provisioning Flow

Provisioning = push lisensi ke FlashERP instance agar instance tahu batasan-batasannya.

### Sequence

```
Developer App                Vernon API                  FlashERP Instance
     |                           |                             |
     | POST /licenses/{id}/      |                             |
     |        provision           |                             |
     |-------------------------->|                             |
     |                           | POST {flasherp_url}/api/    |
     |                           |        v1/license/setup     |
     |                           |---------------------------->|
     |                           |   Headers:                  |
     |                           |     X-API-Key: {provision_  |
     |                           |                api_key}     |
     |                           |   Body:                     |
     |                           |     { license_key,          |
     |                           |       product, plan,        |
     |                           |       status, modules,      |
     |                           |       apps, max_users,      |
     |                           |       max_trans_per_month,   |
     |                           |       max_trans_per_day,     |
     |                           |       max_items,            |
     |                           |       max_customers,        |
     |                           |       max_branches,         |
     |                           |       expires_at }          |
     |                           |                             |
     |                           |      200 OK                 |
     |                           |<----------------------------|
     |                           |                             |
     |                           | UPDATE client_licenses      |
     |                           |   SET is_provisioned = true |
     |                           |                             |
     |                           | INSERT audit_logs           |
     |                           |   action = 'license_        |
     |                           |            provisioned'     |
     |                           |                             |
     |     200 OK                |                             |
     |<--------------------------|                             |
```

### Key Points
- `flasherp_url` harus sudah diisi sebelum provisioning
- `provision_api_key` dikirim via `X-API-Key` header — **tidak pernah** di-return ke client
- Jika FlashERP instance return non-200 → return error `LICENSE_PROVISION_FAILED` (502)
- Jika gagal → kirim notification `provision_failed` ke user
- Jika berhasil → set `is_provisioned = true` + kirim notification `provision_success`
- Re-provision (jika `is_provisioned` sudah true) diperbolehkan — untuk update constraints

### Command Handler

```go
// internal/command/provision_license/handler.go
type ProvisionLicense struct {
    LicenseID uuid.UUID
    ActorID   uuid.UUID
}
// CommandBus key: "provision_license.ProvisionLicense"
```

### Failure Handling
- HTTP timeout: 30 detik
- Tidak ada automatic retry — developer bisa retry manual dari app
- Semua provisioning attempt (success + failure) dicatat di audit log
- Metadata audit: `{ "status_code": 200, "response_time_ms": 450 }` atau `{ "error": "connection timeout" }`

---

## Rate Limiting

### Public Endpoints
- 60 requests/menit per IP
- Middleware: `pkg/middleware/ratelimit.go`
- Implementation: `golang.org/x/time/rate` dengan `map[string]*rate.Limiter` per IP
- Cleanup: background goroutine setiap 5 menit hapus stale entries

### Protected Endpoints
- 300 requests/menit per user (dari JWT `sub` claim)

### Response Headers
```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1711234627
```

### Jika Exceeded
```
HTTP 429 Too Many Requests
{ "error": { "code": "RATE_LIMIT_EXCEEDED", "message": "Too many requests" } }
```
