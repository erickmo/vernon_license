# Domain Models

## Company

| Field | Type | Keterangan |
|---|---|---|
| `id` | UUID | PK |
| `name` | string | Nama perusahaan |
| `email` | string? | Email kontak |
| `phone` | string? | |
| `address` | text? | |
| `pic_name` | string? | Person in charge |
| `pic_email` | string? | |
| `pic_phone` | string? | |
| `notes` | text? | |
| `created_by` | UUID | FK users |
| `created_at` | timestamptz | |
| `updated_at` | timestamptz | |
| `deleted_at` | time? | Soft delete |

---

## Project

Groups licenses + proposals under one engagement.

| Field | Type | Keterangan |
|---|---|---|
| `id` | UUID | PK |
| `company_id` | UUID | FK companies |
| `name` | string | |
| `description` | text? | |
| `status` | string | `active` \| `completed` \| `cancelled` |
| `created_by` | UUID | FK users |
| `created_at` | timestamptz | |
| `updated_at` | timestamptz | |
| `deleted_at` | time? | |

Relationships: Company 1→N Projects, Project 1→N Licenses, Project 1→N Proposals

---

## License

On/off toggle. Client app registers → gets license_key → periodically validates.

| Field | Type | Keterangan |
|---|---|---|
| `id` | UUID | PK |
| `license_key` | string | `FL-XXXXXXXX` (auto-generate, unique) |
| `project_id` | UUID | FK projects |
| `company_id` | UUID | FK companies (denormalized) |
| `product_id` | UUID | FK products |
| `plan` | string | `saas` \| `dedicated` |
| `status` | string | `active` \| `trial` \| `suspended` \| `expired` \| `pending` |
| `modules` | TEXT[] | Fitur aktif (internal tracking) |
| `apps` | TEXT[] | Apps aktif (internal tracking) |
| `contract_amount` | decimal? | Nilai kontrak |
| `description` | text? | |
| `max_users` | int? | Batas pengguna |
| `max_trans_per_month` | int? | |
| `max_trans_per_day` | int? | |
| `max_items` | int? | |
| `max_customers` | int? | |
| `max_branches` | int? | |
| `max_storage` | int? | MB |
| `expires_at` | time? | |
| `instance_url` | string? | URL deployment (set saat register) |
| `instance_name` | string? | Nama instance (set saat register) |
| `provision_api_key` | string? | Key untuk register auth (json `"-"`) |
| `check_interval` | string | `1h` \| `6h` \| `24h` (default `6h`) |
| `last_pull_at` | time? | Timestamp validate terakhir |
| `is_registered` | bool | Client app sudah call register |
| `proposal_id` | UUID? | FK proposals (NULL jika direct create) |
| `created_by` | UUID | FK users |
| `created_at` | timestamptz | |
| `updated_at` | timestamptz | |
| `deleted_at` | time? | |
| `archived_at` | time? | |

### Status Flow
```
pending   → active       (PO approve — atau auto jika pre-approved)
pending   → trial        (PO set trial)
active    ←→ suspended   (PO toggle)
active    → expired      (by expires_at or manual)
trial     → active       (PO approve)
expired   → active       (PO renew)
any       → archived
```

### Validate Logic (public API)
```
valid = true  IF status IN (active) AND (expires_at IS NULL OR expires_at > now)
valid = true  IF status IN (trial) AND approved
valid = false IF status IN (suspended, expired, pending)
```

### Two Creation Paths
1. **Via Proposal**: Sales → PO approve → auto-create license (status: pending until client registers, or active if pre-provisioned)
2. **Direct**: PO creates in Vernon App (sets provision_api_key, client registers later)

### Registration Flow
1. PO creates license in Vernon App → sets `provision_api_key` → gives key to client
2. Client app calls `POST /api/v1/register` with the key
3. Vernon matches by `provision_api_key` + `product_slug`
4. Sets `instance_url`, `instance_name`, `is_registered = true`
5. Returns `license_key` + `valid` status
6. Client uses `license_key` for periodic `GET /api/v1/validate`

---

## Product (superuser CRUD — internal)

| Field | Type | Keterangan |
|---|---|---|
| `id` | UUID | PK |
| `name` | string | |
| `slug` | string | Unique (used in register) |
| `description` | text? | |
| `available_modules` | JSONB | `[{ key, name, description }]` |
| `available_apps` | JSONB | `[{ key, name, description }]` |
| `available_plans` | TEXT[] | `["saas", "dedicated"]` |
| `base_pricing` | JSONB | `{ "saas": { base_price, per_user_price, currency } }` |
| `is_active` | bool | |
| `created_at` | timestamptz | |
| `updated_at` | timestamptz | |
| `deleted_at` | time? | |

---

## Proposal (versioned + changelog — internal)

| Field | Type | Keterangan |
|---|---|---|
| `id` | UUID | PK |
| `project_id` | UUID | FK projects |
| `company_id` | UUID | FK companies |
| `product_id` | UUID | FK products |
| `version` | int | Auto-increment per project+product |
| `status` | string | `draft` \| `submitted` \| `approved` \| `rejected` |
| `modules` | TEXT[] | |
| `apps` | TEXT[] | |
| `plan` | string | |
| `max_users` | int? | |
| `max_trans_per_month` | int? | |
| `max_trans_per_day` | int? | |
| `max_items` | int? | |
| `max_customers` | int? | |
| `max_branches` | int? | |
| `max_storage` | int? | |
| `contract_amount` | decimal? | |
| `expires_at` | time? | |
| `notes` | text? | Dari sales |
| `owner_notes` | text? | Dari PO |
| `rejection_reason` | text? | |
| `changelog` | JSONB? | Auto-computed diff vs previous |
| `pdf_path` | string? | |
| `pdf_generated_at` | time? | |
| `submitted_by` | UUID | |
| `reviewed_by` | UUID? | |
| `reviewed_at` | time? | |
| `created_at` | timestamptz | |
| `updated_at` | timestamptz | |

### Approval → auto-create license + generate PDF
PO can edit submitted proposals before approving.

### Changelog (auto-computed)
```json
{ "compared_to_version": 1, "summary": "...", "changes": [...], "unchanged": [...] }
```

---

## User, AuditLog, Notification

Same as before — see `docs/MIGRATIONS.md` for schemas. All managed internally via App.

### Audit Actions
`license_created` · `license_updated` · `status_changed` · `license_renewed` · `client_registered` · `proposal_created` · `proposal_submitted` · `proposal_edited_by_owner` · `proposal_approved` · `proposal_rejected` · `product_created` · `product_updated` · `company_created` · `project_created` · `user_created` · `user_login`

### Notification Types
`proposal_submitted` → PO · `proposal_approved` → sales · `proposal_rejected` → sales · `proposal_edited` → sales · `license_expiring_30d` → PO+sales · `license_expiring_7d` → PO+sales · `license_expired` → PO · `client_registered` → PO (new client registered)
