//go:build wasm

package pages

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// PublicAPIDocsPage menampilkan dokumentasi Public API Vernon License.
type PublicAPIDocsPage struct {
	app.Compo
}

func (p *PublicAPIDocsPage) Render() app.UI {
	return app.Div().
		Style("min-height", "100vh").
		Style("background", "#0F0A1A").
		Style("color", "#E2D9F3").
		Style("font-family", "'Inter', -apple-system, sans-serif").
		Style("padding", "0").
		Body(
			p.renderNav(),
			app.Div().
				Style("max-width", "960px").
				Style("margin", "0 auto").
				Style("padding", "40px 24px 80px").
				Body(
					p.renderTitle(),
					p.renderBaseInfo(),
					p.renderRateLimit(),
					p.renderRegisterSection(),
					p.renderValidateSection(),
					p.renderValidateOTPSection(),
					p.renderErrorCodes(),
				),
		)
}

func (p *PublicAPIDocsPage) renderNav() app.UI {
	return app.Div().
		Style("background", "#1A1035").
		Style("border-bottom", "1px solid rgba(77,41,117,0.4)").
		Style("padding", "16px 24px").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "space-between").
		Body(
			app.Div().
				Style("display", "flex").
				Style("align-items", "center").
				Style("gap", "12px").
				Body(
					app.Div().
						Style("width", "32px").
						Style("height", "32px").
						Style("background", "linear-gradient(135deg, #4D2975, #26B8B0)").
						Style("border-radius", "8px").
						Style("display", "flex").
						Style("align-items", "center").
						Style("justify-content", "center").
						Style("font-size", "14px").
						Text("V"),
					app.Span().
						Style("font-size", "15px").
						Style("font-weight", "600").
						Style("color", "#E2D9F3").
						Text("Vernon License API"),
				),
			app.Div().
				Style("display", "flex").
				Style("gap", "24px").
				Style("font-size", "13px").
				Body(
					app.A().
						Href("#register").
						Style("color", "#9B8DB5").
						Style("text-decoration", "none").
						Text("Register"),
					app.A().
						Href("#validate").
						Style("color", "#9B8DB5").
						Style("text-decoration", "none").
						Text("Validate"),
					app.A().
						Href("#errors").
						Style("color", "#9B8DB5").
						Style("text-decoration", "none").
						Text("Errors"),
				),
		)
}

func (p *PublicAPIDocsPage) renderTitle() app.UI {
	return app.Div().
		Style("margin-bottom", "40px").
		Body(
			app.H1().
				Style("font-size", "32px").
				Style("font-weight", "700").
				Style("color", "#E2D9F3").
				Style("margin", "0 0 12px").
				Style("letter-spacing", "-0.5px").
				Text("Public API Reference"),
			app.P().
				Style("font-size", "15px").
				Style("color", "#9B8DB5").
				Style("line-height", "1.6").
				Style("margin", "0").
				Text("Endpoints untuk registrasi lisensi dan validasi. Tidak memerlukan autentikasi (no Bearer token)."),
		)
}

func (p *PublicAPIDocsPage) renderBaseInfo() app.UI {
	return app.Div().
		Style("margin-bottom", "32px").
		Body(
			p.renderInfoCard("Base URL", "http://localhost:8081"),
			p.renderInfoCard("Content-Type", "application/json"),
		)
}

func (p *PublicAPIDocsPage) renderInfoCard(label, value string) app.UI {
	return app.Div().
		Style("display", "inline-flex").
		Style("align-items", "center").
		Style("gap", "8px").
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.4)").
		Style("border-radius", "8px").
		Style("padding", "8px 16px").
		Style("margin-right", "12px").
		Style("margin-bottom", "8px").
		Body(
			app.Span().
				Style("font-size", "12px").
				Style("color", "#9B8DB5").
				Text(label+": "),
			app.Span().
				Style("font-size", "13px").
				Style("color", "#26B8B0").
				Style("font-family", "monospace").
				Text(value),
		)
}

func (p *PublicAPIDocsPage) renderRateLimit() app.UI {
	return app.Div().
		Style("background", "rgba(38,184,176,0.08)").
		Style("border", "1px solid rgba(38,184,176,0.3)").
		Style("border-radius", "10px").
		Style("padding", "14px 20px").
		Style("margin-bottom", "40px").
		Style("font-size", "13px").
		Style("color", "#26B8B0").
		Style("display", "flex").
		Style("align-items", "center").
		Style("gap", "8px").
		Body(
			app.Span().Style("font-weight", "600").Text("Rate Limit:"),
			app.Span().Text("60 requests / menit per IP address"),
		)
}

