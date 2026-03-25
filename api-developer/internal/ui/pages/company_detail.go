//go:build wasm

package pages

import (
	"context"
	"strings"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/components"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// companyLicenseItem adalah ringkasan license untuk ditampilkan di company detail.
type companyLicenseItem struct {
	ID           string  `json:"id"`
	LicenseKey   string  `json:"license_key"`
	ProductName  string  `json:"product_name"`
	Plan         string  `json:"plan"`
	Status       string  `json:"status"`
	IsRegistered bool    `json:"is_registered"`
	ExpiresAt    *string `json:"expires_at"`
	InstanceURL  string  `json:"instance_url"`
}

// CompanyDetailPage menampilkan detail company beserta lisensi-lisensinya.
type CompanyDetailPage struct {
	app.Compo
	companyID string
	company   *CompanyItem
	licenses  []companyLicenseItem
	loading   bool
	errMsg    string
	authStore store.AuthStore
}

// OnNav dipanggil saat halaman ini di-navigasi.
func (p *CompanyDetailPage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}
	path := ctx.Page().URL().Path
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) >= 2 {
		p.companyID = parts[1]
	}
	p.fetchAll(ctx)
}

// fetchAll mengambil company detail dan licenses secara paralel.
func (p *CompanyDetailPage) fetchAll(ctx app.Context) {
	if p.companyID == "" {
		return
	}
	p.loading = true
	p.errMsg = ""

	token := p.authStore.GetToken()
	id := p.companyID

	ctx.Async(func() {
		client := api.NewClient("", token)

		var detail CompanyItem
		errCompany := client.Get(context.Background(), "/api/internal/companies/"+id, &detail)

		var lics []companyLicenseItem
		errLic := client.Get(context.Background(), "/api/internal/companies/"+id+"/licenses", &lics)

		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			if errCompany != nil {
				p.errMsg = errCompany.Error()
				return
			}
			p.company = &detail
			if errLic == nil {
				p.licenses = lics
			}
		})
	})
}

// onBackClick navigasi kembali ke companies list.
func (p *CompanyDetailPage) onBackClick(ctx app.Context, e app.Event) {
	ctx.Navigate("/companies")
}

// onLicenseClick navigasi ke detail license.
func (p *CompanyDetailPage) onLicenseClick(id string) func(app.Context, app.Event) {
	return func(ctx app.Context, e app.Event) {
		ctx.Navigate("/licenses/" + id)
	}
}

// Render menampilkan halaman detail company.
func (p *CompanyDetailPage) Render() app.UI {
	if !p.authStore.IsLoggedIn() {
		return app.Div()
	}
	return app.Elem("x-shell").Body(
		&components.Shell{Content: p.renderContent()},
	)
}

func (p *CompanyDetailPage) renderContent() app.UI {
	return app.Div().
		Style("padding", "32px").
		Body(
			// Back
			app.Button().
				Style("background", "none").
				Style("border", "none").
				Style("color", "#26B8B0").
				Style("font-size", "14px").
				Style("cursor", "pointer").
				Style("padding", "0 0 20px").
				Style("display", "flex").
				Style("align-items", "center").
				Style("gap", "6px").
				OnClick(p.onBackClick).
				Body(
					app.Raw(`<svg style="width:16px;height:16px" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/></svg>`),
					app.Text("Kembali ke Companies"),
				),

			app.If(p.loading, func() app.UI {
				return app.Div().
					Style("text-align", "center").Style("padding", "80px").Style("color", "#9B8DB5").
					Text("Memuat...")
			}),

			app.If(!p.loading && p.errMsg != "", func() app.UI {
				return app.Div().
					Style("background", "rgba(239,68,68,0.1)").Style("border", "1px solid #EF4444").
					Style("border-radius", "8px").Style("padding", "12px 16px").Style("color", "#EF4444").
					Text(p.errMsg)
			}),

			app.If(!p.loading && p.company != nil, func() app.UI {
				return p.renderDetail()
			}),
		)
}

