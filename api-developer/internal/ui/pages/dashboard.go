//go:build wasm

package pages

import (
	"fmt"
	"math"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/components"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// provisionKeyItem adalah satu item provision key untuk tampilan superuser.
type provisionKeyItem struct {
	LicenseKey      string `json:"license_key"`
	CompanyName     string `json:"company_name"`
	ProvisionAPIKey string `json:"provision_api_key"`
}

// dashboardStats adalah response dari GET /api/internal/dashboard.
type dashboardStats struct {
	TotalLicenses     int                `json:"total_licenses"`
	ActiveLicenses    int                `json:"active_licenses"`
	PendingLicenses   int                `json:"pending_licenses"`
	SuspendedLicenses int                `json:"suspended_licenses"`
	ExpiredLicenses   int                `json:"expired_licenses"`
	TotalCompanies    int                `json:"total_companies"`
	TotalProposals    int                `json:"total_proposals"`
	PendingProposals  int                `json:"pending_proposals"`
	TotalRevenue      float64            `json:"total_revenue"`
	ExpiringLicenses  []expiringLicense  `json:"expiring_licenses"`
	RecentActivity    []activityItem     `json:"recent_activity"`
	ProvisionKeys     []provisionKeyItem `json:"provision_keys"`
}

// expiringLicense adalah lisensi yang akan expired dalam 30 hari.
type expiringLicense struct {
	ID         string `json:"id"`
	LicenseKey string `json:"license_key"`
	Company    string `json:"company"`
	ExpiresAt  string `json:"expires_at"`
	DaysLeft   int    `json:"days_left"`
}

// activityItem adalah satu item aktivitas terbaru.
type activityItem struct {
	EntityType string `json:"entity_type"`
	Action     string `json:"action"`
	ActorName  string `json:"actor_name"`
	CreatedAt  string `json:"created_at"`
}

// chartSegment adalah satu segmen untuk donut chart.
type chartSegment struct {
	Label string
	Value int
	Color string
}

// DashboardPage menampilkan summary dan statistik Vernon License.
type DashboardPage struct {
	app.Compo
	stats     *dashboardStats
	loading   bool
	errMsg    string
	apiBase   string // origin URL, e.g. "http://localhost:8081"
	authStore store.AuthStore
}

// OnNav dipanggil saat halaman ini di-navigasi.
// Redirect ke /login jika user belum login.
func (p *DashboardPage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}
	p.apiBase = app.Window().Get("location").Get("origin").String()
	p.loadStats(ctx)
}

// loadStats mengambil dashboard stats dari API.
func (p *DashboardPage) loadStats(ctx app.Context) {
	p.loading = true
	p.errMsg = ""

	token := p.authStore.GetToken()

	ctx.Async(func() {
		client := api.NewClient("", token)
		var stats dashboardStats
		err := client.Get(ctx, "/api/internal/dashboard", &stats)

		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			if err != nil {
				p.errMsg = "Gagal memuat dashboard: " + err.Error()
				return
			}
			p.stats = &stats
		})
	})
}

