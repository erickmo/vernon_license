# Flutter Features

## Project Structure

```
app-developer/lib/
├── core/
│   ├── auth/            ← AuthNotifier (ChangeNotifier)
│   ├── constants/       ← AppColors, AppDimensions, AppRoles, ErrorCodes
│   ├── errors/          ← Failure types + error code mapping
│   ├── network/         ← ApiClient (Dio + Bearer interceptor)
│   ├── notifiers/       ← PendingCountNotifier
│   ├── setup/           ← SetupNotifier
│   ├── theme/           ← AppTheme
│   └── widgets/         ← MainShell (root navigation)
├── features/
│   ├── auth/            ← Login
│   ├── clients/         ← Klien aktif list + detail + audit timeline
│   ├── registrations/   ← List, approve, reject, create client
│   ├── dashboard/       ← Analytics + charts
│   ├── app_updates/     ← Release management
│   ├── setup/           ← First-run setup wizard
│   ├── settings/
│   ├── invoices/
│   └── notifications/   ← Push notification list + badge
├── app.dart             ← GoRouter config
├── injection_container.dart ← get_it wiring
└── main.dart
```

---

## Feature: Auth

### Flow
```
App start
  → AuthNotifier.init()
  → Cek token di FlutterSecureStorage
      ↓ ada token   → redirect ke /clients
      ↓ tidak ada   → /login
Login
  → POST /api/v1/auth/login
  → Validasi role in [sales, project_owner, superuser]
  → Simpan token di SecureStorage
  → Redirect ke /clients
Token expired (401 from API)
  → Dio interceptor hapus token dari SecureStorage
  → Redirect ke /login
```

### Files
| File | Type | Deskripsi |
|---|---|---|
| `auth/domain/entities/user.dart` | Entity | User model |
| `auth/domain/usecases/login_usecase.dart` | UseCase | Login logic |
| `auth/presentation/cubit/auth_cubit.dart` | Cubit | States: Initial, Loading, Success, Error |
| `auth/presentation/pages/login_page.dart` | Page | Login form |

---

## Feature: Clients

### Screens
| Screen | File | Deskripsi |
|---|---|---|
| ClientsPage | `clients/presentation/pages/clients_page.dart` | Tab: Klien Aktif + Pengajuan |
| ClientDetailPage | `clients/presentation/pages/client_detail_page.dart` | Detail + edit + audit tab + renew button |
| CreateClientPage | `registrations/presentation/pages/create_client_page.dart` | 4-step wizard |

### Tab Klien Aktif — Filter Status
Filter chips: **Semua** · **Aktif** · **Trial** · **Ditangguhkan** · **Kedaluwarsa**

### CompanyCubit
```dart
load({search, status, page, sortBy, order})  // load dengan filter + pagination
filterByStatus(status)     // ganti status filter
search(query)              // search by name
refresh()                  // reload filter saat ini
loadMore()                 // infinite scroll / next page
toggleSelect(id)           // multi-select mode
bulkUpdateStatus(ids, status)  // bulk action
```

### States
```dart
CompanyInitial
CompanyLoading
CompanyLoaded(List<ClientLicense>, PaginationMeta, Set<String> selectedIds)
CompanyError(Failure)
```

### ClientDetailPage — Tabs
| Tab | Deskripsi |
|---|---|
| Info | Detail lisensi, modul, apps, constraints |
| Activity | Audit timeline (AuditCubit) |

### Action Buttons (ClientDetailPage)
- **Edit** → Navigate ke edit form (update constraints)
- **Download Proposal** → `GET /licenses/{id}/proposal` → open PDF / share
- **Regenerate Proposal** → `POST /licenses/{id}/proposal/regenerate` → re-generate jika data berubah
- **Renew** → BottomSheet (date picker + contract amount + notes) — **hanya `project_owner` / `superuser`**
- **Suspend/Activate** → Confirmation dialog — **hanya `project_owner` / `superuser`**
- **Approve** → Approve client (trial → active) — **hanya `project_owner` / `superuser`**
- **Provision** → Push ke FlashERP — **hanya `project_owner` / `superuser`**
- **Copy License Key** → Clipboard