func (p *CompanyDetailPage) renderDetail() app.UI {
	c := p.company
	initial := "?"
	if len(c.Name) > 0 {
		initial = strings.ToUpper(string([]rune(c.Name)[0]))
	}

	// Count licenses by status
	total := len(p.licenses)
	active, pending, suspended := 0, 0, 0
	for _, l := range p.licenses {
		switch l.Status {
		case "active":
			active++
		case "pending":
			pending++
		case "suspended":
			suspended++
		}
	}

	return app.Div().Body(
		// ── Header card ────────────────────────────────────────
		app.Div().
			Style("background", "#1A1035").
			Style("border", "1px solid rgba(77,41,117,0.3)").
			Style("border-radius", "12px").
			Style("padding", "28px 32px").
			Style("display", "flex").
			Style("align-items", "center").
			Style("gap", "24px").
			Style("margin-bottom", "20px").
			Body(
				// Avatar
				app.Div().
					Style("width", "64px").Style("height", "64px").
					Style("border-radius", "14px").
					Style("background", "linear-gradient(135deg, #4D2975, #26B8B0)").
					Style("display", "flex").Style("align-items", "center").Style("justify-content", "center").
					Style("color", "#E2D9F3").Style("font-size", "26px").Style("font-weight", "700").
					Style("flex-shrink", "0").
					Text(initial),
				app.Div().Style("flex", "1").Body(
					app.Div().
						Style("color", "#E2D9F3").Style("font-size", "22px").Style("font-weight", "700").
						Text(c.Name),
					app.If(c.CreatedAt != "", func() app.UI {
						return app.Div().
							Style("color", "#9B8DB5").Style("font-size", "12px").Style("margin-top", "4px").
							Text("Terdaftar sejak "+c.CreatedAt)
					}),
				),
				// License count badge
				app.Div().
					Style("text-align", "center").Style("flex-shrink", "0").
					Body(
						app.Div().
							Style("color", "#E2D9F3").Style("font-size", "28px").Style("font-weight", "700").
							Text(itoa(total)),
						app.Div().
							Style("color", "#9B8DB5").Style("font-size", "11px").Style("text-transform", "uppercase").
							Text("Lisensi"),
					),
			),

		// ── Stats row ───────────────────────────────────────────
		app.If(total > 0, func() app.UI {
			return app.Div().
				Style("display", "grid").
				Style("grid-template-columns", "repeat(3, 1fr)").
				Style("gap", "12px").
				Style("margin-bottom", "20px").
				Body(
					statCard("Active", itoa(active), "#22C55E"),
					statCard("Pending", itoa(pending), "#F59E0B"),
					statCard("Suspended", itoa(suspended), "#EF4444"),
				)
		}),

		// ── Two-column info ─────────────────────────────────────
		app.Div().
			Style("display", "grid").
			Style("grid-template-columns", "1fr 1fr").
			Style("gap", "16px").
			Style("margin-bottom", "20px").
			Body(
				// Contact info
				infoSection("Informasi Kontak",
					infoRow("Email", safeDeref(c.Email)),
					infoRow("Telepon", safeDeref(c.Phone)),
					infoRow("Alamat", safeDeref(c.Address)),
				),
				// PIC info
				infoSection("PIC (Person In Charge)",
					infoRow("Nama", safeDeref(c.PICName)),
					infoRow("Email", safeDeref(c.PICEmail)),
					infoRow("Telepon", safeDeref(c.PICPhone)),
				),
			),

		// Notes (full width, only if filled)
		app.If(safeDeref(c.Notes) != "", func() app.UI {
			return app.Div().
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(77,41,117,0.3)").
				Style("border-radius", "12px").
				Style("padding", "20px 24px").
				Style("margin-bottom", "20px").
				Body(
					sectionTitle("Catatan"),
					app.Div().
						Style("color", "#E2D9F3").Style("font-size", "14px").
						Style("line-height", "1.6").Style("white-space", "pre-wrap").
						Text(safeDeref(c.Notes)),
				)
		}),

		// ── Licenses table ──────────────────────────────────────
		app.Div().
			Style("background", "#1A1035").
			Style("border", "1px solid rgba(77,41,117,0.3)").
			Style("border-radius", "12px").
			Style("overflow", "hidden").
			Body(
				app.Div().
					Style("padding", "18px 24px").
					Style("border-bottom", "1px solid rgba(77,41,117,0.2)").
					Body(sectionTitle("Lisensi")),
				app.If(len(p.licenses) == 0, func() app.UI {
					return app.Div().
						Style("padding", "40px").Style("text-align", "center").Style("color", "#9B8DB5").
						Text("Belum ada lisensi untuk company ini.")
				}),
				app.If(len(p.licenses) > 0, func() app.UI {
					return p.renderLicenseTable()
				}),
			),
	)
}