// ── Register ────────────────────────────────────────────────────────

func (p *PublicAPIDocsPage) renderRegisterSection() app.UI {
	return app.Div().
		ID("register").
		Style("margin-bottom", "48px").
		Body(
			p.renderEndpointHeader("POST", "/api/v1/register", "Registrasi lisensi baru untuk client app."),
			p.renderSectionTitle("Request Body"),
			p.renderCodeBlock(`{
  "app_name":     "string  (wajib) — slug produk, contoh: \"vernon-erp\"",
  "otp":          "string  (wajib) — kode OTP yang masih aktif",
  "client_name":  "string  (wajib) — nama perusahaan/client, contoh: \"PT Maju Jaya\"",
  "instance_url": "string  (wajib) — URL instance client, contoh: \"https://app.client.com\""
}`),
			p.renderFieldTable([]fieldRow{
				{"app_name", "string", "Ya", "Slug produk yang terdaftar di sistem (= product slug)"},
				{"otp", "string", "Ya", "Kode OTP aktif dari tabel otp (harus belum expired)"},
				{"client_name", "string", "Ya", "Nama perusahaan/client (digunakan untuk find-or-create company)"},
				{"instance_url", "string", "Ya", "URL tempat client di-deploy (digunakan untuk callback superuser creation)"},
			}),
			p.renderSectionTitle("Flow"),
			p.renderFlowSteps([]string{
				"Validasi semua field wajib ada",
				"Capture client_app_ip dari request sender IP",
				"Cek produk ada berdasarkan app_name (slug)",
				"Cek OTP aktif (belum expired) di tabel otp",
				"Find-or-create company berdasarkan client_name",
				"Cek kombinasi company + product belum punya lisensi",
				"Generate license key (format: FL-XXXXXXXX)",
				"Buat lisensi baru dengan status \"pending\" + client_app_ip + instance_url",
				"Catat audit log (action: client_registered)",
			}),
			p.renderSectionTitle("Success Response — 201 Created"),
			p.renderCodeBlock(`{
  "license_key":    "FL-A1B2C3D4",
  "product":        "Vernon ERP",
  "check_interval": "6h",
  "status":         "pending",
  "message":        "Registration received. License is pending approval."
}`),
			p.renderResponseTable([]fieldRow{
				{"license_key", "string", "", "License key unik (FL-XXXXXXXX)"},
				{"product", "string", "", "Nama produk (bukan slug)"},
				{"check_interval", "string", "", "Interval pengecekan validasi, default \"6h\""},
				{"status", "string", "", "Selalu \"pending\" saat baru registrasi"},
				{"message", "string", "", "Pesan konfirmasi"},
			}),
			p.renderSectionTitle("Error Responses"),
			p.renderErrorTable([]errorRow{
				{"400", "VALIDATION_FAILED", "Field wajib kosong atau body JSON invalid"},
				{"403", "PRODUCT_NOT_FOUND", "app_name tidak ditemukan di database"},
				{"403", "INVALID_CLIENT_CODE", "OTP invalid atau sudah expired"},
				{"409", "ALREADY_REGISTERED", "Kombinasi company + product sudah punya lisensi"},
				{"500", "INTERNAL_ERROR", "Kesalahan server internal"},
			}),
			p.renderSectionTitle("Contoh cURL"),
			p.renderCodeBlock(`curl -X POST http://localhost:8081/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "app_name": "vernon-erp",
    "otp": "ABC123",
    "client_name": "PT Maju Jaya",
    "instance_url": "https://app.client.com"
  }'`),
		)
}

// ── Validate OTP ────────────────────────────────────────────────────

func (p *PublicAPIDocsPage) renderValidateOTPSection() app.UI {
	return app.Div().
		ID("validate-otp").
		Style("margin-bottom", "48px").
		Body(
			p.renderEndpointHeader("POST", "/api/v1/validate_otp", "Validasi apakah OTP masih aktif. Digunakan oleh client app saat menerima request dari license app."),
			p.renderSectionTitle("Request Body"),
			p.renderCodeBlock(`{
  "otp": "string (wajib) — kode OTP yang ingin divalidasi"
}`),
			p.renderFieldTable([]fieldRow{
				{"otp", "string", "Ya", "Kode OTP yang ingin dicek validitasnya"},
			}),
			p.renderSectionTitle("Success Response — 200 OK"),
			p.renderCodeBlock(`{
  "status": true
}`),
			p.renderSectionTitle("Invalid OTP — 200 OK"),
			p.renderCodeBlock(`{
  "status": false
}`),
			p.renderSectionTitle("Contoh cURL"),
			p.renderCodeBlock(`curl -X POST http://localhost:8081/api/v1/validate_otp \
  -H "Content-Type: application/json" \
  -d '{"otp": "ABC123"}'`),
		)
}

