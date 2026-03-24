# Claude Code Implementation Prompts

Run sequentially. One prompt per session. CLAUDE.md + docs/ must be in project root.

---

## Phase 1: Database

### Prompt 1.1 — All Migrations
```
Baca docs/MIGRATIONS.md. Create 001–008 migration files exactly as spec. Run make migrate-up.
```

---

## Phase 2: Domain + Repositories

### Prompt 2.1 — All Domain Models + Repos
```
Baca docs/DOMAIN_MODELS.md.
Create domain models + repository implementations for ALL entities:
- internal/domain/company/, project/, product/, client_license/, proposal/, audit/, notification/, user/
- infrastructure/database/ — one repo file per entity
Key: license has provision_api_key, is_registered, last_pull_at, instance_url/name.
Proposal has ComputeChangelog function.
All queries: WHERE deleted_at IS NULL.
```

---

## Phase 3: Public API (2 endpoints)

### Prompt 3.1 — Register + Validate Endpoints
```
Baca docs/PUBLIC_API.md (entire file — the definitive spec).

Create internal/publicapi/:
1. register_handler.go — POST /api/v1/register
   - Validate provision_api_key against client_licenses table
   - Match by provision_api_key + product_slug
   - Set instance_url, instance_name, is_registered = true
   - Return license_key + valid status
   - If license status = active → valid: true. If pending/suspended/expired → valid: false with reason
   - If key invalid → 403. If already registered → 409.

2. validate_handler.go — GET /api/v1/validate?key=FL-XXX
   - Find license by key
   - valid = (status == active AND not expired)
   - Update last_pull_at = now
   - Return { valid, license_key, check_interval, reason? }

3. Rate limiting middleware: 60 req/min per IP

4. Register routes in main.go — these are the ONLY public routes (no JWT required)

5. Audit log: "client_registered" on successful register

NO JWT auth on these endpoints. provision_api_key is the only auth for register.
```

---

## Phase 4: App Internal Services

### Prompt 4.1 — Company + Project + Product Services
```
Create internal/service/ for business logic that the App (WASM) calls directly:
- CompanyService: CRUD
- ProjectService: CRUD
- ProductService: CRUD (superuser only validation)
These are NOT HTTP handlers — they're Go functions called from the WASM app via repository layer.
```

### Prompt 4.2 — License Service
```
Create LicenseService:
- DirectCreate (PO): generate license_key + provision_api_key, set status
- Suspend / Activate / Renew: status transitions with audit
- License key format: FL-XXXXXXXX (8 random uppercase alphanumeric)
- Provision API key: 32-char random string
- All mutations create audit log entries
```

### Prompt 4.3 — Proposal Service (full lifecycle)
```
Baca docs/DOMAIN_MODELS.md section Proposal.
Create ProposalService:
- Create draft (sales) with changelog computation
- Update (sales=draft, PO=submitted — audit: proposal_edited_by_owner)
- Submit for review
- Approve → generate PDF + auto-create license + notify sales
- Reject → set reason + notify sales
- PDF generator: pkg/proposal/ with GeneratePDF function (reference proposal-design.html)
```

### Prompt 4.4 — User + Auth + Audit + Notifications
```
- UserService: CRUD (superuser creates PO/sales), bcrypt passwords
- AuthService: login → JWT (for App only, NOT public API)
- AuditService: log method, inject into all services
- NotificationService: create + list + mark read + unread count + expiry cron (FX lifecycle)
```

---

## Phase 5: Go WASM App

### Prompt 5.1 — Scaffold + Auth + Shell
```
Baca docs/APP_FEATURES.md.
Setup go-app/v10: main.go (serves WASM + 2 public endpoints), main_wasm.go (routes).
Shell with sidebar (role-aware). LoginPage. AuthStore (localStorage). Route guards.
```

### Prompt 5.2 — Companies + Projects Pages
```
CompaniesListPage, CompanyDetailPage, CompanyFormPage.
ProjectDetailPage with 3 tabs (Licenses, Proposals, Activity).
```

### Prompt 5.3 — Licenses Pages
```
LicensesListPage (global). LicenseDetailPage with 3 tabs (Info, Registration Status, Activity).
Registration Status: is_registered, instance_url, last_pull_at, provision_api_key (copy button).
LicenseCreatePage (PO): wizard → generates provision_api_key.
Suspend/activate buttons.
```

### Prompt 5.4 — Proposals Pages
```
ProposalsListPage (per project). ProposalDetailPage (Overview + Changelog tabs).
ProposalFormPage (create/edit, live diff, PO "Save & Approve").
Changelog component. PDF download via JS blob.
```

### Prompt 5.5 — Products + Users + Dashboard + Notifications
```
Products + Users (superuser only, hidden for others).
Dashboard: summary cards, SVG charts, expiring list.
Notifications: list, badge polling.
Audit timeline component (used in license + project detail).
```

### Prompt 5.6 — Role Visibility + PWA
```
Final pass: role checks everywhere, route guards, PWA manifest + service worker.
Test all 3 roles. Test register + validate endpoints from a test client.
```

---

## Notes

- Public API = literally 2 endpoints. Everything else is App-internal.
- `provision_api_key` is generated by Vernon App when PO creates a license, then given to client team
- Client app calls register once → gets license_key → calls validate periodically
- App talks to DB directly via repository layer, NOT through public API
- `make migrate-up` after Phase 1
- `make tidy` after adding deps
- Test: `curl -X POST /api/v1/register` and `curl /api/v1/validate?key=FL-XXX`
