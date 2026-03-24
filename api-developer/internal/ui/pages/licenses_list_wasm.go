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

// LicenseListItem adalah representasi ringkas license untuk tampilan daftar.
type LicenseListItem struct {
	ID           string  `json:"id"`
	LicenseKey   string  `json:"license_key"`
	CompanyName  string  `json:"company_name"`
	ProjectName  string  `json:"project_name"`
	ProductName  string  `json:"product_name"`
	Plan         string  `json:"plan"`
	Status       string  `json:"status"`
	IsRegistered bool    `json:"is_registered"`
	ExpiresAt    *string `json:"expires_at"`
}

// LicensesListPage menampilkan semua license (global view) dengan filter status.
// Tombol "Buat License" hanya tampil untuk project_owner dan superuser.
type LicensesListPage struct {
	app.Compo
	licenses     []LicenseListItem
	filtered     []LicenseListItem
	loading      bool
	errMsg       string
	statusFilter string // "" = all, "active", "pending", "suspended", "expired"
	searchQuery  string
	authStore    store.AuthStore
}

// OnNav dipanggil saat halaman ini di-navigasi.
// Redirect ke /login jika belum login, lalu fetch licenses.
func (p *LicensesListPage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}
	p.statusFilter = ""
	p.searchQuery = ""
	p.fetchLicenses(ctx)
}

// fetchLicenses mengambil daftar license dari API dan menyimpan ke state.
func (p *LicensesListPage) fetchLicenses(ctx app.Context) {
	p.loading = true
	p.errMsg = ""

	token := p.authStore.GetToken()
	ctx.Async(func() {
		client := api.NewClient("", token)
		var items []LicenseListItem
		if err := client.Get(context.Background(), "/api/internal/licenses", &items); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.loading = false
				p.errMsg = err.Error()
			})
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			p.licenses = items
			p.applyFilter()
		})
	})
}

// applyFilter memfilter p.licenses berdasarkan statusFilter dan searchQuery.
func (p *LicensesListPage) applyFilter() {
	result := make([]LicenseListItem, 0, len(p.licenses))
	q := strings.ToLower(p.searchQuery)
	for _, l := range p.licenses {
		if p.statusFilter != "" && l.Status != p.statusFilter {
			continue
		}
		if q != "" {
			match := strings.Contains(strings.ToLower(l.LicenseKey), q) ||
				strings.Contains(strings.ToLower(l.CompanyName), q) ||
				strings.Contains(strings.ToLower(l.ProjectName), q) ||
				strings.Contains(strings.ToLower(l.ProductName), q)
			if !match {
				continue
			}
		}
		result = append(result, l)
	}
	p.filtered = result
}

// onFilterChange dipanggil saat status filter berubah.
func (p *LicensesListPage) onFilterChange(ctx app.Context, e app.Event) {
	p.statusFilter = ctx.JSSrc().Get("value").String()
	p.applyFilter()
}

// onSearchChange dipanggil saat search input berubah.
func (p *LicensesListPage) onSearchChange(ctx app.Context, e app.Event) {
	p.searchQuery = ctx.JSSrc().Get("value").String()
	p.applyFilter()
}

// onCreateClick navigasi ke halaman buat license baru.
func (p *LicensesListPage) onCreateClick(ctx app.Context, e app.Event) {
	ctx.Navigate("/licenses/new")
}

// onViewClick navigasi ke halaman detail license.
func (p *LicensesListPage) onViewClick(id string) func(app.Context, app.Event) {
	return func(ctx app.Context, e app.Event) {
		ctx.Navigate("/licenses/" + id)
	}
}

// Render menampilkan halaman daftar licenses dalam Shell.
func (p *LicensesListPage) Render() app.UI {
	if !p.authStore.IsLoggedIn() {
		return app.Div()
	}

	return app.Elem("x-shell").
		Body(
			&components.Shell{
				Content: p.renderContent(),
			},
		)
}