// ── Validate ────────────────────────────────────────────────────────

func (p *PublicAPIDocsPage) renderValidateSection() app.UI {
	return app.Div().
		ID("validate").
		Style("margin-bottom", "48px").
		Body(
			p.renderEndpointHeader("GET", "/api/v1/validate?key=FL-XXXXXXXX", "Cek apakah lisensi masih valid."),
			p.renderSectionTitle("Query Parameters"),
			p.renderFieldTable([]fieldRow{
				{"key", "string", "Ya", "License key yang didapat saat registrasi (format: FL-XXXXXXXX)"},
			}),
			p.renderSectionTitle("Flow"),
			p.renderFlowSteps([]string{
				"Ambil query param \"key\"",
				"Cari lisensi berdasarkan key",
				"Update last_pull_at untuk monitoring (non-fatal jika gagal)",
				"Gunakan check_interval dari lisensi, fallback ke config default",
				"Cek validitas: status \"active\" + belum expired, atau status \"trial\"",
				"Return valid=true jika valid, atau valid=false dengan reason",
			}),
			p.renderSectionTitle("Valid License — 200 OK"),
			p.renderCodeBlock(`{
  "valid":          true,
  "license_key":    "FL-A1B2C3D4",
  "check_interval": "6h"
}`),
			p.renderSectionTitle("Invalid License — 200 OK"),
			p.renderCodeBlock(`{
  "valid":          false,
  "license_key":    "FL-A1B2C3D4",
  "reason":         "suspended",
  "check_interval": "6h"
}`),
			p.renderSectionTitle("Possible reason Values"),
			p.renderReasonTable([]reasonRow{
				{"suspended", "Lisensi ditangguhkan oleh admin"},
				{"expired", "Lisensi sudah melewati tanggal kadaluarsa (expires_at)"},
				{"pending_approval", "Lisensi belum disetujui (status pending atau trial tanpa active)"},
			}),
			p.renderSectionTitle("Error Responses"),
			p.renderErrorTable([]errorRow{
				{"400", "VALIDATION_FAILED", "Query param 'key' kosong"},
				{"404", "LICENSE_NOT_FOUND", "License key tidak ditemukan di database"},
				{"500", "INTERNAL_ERROR", "Kesalahan server internal"},
			}),
			p.renderSectionTitle("Contoh cURL"),
			p.renderCodeBlock(`curl "http://localhost:8081/api/v1/validate?key=FL-A1B2C3D4"`),
		)
}

// ── Error Codes ─────────────────────────────────────────────────────

func (p *PublicAPIDocsPage) renderErrorCodes() app.UI {
	return app.Div().
		ID("errors").
		Style("margin-bottom", "48px").
		Body(
			app.H2().
				Style("font-size", "22px").
				Style("font-weight", "600").
				Style("color", "#E2D9F3").
				Style("margin", "0 0 16px").
				Text("Error Response Format"),
			app.P().
				Style("font-size", "14px").
				Style("color", "#9B8DB5").
				Style("margin", "0 0 16px").
				Style("line-height", "1.6").
				Text("Semua error response menggunakan format yang sama:"),
			p.renderCodeBlock(`{
  "valid": false,
  "error": {
    "code":    "ERROR_CODE",
    "message": "Human-readable error description"
  }
}`),
			app.H3().
				Style("font-size", "16px").
				Style("font-weight", "600").
				Style("color", "#E2D9F3").
				Style("margin", "24px 0 12px").
				Text("Semua Error Codes"),
			p.renderErrorTable([]errorRow{
				{"400", "VALIDATION_FAILED", "Request body atau parameter tidak valid"},
				{"403", "PRODUCT_NOT_FOUND", "Product slug tidak ditemukan"},
				{"403", "INVALID_CLIENT_CODE", "OTP invalid atau expired"},
				{"404", "LICENSE_NOT_FOUND", "License key tidak ditemukan"},
				{"409", "ALREADY_REGISTERED", "Company+Product sudah terdaftar"},
				{"500", "INTERNAL_ERROR", "Kesalahan server internal"},
			}),
		)
}