// Render menampilkan shell dengan konten dashboard.
func (p *DashboardPage) Render() app.UI {
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

// renderContent merender area konten dashboard.
func (p *DashboardPage) renderContent() app.UI {
	return app.Div().
		Style("padding", "32px").
		Style("max-width", "1200px").
		Body(
			// Header
			app.Div().
				Style("margin-bottom", "28px").
				Body(
					app.H1().
						Style("color", "#E2D9F3").
						Style("font-size", "24px").
						Style("font-weight", "700").
						Style("margin", "0 0 4px").
						Text("Dashboard"),
					app.P().
						Style("color", "#9B8DB5").
						Style("font-size", "14px").
						Style("margin", "0").
						Text("Ringkasan sistem Vernon License"),
				),

			// Loading state
			app.If(p.loading,
				func() app.UI {
					return app.Div().
						Style("color", "#9B8DB5").
						Style("text-align", "center").
						Style("padding", "48px 0").
						Text("Memuat data...")
				},
			),

			// Error state
			app.If(p.errMsg != "",
				func() app.UI {
					return app.Div().
						Style("background", "rgba(239,68,68,0.1)").
						Style("border", "1px solid #EF4444").
						Style("border-radius", "8px").
						Style("padding", "12px 16px").
						Style("color", "#EF4444").
						Style("font-size", "14px").
						Style("margin-bottom", "24px").
						Text(p.errMsg)
				},
			),

			// Konten utama — hanya tampil jika stats sudah ada
			app.If(!p.loading && p.stats != nil,
				func() app.UI {
					return app.Div().
						Body(
							p.renderSummaryCards(),
							p.renderChartsRow(),
						p.renderAPIInfo(),
							p.renderExpiringTable(),
							p.renderActivityFeed(),
							p.renderProvisionKeys(),
						)
				},
			),
		)
}

// renderSummaryCards merender 4 summary card di baris teratas.
func (p *DashboardPage) renderSummaryCards() app.UI {
	s := p.stats
	return app.Div().
		Style("display", "grid").
		Style("grid-template-columns", "repeat(auto-fit, minmax(220px, 1fr))").
		Style("gap", "16px").
		Style("margin-bottom", "28px").
		Body(
			p.summaryCard(
				"Total Licenses",
				fmt.Sprintf("%d", s.TotalLicenses),
				fmt.Sprintf("%d active · %d pending", s.ActiveLicenses, s.PendingLicenses),
				"#4D2975",
				"M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z",
			),
			p.summaryCard(
				"Revenue",
				fmt.Sprintf("Rp %.0f", s.TotalRevenue),
				"dari lisensi aktif",
				"#22C55E",
				"M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
			),
			p.summaryCard(
				"Expiring Soon",
				fmt.Sprintf("%d", len(s.ExpiringLicenses)),
				"dalam 30 hari ke depan",
				"#F59E0B",
				"M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z",
			),
			p.summaryCard(
				"Pending Proposals",
				fmt.Sprintf("%d", s.PendingProposals),
				fmt.Sprintf("dari %d total proposal", s.TotalProposals),
				"#26B8B0",
				"M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z",
			),
		)
}

// summaryCard merender satu summary card.
func (p *DashboardPage) summaryCard(title, value, subtitle, accentColor, iconPath string) app.UI {
	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "20px").
		Body(
			app.Div().
				Style("display", "flex").
				Style("align-items", "flex-start").
				Style("justify-content", "space-between").
				Style("margin-bottom", "12px").
				Body(
					app.Div().
						Style("flex", "1").
						Style("min-width", "0").
						Style("overflow", "hidden").
						Body(
							app.Div().
								Style("color", "#9B8DB5").
								Style("font-size", "12px").
								Style("font-weight", "500").
								Style("text-transform", "uppercase").
								Style("letter-spacing", "0.05em").
								Style("margin-bottom", "6px").
								Style("overflow", "hidden").
								Style("text-overflow", "ellipsis").
								Style("white-space", "nowrap").
								Text(title),
							app.Div().
								Style("color", "#E2D9F3").
								Style("font-size", "28px").
								Style("font-weight", "700").
								Style("line-height", "1").
								Style("overflow", "hidden").
								Style("text-overflow", "ellipsis").
								Style("white-space", "nowrap").
								Text(value),
						),
					app.Div().
						Style("width", "40px").
						Style("height", "40px").
						Style("background", "rgba(77,41,117,0.2)").
						Style("border-radius", "10px").
						Style("display", "flex").
						Style("align-items", "center").
						Style("justify-content", "center").
						Style("color", accentColor).
						Style("flex-shrink", "0").
						Body(
							app.Raw(`<svg style="width:20px;height:20px" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="` + iconPath + `"/></svg>`),
						),
				),
			app.Div().
				Style("color", "#9B8DB5").
				Style("font-size", "12px").
				Text(subtitle),
		)
}