### Role-based UI Visibility

| Element | `sales` | `project_owner` | `superuser` |
|---|---|---|---|
| List / detail lisensi | ✅ | ✅ | ✅ |
| **Create client / proposal** (FAB) | ✅ | ✅ | ✅ |
| **Download / regenerate proposal** | ✅ | ✅ | ✅ |
| Edit / update constraints | ✅ | ✅ | ✅ |
| Export CSV | ✅ | ✅ | ✅ |
| **Approve client** (trial → active) | ❌ | ✅ | ✅ |
| **Suspend / activate** | ❌ | ✅ | ✅ |
| **Renew lisensi** | ❌ | ✅ | ✅ |
| **Provision** ke FlashERP | ❌ | ✅ | ✅ |
| **Bulk update status** | ❌ | ✅ | ✅ |
| **Manage products** (CRUD) | ❌ | ✅ | ✅ |
| Dashboard | ✅ | ✅ | ✅ |
| Notifications | ✅ | ✅ | ✅ |
| Global audit log | ❌ | ❌ | ✅ |

### Multi-select Mode (ClientsPage)
- Long-press → masuk multi-select
- Floating action bar: Export (all roles), Suspend/Activate (**`project_owner`/`superuser` only**)
- CompanyCubit.selectedIds tracking

---

## Feature: Dashboard

### DashboardCubit
```dart
load()      // fetch GET /api/v1/dashboard
refresh()   // pull-to-refresh
```

### States
```dart
DashboardInitial
DashboardLoading
DashboardLoaded(DashboardData)
DashboardError(Failure)
```

### Widgets
| Widget | Deskripsi |
|---|---|
| SummaryCards | 4 kartu: Active, Trial, Suspended, Expired (color-coded) |
| RevenueChart | Bar chart by product (`fl_chart`) |
| ExpiringSoonList | List tile + countdown badge (merah < 7d, amber < 30d) |
| RecentActivityTimeline | Compact timeline dari audit log |
| DistributionDonut | SaaS vs Dedicated pie chart |

---

## Feature: Notifications

### NotificationCubit
```dart
load({page})            // load paginated
markAsRead(id)          // PUT /notifications/{id}/read
markAllAsRead()         // PUT /notifications/read-all
getUnreadCount()        // GET /notifications/unread-count
```

### FCM Integration
- Register token: `POST /api/v1/devices` saat login
- Unregister: `DELETE /api/v1/devices/{token}` saat logout
- Handle foreground: show in-app banner
- Handle background: system notification → deep link ke detail

### Badge Count
`PendingCountNotifier` di `core/notifiers/` — poll `unread-count` setiap 60 detik, show badge di bottom nav.

---

## Feature: Registrations (Create Client Wizard)

### 5-step Wizard
| Step | Fields | Sumber Data |
|---|---|---|
| 1. Info Klien | client_name, client_email, description | Manual input |
| 2. Produk | product (dropdown dari API), plan | `GET /api/v1/products` → dropdown |
| 3. Modul & Apps | modules[] (checkbox), apps[] (checkbox) | Filtered dari `product.available_modules` / `available_apps` |
| 4. Batasan & Pricing | max_users, max_trans_per_month, max_trans_per_day, max_items, max_customers, max_branches, contract_amount (auto-calculated + editable) | `product.base_pricing` untuk auto-calc |
| 5. Review & Submit | Summary semua field + expires_at + confirm → submit → proposal PDF ready | |

### Pricing Auto-calculation (Step 4)
```
contract_amount = base_price + (per_user_price × max_users)
```
- Auto-filled dari `product.base_pricing[selected_plan]`
- User bisa override `contract_amount` secara manual (custom deal)
- Show breakdown: base + per-user = total