// ── Shared Components ───────────────────────────────────────────────

func (p *PublicAPIDocsPage) renderEndpointHeader(method, path, desc string) app.UI {
	methodColor := "#26B8B0"
	if method == "POST" {
		methodColor = "#F59E0B"
	}
	return app.Div().
		Style("margin-bottom", "24px").
		Body(
			app.Div().
				Style("display", "flex").
				Style("align-items", "center").
				Style("gap", "12px").
				Style("margin-bottom", "8px").
				Body(
					app.Span().
						Style("background", methodColor).
						Style("color", "#0F0A1A").
						Style("font-size", "12px").
						Style("font-weight", "700").
						Style("padding", "4px 10px").
						Style("border-radius", "6px").
						Style("font-family", "monospace").
						Text(method),
					app.Span().
						Style("font-size", "16px").
						Style("font-weight", "600").
						Style("color", "#E2D9F3").
						Style("font-family", "monospace").
						Text(path),
				),
			app.P().
				Style("font-size", "14px").
				Style("color", "#9B8DB5").
				Style("margin", "0").
				Style("line-height", "1.6").
				Text(desc),
		)
}

func (p *PublicAPIDocsPage) renderSectionTitle(title string) app.UI {
	return app.H3().
		Style("font-size", "15px").
		Style("font-weight", "600").
		Style("color", "#E2D9F3").
		Style("margin", "24px 0 12px").
		Style("padding-top", "8px").
		Text(title)
}

func (p *PublicAPIDocsPage) renderCodeBlock(code string) app.UI {
	return app.Pre().
		Style("background", "#130E22").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "10px").
		Style("padding", "16px 20px").
		Style("margin", "0 0 16px").
		Style("overflow-x", "auto").
		Style("font-size", "13px").
		Style("line-height", "1.6").
		Style("color", "#C4B5D9").
		Style("font-family", "'JetBrains Mono', 'Fira Code', monospace").
		Body(
			app.Code().Text(code),
		)
}

type fieldRow struct {
	Name     string
	Type     string
	Required string
	Desc     string
}

func (p *PublicAPIDocsPage) renderFieldTable(rows []fieldRow) app.UI {
	headerStyle := func(ui app.UI) app.UI {
		return app.Th().
			Style("text-align", "left").
			Style("padding", "10px 16px").
			Style("font-size", "11px").
			Style("font-weight", "600").
			Style("color", "#9B8DB5").
			Style("text-transform", "uppercase").
			Style("letter-spacing", "0.5px").
			Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
			Body(ui)
	}

	tableRows := []app.UI{
		app.Tr().Body(
			headerStyle(app.Text("Field")),
			headerStyle(app.Text("Type")),
			headerStyle(app.Text("Wajib")),
			headerStyle(app.Text("Deskripsi")),
		),
	}

	for _, row := range rows {
		tableRows = append(tableRows, app.Tr().Body(
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("font-family", "monospace").
				Style("color", "#26B8B0").
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Name),
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("color", "#C4B5D9").
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Type),
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("color", func() string {
					if row.Required == "Ya" {
						return "#F59E0B"
					}
					return "#9B8DB5"
				}()).
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Required),
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("color", "#9B8DB5").
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Desc),
		))
	}

	return app.Div().
		Style("overflow-x", "auto").
		Style("margin-bottom", "16px").
		Body(
			app.Table().
				Style("width", "100%").
				Style("border-collapse", "collapse").
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(77,41,117,0.3)").
				Style("border-radius", "10px").
				Style("overflow", "hidden").
				Body(tableRows...),
		)
}

func (p *PublicAPIDocsPage) renderResponseTable(rows []fieldRow) app.UI {
	headerStyle := func(ui app.UI) app.UI {
		return app.Th().
			Style("text-align", "left").
			Style("padding", "10px 16px").
			Style("font-size", "11px").
			Style("font-weight", "600").
			Style("color", "#9B8DB5").
			Style("text-transform", "uppercase").
			Style("letter-spacing", "0.5px").
			Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
			Body(ui)
	}

	tableRows := []app.UI{
		app.Tr().Body(
			headerStyle(app.Text("Field")),
			headerStyle(app.Text("Type")),
			headerStyle(app.Text("Deskripsi")),
		),
	}

	for _, row := range rows {
		tableRows = append(tableRows, app.Tr().Body(
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("font-family", "monospace").
				Style("color", "#26B8B0").
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Name),
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("color", "#C4B5D9").
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Type),
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("color", "#9B8DB5").
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Desc),
		))
	}

	return app.Div().
		Style("overflow-x", "auto").
		Style("margin-bottom", "16px").
		Body(
			app.Table().
				Style("width", "100%").
				Style("border-collapse", "collapse").
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(77,41,117,0.3)").
				Style("border-radius", "10px").
				Style("overflow", "hidden").
				Body(tableRows...),
		)
}