// renderChartsRow merender baris chart donut dan stats breakdown.
func (p *DashboardPage) renderChartsRow() app.UI {
	s := p.stats
	segments := []chartSegment{
		{Label: "Active", Value: s.ActiveLicenses, Color: "#22C55E"},
		{Label: "Pending", Value: s.PendingLicenses, Color: "#F59E0B"},
		{Label: "Suspended", Value: s.SuspendedLicenses, Color: "#EF4444"},
		{Label: "Expired", Value: s.ExpiredLicenses, Color: "#9B8DB5"},
	}

	return app.Div().
		Style("display", "grid").
		Style("grid-template-columns", "1fr 1fr").
		Style("gap", "16px").
		Style("margin-bottom", "28px").
		Body(
			// Donut chart card
			app.Div().
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(77,41,117,0.3)").
				Style("border-radius", "12px").
				Style("padding", "20px").
				Body(
					app.Div().
						Style("color", "#E2D9F3").
						Style("font-size", "15px").
						Style("font-weight", "600").
						Style("margin-bottom", "20px").
						Text("License Status"),
					app.Div().
						Style("display", "flex").
						Style("align-items", "center").
						Style("gap", "24px").
						Body(
							renderDonutChart(segments, 120),
							p.renderChartLegend(segments),
						),
				),

			// Company & proposal stats card
			app.Div().
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(77,41,117,0.3)").
				Style("border-radius", "12px").
				Style("padding", "20px").
				Body(
					app.Div().
						Style("color", "#E2D9F3").
						Style("font-size", "15px").
						Style("font-weight", "600").
						Style("margin-bottom", "20px").
						Text("Overview"),
					p.renderOverviewStats(),
				),
		)
}

// renderChartLegend merender legenda untuk donut chart.
func (p *DashboardPage) renderChartLegend(segments []chartSegment) app.UI {
	items := make([]app.UI, 0, len(segments))
	for _, seg := range segments {
		seg := seg
		items = append(items, app.Div().
			Style("display", "flex").
			Style("align-items", "center").
			Style("gap", "8px").
			Style("margin-bottom", "8px").
			Body(
				app.Div().
					Style("width", "10px").
					Style("height", "10px").
					Style("border-radius", "50%").
					Style("background", seg.Color).
					Style("flex-shrink", "0"),
				app.Span().
					Style("color", "#9B8DB5").
					Style("font-size", "13px").
					Text(fmt.Sprintf("%s: %d", seg.Label, seg.Value)),
			))
	}
	return app.Div().Body(items...)
}

// renderOverviewStats merender statistik overview dalam grid.
func (p *DashboardPage) renderOverviewStats() app.UI {
	s := p.stats
	stats := []struct {
		Label string
		Value string
		Color string
	}{
		{"Total Companies", fmt.Sprintf("%d", s.TotalCompanies), "#26B8B0"},
		{"Total Proposals", fmt.Sprintf("%d", s.TotalProposals), "#4D2975"},
		{"Pending Proposals", fmt.Sprintf("%d", s.PendingProposals), "#F59E0B"},
		{"Suspended", fmt.Sprintf("%d", s.SuspendedLicenses), "#EF4444"},
	}

	items := make([]app.UI, 0, len(stats))
	for _, st := range stats {
		st := st
		items = append(items, app.Div().
			Style("display", "flex").
			Style("align-items", "center").
			Style("justify-content", "space-between").
			Style("padding", "10px 0").
			Style("border-bottom", "1px solid rgba(77,41,117,0.2)").
			Body(
				app.Span().
					Style("color", "#9B8DB5").
					Style("font-size", "13px").
					Text(st.Label),
				app.Span().
					Style("color", st.Color).
					Style("font-size", "16px").
					Style("font-weight", "700").
					Text(st.Value),
			))
	}
	return app.Div().Body(items...)
}