### CreateClientCubit
```dart
loadProducts()           // GET /api/v1/products → populate dropdown
selectProduct(productId) // filter available modules/apps/plans
updateStep(step, data)   // simpan data per step
calculatePricing()       // auto-calc dari product.base_pricing + max_users
submit()                 // POST /api/v1/licenses → auto-generates proposal
downloadProposal(id)     // GET /api/v1/licenses/{id}/proposal → open/share PDF
reset()                  // clear wizard
```

### States
```dart
CreateClientInitial
CreateClientProductsLoaded(List<Product>)
CreateClientInProgress(step, data, calculatedPrice)
CreateClientSubmitting
CreateClientSuccess(license, proposalUrl)  // ← show "Download Proposal" button
CreateClientError(Failure)
```

### Post-submit Flow
1. Submit → `POST /api/v1/licenses` → response includes `proposal_url`
2. Show success dialog with "Download Proposal" & "View Client" buttons
3. Download → `GET /api/v1/licenses/{id}/proposal` → open PDF viewer / share sheet

---

## Feature: Products (Management — `project_owner`/`superuser` only)

### Screens
| Screen | File | Deskripsi |
|---|---|---|
| ProductsPage | `products/presentation/pages/products_page.dart` | List semua produk |
| ProductFormPage | `products/presentation/pages/product_form_page.dart` | Create / edit produk |

### ProductCubit
```dart
load()                    // GET /api/v1/products
create(productData)       // POST /api/v1/products
update(id, productData)   // PUT /api/v1/products/{id}
delete(id)                // DELETE /api/v1/products/{id}
```

### Product Form Fields
- name, slug (auto-generate dari name, editable)
- description
- available_plans (checkbox: saas, dedicated)
- available_modules (dynamic list — add/remove `{ key, name, description }`)
- available_apps (dynamic list — add/remove `{ key, name, description }`)
- base_pricing per plan (base_price, per_user_price, currency)

### Visibility
- Products menu item: **hidden** untuk `sales` role
- Products dropdown in create wizard: **visible** for all roles (read-only list)

---

## Feature: Audit (di ClientDetailPage)

### AuditCubit
```dart
loadForLicense(licenseId, {page})   // GET /licenses/{id}/audit
loadMore()                           // pagination
```

### States
```dart
AuditInitial
AuditLoading
AuditLoaded(List<AuditEntry>, PaginationMeta)
AuditError(Failure)
```

### Widget: ActivityTimeline
- Chronological ListView
- Colored icons per action type:
  - `status_changed` → amber
  - `constraints_updated` → blue
  - `license_provisioned` → green
  - `license_renewed` → purple
- Actor name + relative timestamp
- Expandable change diff (old → new)

---

## Core: Error Handling

### Error Mapper
```dart
// core/errors/error_mapper.dart
Failure mapApiError(DioException e) {
  final code = e.response?.data['error']['code'];
  return switch (code) {
    'AUTH_INVALID_CREDENTIALS' => AuthFailure('Email atau password salah'),
    'AUTH_TOKEN_EXPIRED'       => AuthFailure('Sesi habis, silakan login kembali'),
    'AUTH_INSUFFICIENT_ROLE'   => AuthFailure('Anda tidak memiliki akses'),
    'LICENSE_NOT_FOUND'        => LicenseFailure('Lisensi tidak ditemukan'),
    'LICENSE_INVALID_TRANSITION' => LicenseFailure('Perubahan status tidak valid'),
    'VALIDATION_FAILED'        => ValidationFailure('Data yang dikirim tidak valid'),
    'RATE_LIMIT_EXCEEDED'      => NetworkFailure('Terlalu banyak request'),
    _                          => NetworkFailure('Terjadi kesalahan'),
  };
}
```

### Failure Hierarchy
```dart
sealed class Failure {
  String get messageId;
}
class AuthFailure extends Failure { ... }
class LicenseFailure extends Failure { ... }
class NetworkFailure extends Failure { ... }
class ValidationFailure extends Failure { ... }
```