type errorRow struct {
	Status  string
	Code    string
	Message string
}

func (p *PublicAPIDocsPage) renderErrorTable(rows []errorRow) app.UI {
	headerStyle := func(ui app.UI) app.UI {
		return app.Th().
			Style("text-align", "left").
			Style("padding", "10px 16px").
			Style("font-size", "11px").
			Style("font-weight", "600").
			Style("color", "#9B8DB5").
			Style("text-transform", "uppercase").
			Style("letter-spacing", "0.5px").
			Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
			Body(ui)
	}

	tableRows := []app.UI{
		app.Tr().Body(
			headerStyle(app.Text("HTTP")),
			headerStyle(app.Text("Code")),
			headerStyle(app.Text("Keterangan")),
		),
	}

	for _, row := range rows {
		tableRows = append(tableRows, app.Tr().Body(
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("font-weight", "600").
				Style("color", "#EF4444").
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Status),
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("font-family", "monospace").
				Style("color", "#C4B5D9").
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Code),
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("color", "#9B8DB5").
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Message),
		))
	}

	return app.Div().
		Style("overflow-x", "auto").
		Style("margin-bottom", "16px").
		Body(
			app.Table().
				Style("width", "100%").
				Style("border-collapse", "collapse").
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(77,41,117,0.3)").
				Style("border-radius", "10px").
				Style("overflow", "hidden").
				Body(tableRows...),
		)
}

type reasonRow struct {
	Value string
	Desc  string
}

func (p *PublicAPIDocsPage) renderReasonTable(rows []reasonRow) app.UI {
	headerStyle := func(ui app.UI) app.UI {
		return app.Th().
			Style("text-align", "left").
			Style("padding", "10px 16px").
			Style("font-size", "11px").
			Style("font-weight", "600").
			Style("color", "#9B8DB5").
			Style("text-transform", "uppercase").
			Style("letter-spacing", "0.5px").
			Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
			Body(ui)
	}

	tableRows := []app.UI{
		app.Tr().Body(
			headerStyle(app.Text("Reason")),
			headerStyle(app.Text("Deskripsi")),
		),
	}

	for _, row := range rows {
		tableRows = append(tableRows, app.Tr().Body(
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("font-family", "monospace").
				Style("color", "#F59E0B").
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Value),
			app.Td().
				Style("padding", "10px 16px").
				Style("font-size", "13px").
				Style("color", "#9B8DB5").
				Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
				Text(row.Desc),
		))
	}

	return app.Div().
		Style("overflow-x", "auto").
		Style("margin-bottom", "16px").
		Body(
			app.Table().
				Style("width", "100%").
				Style("border-collapse", "collapse").
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(77,41,117,0.3)").
				Style("border-radius", "10px").
				Style("overflow", "hidden").
				Body(tableRows...),
		)
}

func (p *PublicAPIDocsPage) renderFlowSteps(steps []string) app.UI {
	stepUIs := make([]app.UI, len(steps))
	for i, step := range steps {
		num := string(rune('1' + i))
		stepUIs[i] = app.Div().
			Style("display", "flex").
			Style("align-items", "flex-start").
			Style("gap", "12px").
			Style("padding", "8px 0").
			Body(
				app.Span().
					Style("min-width", "24px").
					Style("height", "24px").
					Style("background", "rgba(77,41,117,0.4)").
					Style("border-radius", "6px").
					Style("display", "flex").
					Style("align-items", "center").
					Style("justify-content", "center").
					Style("font-size", "12px").
					Style("font-weight", "600").
					Style("color", "#26B8B0").
					Text(num),
				app.Span().
					Style("font-size", "13px").
					Style("color", "#C4B5D9").
					Style("line-height", "24px").
					Text(step),
			)
	}
	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "10px").
		Style("padding", "12px 20px").
		Style("margin-bottom", "16px").
		Body(stepUIs...)
}
