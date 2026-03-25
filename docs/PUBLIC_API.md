# Public API

2 endpoints only. No auth. Rate limit: 60 req/min per IP.

---

## POST /api/v1/register

**Request**
```json
{ "otp": "abc123", "product_slug": "flasherp", "instance_url": "https://co.flasherp.id", "instance_name": "PT Maju Bersama" }
```

**201 Created**
```json
{ "license_key": "FL-A1B2C3D4", "product": "FlashERP", "check_interval": "6h", "status": "pending", "message": "License registered" }
```

**Errors**: `403 INVALID_CLIENT_CODE` · `403 PRODUCT_NOT_FOUND` · `409 ALREADY_REGISTERED`

**Server flow**: validate OTP → find product → find-or-create company → check duplicate → create license (status: pending) → return

---

## GET /api/v1/validate?key=FL-XXXXXXXX

**200 valid**: `{ "valid": true, "license_key": "...", "check_interval": "6h" }`

**200 invalid**: `{ "valid": false, "reason": "suspended|expired|pending_approval", "check_interval": "6h" }`

**404**: license_key tidak ditemukan

**Valid logic**: `status=active AND (no expiry OR expiry > now)` → true. `status=trial` → true. Else → false.

---

## Client Integration

```
Startup → POST /register → simpan license_key + check_interval
Loop    → GET /validate?key={k} setiap check_interval
          valid=true  → lanjut
          valid=false → block + tampilkan reason
          error       → gunakan cached
```
