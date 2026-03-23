# app-developer — FlashERP Developer Portal

## Purpose & Role

App mobile khusus untuk tim `developer_sales` FlashLab. Digunakan untuk meninjau dan menyetujui/menolak permintaan registrasi perusahaan baru dari calon klien.

**Users:** Developer sales FlashLab (role: `developer_sales`)
**Platform:** Android + iOS
**Package:** `flash_developer`

## Tech Stack

| Component | Library |
|---|---|
| State | flutter_bloc / Cubit (^9.1.1) |
| Navigation | go_router ^14.0.0 |
| DI | get_it ^9.2.1 (manual) |
| HTTP | Dio ^5.4.3+1 |
| Token | flutter_secure_storage ^10.0.0 |

## Project Structure

```
lib/
├── core/
│   ├── auth/auth_notifier.dart     ← ChangeNotifier status auth
│   ├── constants/                  ← app_constants, app_colors
│   ├── errors/failures.dart        ← Failure types
│   ├── network/api_client.dart     ← Dio + auth interceptor
│   └── theme/app_theme.dart
├── features/
│   ├── auth/                       ← Login (domain/data/presentation)
│   └── registrations/              ← List + Approve + Reject
├── app.dart                        ← GoRouter (2 routes: /login, /registrations)
├── injection_container.dart        ← get_it manual DI
└── main.dart
```

## API Endpoints

```
POST /api/v1/auth/login                           ← Login (role check: developer_sales)
GET  /api/v1/developer/registrations              ← List registrasi (?status=pending|approved|rejected)
POST /api/v1/developer/registrations/{id}/approve ← Setujui (body: {company_code, company_name})
POST /api/v1/developer/registrations/{id}/reject  ← Tolak (body: {reason})
```

## Key Features

1. **Login** — validasi role `developer_sales`, akses ditolak jika bukan
2. **List Registrasi** — filter: pending / disetujui / ditolak / semua
3. **Approve** — input kode + nama perusahaan → buat company baru di sistem
4. **Reject** — input alasan penolakan

## Design Tokens

```
Primary:   #4D2975 (Deep Purple)
Secondary: #E9A800 (Amber Gold)
Accent:    #26B8B0 (Teal)
Success:   #22C55E
Error:     #EF4444
```

## Build & Run

```bash
flutter pub get
flutter run --dart-define=API_BASE_URL=http://localhost:8080
flutter build apk --dart-define=API_BASE_URL=https://api.flasherp.id
```

## Auth Flow

```
App start → AuthNotifier.init() → cek token di FlutterSecureStorage
  ↓ ada token     → GoRouter redirect → /registrations
  ↓ tidak ada     → /login
Login → validasi role developer_sales → simpan token → /registrations
```