func (p *CompanyDetailPage) renderLicenseTable() app.UI {
	rows := make([]app.UI, 0, len(p.licenses))
	for _, l := range p.licenses {
		l := l
		exp := "—"
		if l.ExpiresAt != nil && len(*l.ExpiresAt) >= 10 {
			exp = (*l.ExpiresAt)[:10]
		}
		regText, regColor := "Belum", "#F59E0B"
		if l.IsRegistered {
			regText, regColor = "Ya", "#22C55E"
		}
		url := l.InstanceURL
		if url == "" {
			url = "—"
		}
		rows = append(rows, app.Tr().
			Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
			Style("cursor", "pointer").
			OnClick(p.onLicenseClick(l.ID)).
			Body(
				app.Td().Style("padding", "11px 14px").Style("font-family", "monospace").Style("font-size", "12px").Style("color", "#E2D9F3").Text(l.LicenseKey),
				app.Td().Style("padding", "11px 14px").Style("font-size", "13px").Style("color", "#9B8DB5").Text(l.ProductName),
				app.Td().Style("padding", "11px 14px").Style("font-size", "12px").Style("color", "#9B8DB5").Text(l.Plan),
				app.Td().Style("padding", "11px 14px").Body(statusBadge(l.Status)),
				app.Td().Style("padding", "11px 14px").Style("font-size", "12px").Style("color", regColor).Text(regText),
				app.Td().Style("padding", "11px 14px").Style("font-size", "12px").Style("color", "#9B8DB5").Style("max-width", "180px").Style("overflow", "hidden").Style("text-overflow", "ellipsis").Style("white-space", "nowrap").Text(url),
				app.Td().Style("padding", "11px 14px").Style("font-size", "12px").Style("color", "#9B8DB5").Text(exp),
			),
		)
	}
	return app.Table().
		Style("width", "100%").
		Style("border-collapse", "collapse").
		Body(
			app.THead().Style("background", "rgba(77,41,117,0.15)").Body(
				app.Tr().Body(
					thCell("License Key"),
					thCell("Product"),
					thCell("Plan"),
					thCell("Status"),
					thCell("Registered"),
					thCell("Instance URL"),
					thCell("Expires"),
				),
			),
			app.TBody().Body(rows...),
		)
}

// ── helpers ────────────────────────────────────────────────────────────

func statCard(label, value, color string) app.UI {
	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "10px").
		Style("padding", "16px 20px").
		Style("text-align", "center").
		Body(
			app.Div().Style("color", color).Style("font-size", "26px").Style("font-weight", "700").Text(value),
			app.Div().Style("color", "#9B8DB5").Style("font-size", "11px").Style("margin-top", "4px").Style("text-transform", "uppercase").Text(label),
		)
}

func infoSection(title string, rows ...app.UI) app.UI {
	body := make([]app.UI, 0, len(rows)+1)
	body = append(body, sectionTitle(title))
	body = append(body, rows...)
	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "20px 24px").
		Body(body...)
}

func sectionTitle(title string) app.UI {
	return app.Div().
		Style("color", "#9B8DB5").
		Style("font-size", "11px").
		Style("font-weight", "600").
		Style("text-transform", "uppercase").
		Style("letter-spacing", "0.06em").
		Style("margin-bottom", "14px").
		Text(title)
}

func infoRow(label, value string) app.UI {
	if value == "" {
		value = "—"
	}
	return app.Div().
		Style("display", "flex").
		Style("justify-content", "space-between").
		Style("align-items", "flex-start").
		Style("padding", "7px 0").
		Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
		Body(
			app.Span().Style("color", "#9B8DB5").Style("font-size", "13px").Text(label),
			app.Span().Style("color", "#E2D9F3").Style("font-size", "13px").Style("text-align", "right").Style("max-width", "60%").Style("word-break", "break-word").Text(value),
		)
}

func thCell(label string) app.UI {
	return app.Th().
		Style("padding", "10px 14px").
		Style("color", "#9B8DB5").
		Style("font-size", "11px").
		Style("font-weight", "600").
		Style("text-align", "left").
		Style("text-transform", "uppercase").
		Style("white-space", "nowrap").
		Text(label)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}