// renderDonutChart menghasilkan SVG donut chart sederhana.
// segments: [{label, value, color}]
func renderDonutChart(segments []chartSegment, size int) app.UI {
	total := 0
	for _, s := range segments {
		total += s.Value
	}

	cx := size / 2
	cy := size / 2
	r := (size - 20) / 2
	circumference := 2 * math.Pi * float64(r)

	// Bangun semua circle elements sebagai string SVG
	svgContent := fmt.Sprintf(`<svg width="%d" height="%d" style="flex-shrink:0">`, size, size)

	// Background circle
	svgContent += fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d" fill="none" stroke="rgba(77,41,117,0.2)" stroke-width="12"/>`, cx, cy, r)

	if total == 0 {
		// Tampilkan lingkaran penuh abu-abu jika tidak ada data
		svgContent += fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d" fill="none" stroke="#9B8DB5" stroke-width="12"/>`, cx, cy, r)
	} else {
		offset := 0.0
		// Mulai dari atas (rotate -90 degrees = -PI/2)
		startAngle := -math.Pi / 2

		for _, seg := range segments {
			if seg.Value == 0 {
				continue
			}
			fraction := float64(seg.Value) / float64(total)
			dashLen := fraction * circumference
			dashGap := circumference - dashLen

			// Hitung rotasi dalam derajat
			rotateDeg := (offset/circumference)*360 + (startAngle * 180 / math.Pi) + 90

			svgContent += fmt.Sprintf(
				`<circle cx="%d" cy="%d" r="%d" fill="none" stroke="%s" stroke-width="12" stroke-dasharray="%.2f %.2f" transform="rotate(%.2f %d %d)"/>`,
				cx, cy, r, seg.Color, dashLen, dashGap, rotateDeg, cx, cy,
			)
			offset += dashLen
		}
	}

	// Center text
	svgContent += fmt.Sprintf(
		`<text x="%d" y="%d" text-anchor="middle" dominant-baseline="middle" fill="#E2D9F3" font-size="18" font-weight="700">%d</text>`,
		cx, cy, total,
	)
	svgContent += `</svg>`

	return app.Raw(svgContent)
}

// renderExpiringTable merender tabel lisensi yang akan expired.
func (p *DashboardPage) renderExpiringTable() app.UI {
	s := p.stats
	if len(s.ExpiringLicenses) == 0 {
		return app.Div()
	}

	rows := make([]app.UI, 0, len(s.ExpiringLicenses))
	for _, lic := range s.ExpiringLicenses {
		lic := lic
		daysColor := "#22C55E"
		if lic.DaysLeft <= 7 {
			daysColor = "#EF4444"
		} else if lic.DaysLeft <= 14 {
			daysColor = "#F59E0B"
		}

		rows = append(rows, app.Tr().
			Style("border-bottom", "1px solid rgba(77,41,117,0.2)").
			Body(
				app.Td().
					Style("padding", "10px 12px").
					Style("color", "#E2D9F3").
					Style("font-size", "13px").
					Style("font-family", "monospace").
					Text(lic.LicenseKey),
				app.Td().
					Style("padding", "10px 12px").
					Style("color", "#9B8DB5").
					Style("font-size", "13px").
					Text(lic.Company),
				app.Td().
					Style("padding", "10px 12px").
					Style("color", "#9B8DB5").
					Style("font-size", "13px").
					Text(lic.ExpiresAt),
				app.Td().
					Style("padding", "10px 12px").
					Style("color", daysColor).
					Style("font-size", "13px").
					Style("font-weight", "600").
					Text(fmt.Sprintf("%d hari", lic.DaysLeft)),
			))
	}

	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "20px").
		Style("margin-bottom", "28px").
		Body(
			app.Div().
				Style("color", "#E2D9F3").
				Style("font-size", "15px").
				Style("font-weight", "600").
				Style("margin-bottom", "16px").
				Text("Expiring Soon (30 Hari)"),
			app.Div().
				Style("overflow-x", "auto").
				Body(
					app.Table().
						Style("width", "100%").
						Style("border-collapse", "collapse").
						Body(
							app.THead().
								Body(
									app.Tr().
										Style("border-bottom", "1px solid rgba(77,41,117,0.4)").
										Body(
											app.Th().Style("padding", "8px 12px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("License Key"),
											app.Th().Style("padding", "8px 12px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Company"),
											app.Th().Style("padding", "8px 12px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Expires"),
											app.Th().Style("padding", "8px 12px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Days Left"),
										),
								),
							app.TBody().Body(rows...),
						),
				),
		)
}

