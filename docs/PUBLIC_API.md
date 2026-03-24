# Public API — 2 Endpoints

Vernon License exposes **only 2 public endpoints**. Everything else is managed internally via the Vernon App.

---

## 1. Register Client App

Client app calls this to register itself with Vernon and receive a license.

```
POST /api/v1/register
```

### Request
```json
{
  "product_slug": "flasherp",
  "instance_url": "https://majubersama.flasherp.id",
  "instance_name": "PT Maju Bersama",
  "provision_api_key": "key-given-during-setup"
}
```

| Field | Type | Required | Keterangan |
|---|---|---|---|
| `product_slug` | string | yes | Slug produk (e.g. `flasherp`, `flashpos`) |
| `instance_url` | string | yes | URL deployment instance |
| `instance_name` | string | yes | Nama client / company |
| `provision_api_key` | string | yes | API key yang diberikan saat setup (proves legitimacy) |

### Validation
1. `provision_api_key` harus cocok dengan key yang sudah di-generate di Vernon untuk instance ini
2. `product_slug` harus merujuk ke product yang aktif
3. `instance_url` belum terdaftar (prevent duplicate registration)

### Response (201 — Created)
```json
{
  "license_key": "FL-A1B2C3D4",
  "product": "FlashERP",
  "check_interval": "6h",
  "valid": true,
  "message": "License registered successfully"
}
```

### Response (403 — Invalid key)
```json
{
  "valid": false,
  "error": { "code": "INVALID_API_KEY", "message": "Invalid provision API key" }
}
```

### Response (409 — Already registered)
```json
{
  "valid": false,
  "error": { "code": "ALREADY_REGISTERED", "message": "Instance URL already registered" }
}
```

### What happens server-side
1. Validate `provision_api_key`
2. Find or create company record from `instance_name`
3. Find or create license linked to the project
4. Generate `license_key` (FL-XXXXXXXX)
5. Set initial status based on whether license is pre-approved:
   - If PO already created/approved the license → `valid: true`
   - If no pre-approval → `valid: false` (pending approval in Vernon App)
6. Return license_key + validity

### Pre-registration flow
Optionally, PO can **pre-create** a license in Vernon App before the client app registers. When client calls register, Vernon matches by `instance_url` + `product_slug` and returns the existing license.

---

## 2. Validate License

Client app calls this **periodically** to check if license is still valid.

```
GET /api/v1/validate?key=FL-A1B2C3D4
```

### Request

| Param | Type | Required | Keterangan |
|---|---|---|---|
| `key` | string (query) | yes | License key dari register response |

### Response (200 — License found)
```json
{
  "valid": true,
  "license_key": "FL-A1B2C3D4",
  "check_interval": "6h"
}
```

Or if suspended/expired/not-approved:
```json
{
  "valid": false,
  "license_key": "FL-A1B2C3D4",
  "reason": "suspended",
  "check_interval": "6h"
}
```

### Response (404 — Unknown key)
```json
{
  "valid": false,
  "error": { "code": "LICENSE_NOT_FOUND", "message": "Unknown license key" }
}
```

### Validation logic (server-side)
```
1. Find license by key
2. If not found → 404
3. If status == "active" AND (expires_at is NULL OR expires_at > now) → valid: true
4. If status == "trial" AND approved → valid: true
5. If status == "suspended" → valid: false, reason: "suspended"
6. If status == "expired" OR expires_at <= now → valid: false, reason: "expired"
7. If status == "trial" AND not approved → valid: false, reason: "pending_approval"
8. Update last_pull_at = now (for monitoring)
```

### `reason` values
| Reason | Meaning | Client should |
|---|---|---|
| `suspended` | PO suspended the license | Show "Akun ditangguhkan. Hubungi admin." |
| `expired` | License past expiry date | Show "Lisensi kedaluwarsa." |
| `pending_approval` | Registered but not yet approved | Show "Menunggu persetujuan." |

---

## Client Integration (~15 lines any language)

### Startup
```
POST /api/v1/register → store license_key + check_interval locally
```

### Periodic check
```
Every check_interval:
  GET /api/v1/validate?key={license_key}
  If valid: true → continue
  If valid: false → block access, show reason
  If network error → use cached result
```

### Go example
```go
// Register (called once on first startup)
func register() (*RegisterResponse, error) {
    body, _ := json.Marshal(RegisterRequest{
        ProductSlug:     "flasherp",
        InstanceURL:     "https://mycompany.flasherp.id",
        InstanceName:    "PT Maju Bersama",
        ProvisionAPIKey: os.Getenv("VERNON_API_KEY"),
    })
    resp, err := http.Post(vernonURL+"/api/v1/register", "application/json", bytes.NewReader(body))
    if err != nil { return nil, err }
    var result RegisterResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}

// Validate (called every check_interval)
func validate(licenseKey string) bool {
    resp, err := http.Get(vernonURL + "/api/v1/validate?key=" + licenseKey)
    if err != nil { return cachedValid } // network error → use cache
    var result ValidateResponse
    json.NewDecoder(resp.Body).Decode(&result)
    cachedValid = result.Valid
    return result.Valid
}
```

### Node.js example
```javascript
// Register
const { license_key, check_interval } = await fetch(`${VERNON}/api/v1/register`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ product_slug: 'flasherp', instance_url: URL, instance_name: NAME, provision_api_key: KEY })
}).then(r => r.json());

// Validate
async function validate() {
  try {
    const { valid, reason } = await fetch(`${VERNON}/api/v1/validate?key=${LICENSE_KEY}`).then(r => r.json());
    if (!valid) blockAccess(reason);
  } catch { /* use cached */ }
}
setInterval(validate, 6 * 60 * 60 * 1000); // every 6h
```

### PHP example
```php
// Register
$resp = json_decode(file_get_contents(VERNON_URL . '/api/v1/register', false,
    stream_context_create(['http' => [
        'method' => 'POST', 'header' => 'Content-Type: application/json',
        'content' => json_encode(['product_slug' => 'flasherp', 'instance_url' => URL,
            'instance_name' => NAME, 'provision_api_key' => KEY])
    ]])), true);
$licenseKey = $resp['license_key'];

// Validate
function validate($key) {
    $resp = @json_decode(file_get_contents(VERNON_URL . "/api/v1/validate?key=$key"), true);
    return $resp && $resp['valid'] === true;
}
```

---

## Rate Limiting

- 60 requests/minute per IP
- Response: `429 { "error": { "code": "RATE_LIMIT_EXCEEDED" } }`

## No Auth on Public API

These 2 endpoints do NOT require JWT. The `provision_api_key` on register is the only authentication. Validate uses the license_key as identifier (not a secret — it's an opaque lookup key).

---

## What is NOT in the public API

Everything below is managed **internally** via Vernon App (direct DB access):

- Companies, projects CRUD
- Products CRUD
- Proposals (create, edit, approve, reject, PDF)
- License management (suspend, activate, renew, constraints)
- User management
- Dashboard, audit, notifications
- Provisioning (push sync to client app)

The public API is intentionally minimal — client apps only need to know: "am I allowed to run?"
