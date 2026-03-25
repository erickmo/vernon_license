# Domain Models

## Relationships
`Company` 1→N `Project` 1→N `License` + `Proposal`

---

## Company
`id, name, email?, phone?, address?, pic_name?, pic_email?, pic_phone?, notes?, created_by, created_at, updated_at, deleted_at`

## Project
`id, company_id, name, description?, status(active|completed|cancelled), created_by, created_at, updated_at, deleted_at`

## License (client_licenses)
| Field | Notes |
|---|---|
| `license_key` | `FL-XXXXXXXX`, auto-generated |
| `project_id?` | nullable — register flow tidak assign project |
| `company_id?` | nullable — set saat register |
| `product_id` | FK products |
| `status` | `pending|active|trial|suspended|expired` |
| `otp` | registration auth code (json:"-"), per-license |
| `otp_generated_at` | kapan OTP dibuat |
| `is_registered` | true setelah client call /register |
| `instance_url?, instance_name?` | set saat register |
| `check_interval` | `1h|6h|24h` default `6h` |
| `last_pull_at?` | timestamp validate terakhir |
| `proposal_id?` | NULL jika direct create |
| `modules[], apps[]` | TEXT[] |
| `max_users?, max_trans_per_month?, max_trans_per_day?` | constraints |
| `max_items?, max_customers?, max_branches?, max_storage?` | constraints |
| `contract_amount?, expires_at?` | |

**Status flow**: `pending → active ↔ suspended`, `active → expired`, `trial → active`

## Product
`id, name, slug(unique), description?, available_modules(jsonb), available_apps(jsonb), available_plans(text[]), base_pricing(jsonb), is_active`

## Proposal
`id, project_id, company_id, product_id, version(auto), status(draft|submitted|approved|rejected)`
`modules[], apps[], plan, constraints (same as license), contract_amount?, expires_at?`
`notes?, owner_notes?, rejection_reason?, changelog(jsonb)?`
`submitted_by, reviewed_by?, reviewed_at?, pdf_path?`

**Approval → auto-create license**

## OTP (global table)
`id, code(unique), is_active, created_at, expires_at`
Used for dashboard rotating registration display code. Separate from per-license `otp` column.

## AuditLog
`id, entity_type, entity_id, action, actor_id, actor_name, changes(jsonb), metadata(jsonb), created_at`

## Notification
`id, user_id, type, title, message, entity_type?, entity_id?, is_read, created_at`

## User
`id, name, email(unique), password_hash, role(superuser|project_owner|sales), is_active, created_at, updated_at`