// renderAPIInfo merender card informasi public API endpoints.
func (p *DashboardPage) renderAPIInfo() app.UI {
	base := p.apiBase
	if base == "" {
		base = "http://your-server"
	}
	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "20px").
		Style("margin-bottom", "28px").
		Body(
			app.Div().
				Style("color", "#E2D9F3").
				Style("font-size", "15px").
				Style("font-weight", "600").
				Style("margin-bottom", "16px").
				Text("Public API Endpoints"),
			app.Div().
				Style("display", "flex").
				Style("flex-direction", "column").
				Style("gap", "10px").
				Body(
					renderAPIEndpoint("POST", base+"/api/v1/register", "Client app mendaftarkan diri → mendapatkan license key"),
					renderAPIEndpoint("GET", base+"/api/v1/validate", "Client app memvalidasi lisensi → returns { valid: true/false }"),
				),
			app.Div().
				Style("margin-top", "12px").
				Style("padding", "10px 12px").
				Style("background", "rgba(77,41,117,0.1)").
				Style("border-radius", "6px").
				Style("color", "#9B8DB5").
				Style("font-size", "12px").
				Text("Rate limit: 60 req/min per IP · Header: X-Provision-Key (register) · X-License-Key (validate)"),
		)
}

// renderAPIEndpoint merender satu baris endpoint.
func renderAPIEndpoint(method, fullURL, desc string) app.UI {
	methodColor := "#26B8B0"
	methodBg := "rgba(38,184,176,0.15)"
	if method == "POST" {
		methodColor = "#E9A800"
		methodBg = "rgba(233,168,0,0.15)"
	}
	return app.Div().
		Style("display", "flex").
		Style("align-items", "flex-start").
		Style("gap", "12px").
		Style("min-width", "0").
		Body(
			app.Span().
				Style("display", "inline-block").
				Style("background", methodBg).
				Style("color", methodColor).
				Style("font-size", "11px").
				Style("font-weight", "700").
				Style("padding", "2px 8px").
				Style("border-radius", "4px").
				Style("font-family", "monospace").
				Style("flex-shrink", "0").
				Style("margin-top", "1px").
				Text(method),
			app.Div().
				Style("min-width", "0").
				Style("overflow", "hidden").
				Body(
					app.Div().
						Style("color", "#E2D9F3").
						Style("font-size", "13px").
						Style("font-family", "monospace").
						Style("font-weight", "500").
						Style("overflow", "hidden").
						Style("text-overflow", "ellipsis").
						Style("white-space", "nowrap").
						Text(fullURL),
					app.Div().
						Style("color", "#9B8DB5").
						Style("font-size", "12px").
						Style("margin-top", "2px").
						Style("word-break", "break-word").
						Text(desc),
				),
		)
}

// renderActivityFeed merender feed aktivitas terbaru.
func (p *DashboardPage) renderActivityFeed() app.UI {
	s := p.stats
	if len(s.RecentActivity) == 0 {
		return app.Div()
	}

	items := make([]app.UI, 0, len(s.RecentActivity))
	for _, act := range s.RecentActivity {
		act := act
		items = append(items, app.Div().
			Style("display", "flex").
			Style("align-items", "flex-start").
			Style("gap", "12px").
			Style("padding", "12px 0").
			Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
			Body(
				// Dot indicator
				app.Div().
					Style("width", "8px").
					Style("height", "8px").
					Style("border-radius", "50%").
					Style("background", "#26B8B0").
					Style("flex-shrink", "0").
					Style("margin-top", "5px"),
				app.Div().
					Style("flex", "1").
					Style("min-width", "0").
					Style("overflow", "hidden").
					Body(
						app.Div().
							Style("color", "#E2D9F3").
							Style("font-size", "13px").
							Style("font-weight", "500").
							Style("overflow", "hidden").
							Style("text-overflow", "ellipsis").
							Style("white-space", "nowrap").
							Text(fmt.Sprintf("%s: %s", act.EntityType, act.Action)),
						app.Div().
							Style("color", "#9B8DB5").
							Style("font-size", "12px").
							Style("margin-top", "2px").
							Text(fmt.Sprintf("oleh %s · %s", act.ActorName, act.CreatedAt)),
					),
			))
	}

	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "20px").
		Body(
			app.Div().
				Style("color", "#E2D9F3").
				Style("font-size", "15px").
				Style("font-weight", "600").
				Style("margin-bottom", "8px").
				Text("Recent Activity"),
			app.Div().Body(items...),
		)
}