// renderContent merender area utama halaman licenses list.
func (p *LicensesListPage) renderContent() app.UI {
	return app.Div().
		Style("padding", "32px").
		Body(
			// Header
			app.Div().
				Style("display", "flex").
				Style("align-items", "center").
				Style("justify-content", "space-between").
				Style("margin-bottom", "24px").
				Body(
					app.H1().
						Style("color", "#E2D9F3").
						Style("font-size", "24px").
						Style("font-weight", "700").
						Style("margin", "0").
						Text("Licenses"),
					app.If(p.authStore.HasRole("project_owner"),
						func() app.UI {
							return app.Button().
								Style("background", "linear-gradient(135deg, #4D2975, #26B8B0)").
								Style("color", "#E2D9F3").
								Style("border", "none").
								Style("border-radius", "8px").
								Style("padding", "10px 20px").
								Style("font-size", "14px").
								Style("font-weight", "600").
								Style("cursor", "pointer").
								OnClick(p.onCreateClick).
								Text("+ Buat License")
						},
					),
				),

			// Error message
			app.If(p.errMsg != "",
				func() app.UI {
					return app.Div().
						Style("background", "rgba(239,68,68,0.15)").
						Style("border", "1px solid rgba(239,68,68,0.4)").
						Style("border-radius", "8px").
						Style("padding", "12px 16px").
						Style("color", "#EF4444").
						Style("font-size", "14px").
						Style("margin-bottom", "16px").
						Text(p.errMsg)
				},
			),

			// Filter bar
			app.Div().
				Style("display", "flex").
				Style("gap", "12px").
				Style("margin-bottom", "20px").
				Body(
					// Search input
					app.Input().
						Type("text").
						Placeholder("Cari license key, company, project, product...").
						Value(p.searchQuery).
						Style("flex", "1").
						Style("background", "#1A1035").
						Style("border", "1px solid rgba(77,41,117,0.4)").
						Style("border-radius", "8px").
						Style("padding", "10px 14px").
						Style("color", "#E2D9F3").
						Style("font-size", "14px").
						Style("outline", "none").
						OnInput(p.onSearchChange),
					// Status filter
					app.Select().
						Style("background", "#1A1035").
						Style("border", "1px solid rgba(77,41,117,0.4)").
						Style("border-radius", "8px").
						Style("padding", "10px 14px").
						Style("color", "#E2D9F3").
						Style("font-size", "14px").
						Style("cursor", "pointer").
						OnChange(p.onFilterChange).
						Body(
							app.Option().Value("").Text("Semua Status"),
							app.Option().Value("active").Text("Active"),
							app.Option().Value("pending").Text("Pending"),
							app.Option().Value("trial").Text("Trial"),
							app.Option().Value("suspended").Text("Suspended"),
							app.Option().Value("expired").Text("Expired"),
						),
				),

			// Loading state
			app.If(p.loading,
				func() app.UI {
					return app.Div().
						Style("text-align", "center").
						Style("padding", "60px").
						Style("color", "#9B8DB5").
						Text("Memuat licenses...")
				},
			),

			// Table
			app.If(!p.loading,
				func() app.UI {
					return p.renderTable()
				},
			),
		)
}

