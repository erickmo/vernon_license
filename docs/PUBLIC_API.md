# Public API

3 endpoints. No auth. Rate limit: 60 req/min per IP.

---

## POST /api/v1/register

**Request**
```json
{ "app_name": "flasherp", "otp": "abc123", "client_name": "PT Maju Bersama", "instance_url": "https://co.flasherp.id" }
```

| Field | Required | Description |
|---|---|---|
| app_name | Ya | Product slug (= product identifier) |
| otp | Ya | Kode OTP aktif dari dashboard license app |
| client_name | Ya | Nama perusahaan/client |
| instance_url | Ya | URL instance client (untuk callback superuser creation) |

**201 Created**
```json
{ "license_key": "FL-A1B2C3D4", "product": "FlashERP", "check_interval": "6h", "status": "pending", "message": "Registration received. License is pending approval." }
```

**Errors**: `400 VALIDATION_FAILED` · `403 PRODUCT_NOT_FOUND` · `403 INVALID_CLIENT_CODE` · `409 ALREADY_REGISTERED`

**Server flow**: validate fields → capture client_app_ip → find product by app_name → validate OTP → find-or-create company by client_name → check duplicate → create license (status: pending, client_app_ip, instance_url) → return

---

## GET /api/v1/validate?key=FL-XXXXXXXX

**200 valid**: `{ "valid": true, "license_key": "...", "check_interval": "6h" }`

**200 invalid**: `{ "valid": false, "reason": "suspended|expired|pending_approval", "check_interval": "6h" }`

**404**: license_key tidak ditemukan

**Valid logic**: `status=active AND (no expiry OR expiry > now)` → true. `status=trial` → true. Else → false.

---

## POST /api/v1/validate_otp

Digunakan oleh client app untuk memvalidasi OTP yang dikirim oleh license app saat create-superuser.

**Request**
```json
{ "otp": "abc123" }
```

**200 OK** (valid): `{ "status": true }`

**200 OK** (invalid/expired): `{ "status": false }`

---

## Client Integration

```
Install  → POST /register → simpan license_key, license_app_ip, instance_url ke .env
Loop     → GET /validate?key={k} setiap check_interval
           valid=true  → lanjut
           valid=false → block + tampilkan reason
           error       → gunakan cached

Superuser Creation (dipanggil oleh license app):
  POST /api/v1/create-superuser ← license app kirim {otp, license_key, username, password}
  Client validasi: license_key, sender IP = license_app_ip, otp via POST /api/v1/validate_otp
  Jika valid → create/update superuser → return {username}
```