// renderProvisionKeys merender tabel provision keys — hanya untuk superuser.
func (p *DashboardPage) renderProvisionKeys() app.UI {
	if !p.authStore.HasRole("superuser") {
		return app.Div()
	}
	s := p.stats
	if len(s.ProvisionKeys) == 0 {
		return app.Div()
	}

	rows := make([]app.UI, 0, len(s.ProvisionKeys))
	for _, pk := range s.ProvisionKeys {
		pk := pk
		rows = append(rows, app.Tr().
			Style("border-bottom", "1px solid rgba(77,41,117,0.2)").
			Body(
				app.Td().
					Style("padding", "10px 12px").
					Style("color", "#9B8DB5").
					Style("font-size", "13px").
					Style("max-width", "160px").
					Style("overflow", "hidden").
					Style("text-overflow", "ellipsis").
					Style("white-space", "nowrap").
					Text(pk.CompanyName),
				app.Td().
					Style("padding", "10px 12px").
					Style("color", "#E2D9F3").
					Style("font-size", "13px").
					Style("font-family", "monospace").
					Style("max-width", "140px").
					Style("overflow", "hidden").
					Style("text-overflow", "ellipsis").
					Style("white-space", "nowrap").
					Text(pk.LicenseKey),
				app.Td().
					Style("padding", "10px 12px").
					Style("font-size", "13px").
					Style("font-family", "monospace").
					Body(
						app.Span().
							Style("display", "inline-block").
							Style("background", "rgba(77,41,117,0.2)").
							Style("color", "#E9A800").
							Style("border", "1px solid rgba(77,41,117,0.4)").
							Style("border-radius", "6px").
							Style("padding", "3px 10px").
							Style("letter-spacing", "0.05em").
							Style("max-width", "320px").
							Style("overflow", "hidden").
							Style("text-overflow", "ellipsis").
							Style("white-space", "nowrap").
							Style("display", "inline-block").
							Text(pk.ProvisionAPIKey),
					),
			),
		)
	}

	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "20px").
		Style("margin-top", "28px").
		Body(
			app.Div().
				Style("display", "flex").
				Style("align-items", "center").
				Style("gap", "10px").
				Style("margin-bottom", "16px").
				Body(
					app.Div().
						Style("color", "#E2D9F3").
						Style("font-size", "15px").
						Style("font-weight", "600").
						Text("Provision Keys"),
					app.Span().
						Style("background", "rgba(239,68,68,0.15)").
						Style("color", "#EF4444").
						Style("font-size", "11px").
						Style("font-weight", "600").
						Style("padding", "2px 8px").
						Style("border-radius", "4px").
						Text("SUPERUSER ONLY"),
				),
			app.Div().
				Style("overflow-x", "auto").
				Body(
					app.Table().
						Style("width", "100%").
						Style("border-collapse", "collapse").
						Body(
							app.THead().Body(
								app.Tr().
									Style("border-bottom", "1px solid rgba(77,41,117,0.4)").
									Body(
										app.Th().Style("padding", "8px 12px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Company"),
										app.Th().Style("padding", "8px 12px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("License Key"),
										app.Th().Style("padding", "8px 12px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Provision API Key"),
									),
							),
							app.TBody().Body(rows...),
						),
				),
		)
}
