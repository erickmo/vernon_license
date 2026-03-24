# App Features (Go WASM — internal, direct DB access)

Vernon App mengelola semua data secara internal. **Tidak menggunakan public API** — langsung ke PostgreSQL via repository layer.

## Project Structure

```
app-developer/
├── main.go                          ← Server (//go:build !wasm) — serves WASM + 2 public API endpoints
├── main_wasm.go                     ← WASM entry (//go:build wasm) — App routes
├── internal/
│   ├── publicapi/                   ← The 2 public endpoints (register + validate)
│   │   ├── register_handler.go
│   │   └── validate_handler.go
│   ├── auth/                        ← App JWT auth (NOT for public API)
│   ├── repository/                  ← Direct DB access for App
│   ├── service/                     ← Business logic (proposal approval → license creation, PDF, etc.)
│   ├── errors/
│   └── ui/
│       ├── components/
│       └── pages/
│           ├── login.go
│           ├── dashboard.go
│           ├── companies_list.go / company_detail.go / company_form.go
│           ├── project_detail.go     ← Tabs: Licenses, Proposals, Activity
│           ├── licenses_list.go / license_detail.go / license_create.go
│           ├── proposals_list.go / proposal_detail.go / proposal_form.go
│           ├── products_list.go / product_form.go    ← superuser
│           ├── users_list.go / user_form.go          ← superuser
│           └── notifications.go
└── web/css/ + img/
```

Key: `publicapi/` package handles the 2 public endpoints. Everything else is App-internal.

---

## Features

### Companies + Projects
- CRUD, all roles. Project groups licenses + proposals.

### Licenses
- **LicenseDetailPage** tabs: Info, Registration Status, Activity
- **Registration Status tab**: is_registered, instance_url, instance_name, last_pull_at, provision_api_key (copy button), check_interval
- **Create license** (PO): set product, plan, modules, constraints → generate `provision_api_key` → give to client team
- **Suspend/activate** (PO): toggle status → next validate call returns false/true

### Proposals (versioned + changelog)
- Sales create → PO can edit submitted (modify pricing/modules) → PO approve → auto-create license + PDF
- Changelog: auto-computed diff, inline annotations, reviewer auto-focuses changelog tab
- PDF generated on approval only

### Products (superuser)
- Dynamic CRUD: name, slug, modules, apps, plans, pricing

### Users (superuser)
- Create PO / sales, deactivate

### Dashboard
- License counts by status, revenue, expiring soon, recent activity, pending proposals

### Notifications + Audit
- Notifications: proposal events, expiry warnings, client registered
- Audit: all entity changes, tracked per handler

---

## Role Visibility

| Element | `sales` | `project_owner` | `superuser` |
|---|---|---|---|
| Companies / Projects | ✅ | ✅ | ✅ |
| View licenses | ✅ | ✅ | ✅ |
| Create / submit proposal | ✅ | ✅ | ✅ |
| Download approved PDF | ✅ | ✅ | ✅ |
| Edit submitted proposal | ❌ | ✅ | ✅ |
| Approve / reject proposal | ❌ | ✅ | ✅ |
| Create license directly | ❌ | ✅ | ✅ |
| Suspend / activate | ❌ | ✅ | ✅ |
| Renew | ❌ | ✅ | ✅ |
| Manage products | ❌ | ❌ | ✅ |
| Manage users | ❌ | ❌ | ✅ |
| Global audit | ❌ | ❌ | ✅ |
