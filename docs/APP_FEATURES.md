# App Features (Go WASM)

Vernon App adalah internal management tool. Bicara langsung ke DB via repository layer, bukan public API.

## Pages

| Page | Route | File |
|---|---|---|
| Login | `/login` | `login.go` |
| Dashboard | `/` | `dashboard.go` |
| Companies list | `/companies` | `companies_list.go` |
| Company detail | `/companies/{id}` | `company_detail.go` |
| Project detail | `/projects/{id}` | `project_detail.go` |
| Licenses list | `/licenses` | `licenses_list.go` |
| License detail | `/licenses/{id}` | `license_detail.go` |
| License create | `/licenses/create` | `license_create.go` |
| Proposals list | `/proposals` | `proposals_list.go` |
| Proposal detail | `/proposals/{id}` | `proposal_detail.go` |
| Proposal form | `/proposals/create`, `/proposals/{id}/edit` | `proposal_form.go` |
| Products list | `/products` | `products_list.go` |
| Product detail | `/products/{id}` | `product_detail.go` |
| Users list | `/users` | `users_list.go` |
| User detail | `/users/{id}` | `user_detail.go` |
| Notifications | `/notifications` | `notifications.go` |
| Activity log | `/logs` | `activity_log.go` |

## Key Features

**Dashboard**: license count by status, expiring soon, pending proposals, recent activity

**License Detail** (tabs: Info · Registration · Activity)
- Info: status, plan, constraints, company/project/product
- Registration: OTP (superuser), is_registered, instance_url, last_pull_at
- Actions: activate, suspend, renew, set status (superuser dropdown)
- After any action → auto-refresh license data

**Proposal Flow**: Sales create draft → submit → PO review (edit jika perlu) → approve/reject
- Approve → auto-create license
- PDF generated on approval

**License Create** (direct, PO only): set product, plan, modules, constraints → license dibuat dengan status `pending`

## Role Visibility

| Feature | sales | project_owner | superuser |
|---|---|---|---|
| View all | ✅ | ✅ | ✅ |
| Create/submit proposal | ✅ | ✅ | ✅ |
| Edit submitted proposal | ❌ | ✅ | ✅ |
| Approve/reject proposal | ❌ | ✅ | ✅ |
| Create license directly | ❌ | ✅ | ✅ |
| Activate/suspend/renew | ❌ | ✅ | ✅ |
| Manage products | ❌ | ❌ | ✅ |
| Manage users | ❌ | ❌ | ✅ |
| View OTP | ❌ | ❌ | ✅ |
| Global audit log | ❌ | ❌ | ✅ |