// renderTable merender tabel license atau pesan kosong.
func (p *LicensesListPage) renderTable() app.UI {
	if len(p.filtered) == 0 {
		msg := "Belum ada license."
		if p.statusFilter != "" || p.searchQuery != "" {
			msg = "Tidak ada license yang cocok dengan filter."
		}
		return app.Div().
			Style("text-align", "center").
			Style("padding", "60px").
			Style("color", "#9B8DB5").
			Style("background", "#1A1035").
			Style("border-radius", "12px").
			Style("border", "1px solid rgba(77,41,117,0.3)").
			Text(msg)
	}

	// Build rows
	rows := make([]app.UI, 0, len(p.filtered))
	for _, l := range p.filtered {
		l := l // capture
		expiresText := "—"
		if l.ExpiresAt != nil && *l.ExpiresAt != "" {
			// Format: show only date portion
			s := *l.ExpiresAt
			if len(s) >= 10 {
				expiresText = s[:10]
			}
		}
		registeredText := "Belum"
		registeredColor := "#F59E0B"
		if l.IsRegistered {
			registeredText = "Ya"
			registeredColor = "#22C55E"
		}

		rows = append(rows, app.Tr().
			Style("border-bottom", "1px solid rgba(77,41,117,0.2)").
			Style("cursor", "pointer").
			OnClick(p.onViewClick(l.ID)).
			Body(
				app.Td().Style("padding", "12px 14px").Style("color", "#E2D9F3").Style("font-family", "monospace").Style("font-size", "13px").Text(l.LicenseKey),
				app.Td().Style("padding", "12px 14px").Style("color", "#E2D9F3").Style("font-size", "13px").Text(l.CompanyName),
				app.Td().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "13px").Text(l.ProjectName),
				app.Td().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "13px").Text(l.ProductName),
				app.Td().Style("padding", "12px 14px").Style("font-size", "13px").Text(l.Plan),
				app.Td().Style("padding", "12px 14px").Body(statusBadge(l.Status)),
				app.Td().Style("padding", "12px 14px").Style("color", registeredColor).Style("font-size", "13px").Text(registeredText),
				app.Td().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "13px").Text(expiresText),
				app.Td().Style("padding", "12px 14px").Body(
					app.Button().
						Style("background", "rgba(77,41,117,0.3)").
						Style("color", "#E2D9F3").
						Style("border", "1px solid rgba(77,41,117,0.4)").
						Style("border-radius", "6px").
						Style("padding", "5px 12px").
						Style("font-size", "12px").
						Style("cursor", "pointer").
						OnClick(p.onViewClick(l.ID)).
						Text("Detail"),
				),
			),
		)
	}

	return app.Div().
		Style("background", "#1A1035").
		Style("border-radius", "12px").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("overflow", "hidden").
		Body(
			app.Table().
				Style("width", "100%").
				Style("border-collapse", "collapse").
				Body(
					app.THead().
						Style("background", "rgba(77,41,117,0.2)").
						Body(
							app.Tr().Body(
								app.Th().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "600").Style("text-align", "left").Style("text-transform", "uppercase").Text("License Key"),
								app.Th().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "600").Style("text-align", "left").Style("text-transform", "uppercase").Text("Company"),
								app.Th().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "600").Style("text-align", "left").Style("text-transform", "uppercase").Text("Project"),
								app.Th().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "600").Style("text-align", "left").Style("text-transform", "uppercase").Text("Product"),
								app.Th().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "600").Style("text-align", "left").Style("text-transform", "uppercase").Text("Plan"),
								app.Th().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "600").Style("text-align", "left").Style("text-transform", "uppercase").Text("Status"),
								app.Th().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "600").Style("text-align", "left").Style("text-transform", "uppercase").Text("Registered"),
								app.Th().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "600").Style("text-align", "left").Style("text-transform", "uppercase").Text("Expires"),
								app.Th().Style("padding", "12px 14px").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "600").Style("text-align", "left").Style("text-transform", "uppercase").Text("Aksi"),
							),
						),
					app.TBody().Body(rows...),
				),
		)
}

// statusBadge mengembalikan badge UI dengan warna sesuai status license.
func statusBadge(status string) app.UI {
	color := "#9B8DB5"
	bg := "rgba(155,141,181,0.15)"
	switch status {
	case "active":
		color = "#22C55E"
		bg = "rgba(34,197,94,0.15)"
	case "trial":
		color = "#26B8B0"
		bg = "rgba(38,184,176,0.15)"
	case "pending":
		color = "#F59E0B"
		bg = "rgba(245,158,11,0.15)"
	case "suspended":
		color = "#EF4444"
		bg = "rgba(239,68,68,0.15)"
	case "expired":
		color = "#9B8DB5"
		bg = "rgba(155,141,181,0.15)"
	}

	return app.Span().
		Style("display", "inline-block").
		Style("background", bg).
		Style("color", color).
		Style("border-radius", "6px").
		Style("padding", "3px 10px").
		Style("font-size", "12px").
		Style("font-weight", "600").
		Style("text-transform", "capitalize").
		Text(status)
}
