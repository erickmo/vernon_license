//go:build wasm

package pages

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/components"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// LicenseDetail adalah representasi lengkap license untuk tampilan detail.
type LicenseDetail struct {
	ID                     string   `json:"id"`
	LicenseKey             string   `json:"license_key"`
	CompanyID              string   `json:"company_id"`
	CompanyName            string   `json:"company_name"`
	ProjectID              string   `json:"project_id"`
	ProjectName            string   `json:"project_name"`
	ProductName            string   `json:"product_name"`
	Plan                   string   `json:"plan"`
	Status                 string   `json:"status"`
	Modules                []string `json:"modules"`
	Apps                   []string `json:"apps"`
	ContractAmount         *float64 `json:"contract_amount"`
	Description            string   `json:"description"`
	MaxUsers               *int     `json:"max_users"`
	MaxTransPerMonth       *int     `json:"max_trans_per_month"`
	MaxTransPerDay         *int     `json:"max_trans_per_day"`
	MaxItems               *int     `json:"max_items"`
	MaxCustomers           *int     `json:"max_customers"`
	MaxBranches            *int     `json:"max_branches"`
	MaxStorage             *int     `json:"max_storage"`
	ExpiresAt              *string  `json:"expires_at"`
	IsRegistered           bool     `json:"is_registered"`
	InstanceURL            string   `json:"instance_url"`
	InstanceName           string   `json:"instance_name"`
	ClientAppIP       string  `json:"client_app_ip"`
	SuperuserUsername  string  `json:"superuser_username"`
	CheckInterval     string  `json:"check_interval"`
	LastPullAt        *string `json:"last_pull_at"`
}

// LicenseDetailPage menampilkan detail license dengan 3 tabs: Info, Registration, Activity.
type LicenseDetailPage struct {
	app.Compo
	licenseID          string
	license            *LicenseDetail
	activeTab          string // "info" | "registration" | "activity"
	auditLogs          []AuditItem
	loading            bool
	errMsg             string
	showSuspendConfirm bool
	authStore          store.AuthStore
	showStatusDropdown bool

	// Credential dialog — shared untuk activate, reset-superuser
	showCredDialog   bool
	credAction       string // "activate" | "reset-superuser"
	credUsername     string
	credPassword     string
	credShowPassword bool
	credLoading      bool

	// Set status confirm dialog (superuser, untuk status non-credential)
	showSetStatusConfirm bool
	pendingStatus        string
	setStatusLoading     bool

	// Success message (auto-cleared after 3s)
	successMsg string
}

// OnNav dipanggil saat halaman ini di-navigasi.
// Redirect ke /login jika belum login, ambil ID dari URL, lalu fetch data.
func (p *LicenseDetailPage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}
	// Extract license ID from URL path: /licenses/{id}
	path := ctx.Page().URL().Path
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) >= 2 {
		p.licenseID = parts[1]
	}
	p.activeTab = "info"
	p.showSuspendConfirm = false
	p.fetchLicense(ctx)
}

// fetchLicense mengambil detail license dari API.
func (p *LicenseDetailPage) fetchLicense(ctx app.Context) {
	if p.licenseID == "" {
		return
	}
	p.loading = true
	p.errMsg = ""

	token := p.authStore.GetToken()
	licenseID := p.licenseID
	ctx.Async(func() {
		client := api.NewClient("", token)
		var detail LicenseDetail
		if err := client.Get(context.Background(), "/api/internal/licenses/"+licenseID, &detail); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.loading = false
				p.errMsg = err.Error()
			})
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			p.license = &detail
		})
	})
}

// fetchAuditLogs mengambil audit log untuk license ini.
func (p *LicenseDetailPage) fetchAuditLogs(ctx app.Context) {
	if p.licenseID == "" {
		return
	}
	token := p.authStore.GetToken()
	licenseID := p.licenseID
	ctx.Async(func() {
		client := api.NewClient("", token)
		var logs []AuditItem
		if err := client.Get(context.Background(), "/api/internal/licenses/"+licenseID+"/audit", &logs); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.auditLogs = nil
				p.errMsg = err.Error()
			})
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.auditLogs = logs
			p.errMsg = ""
		})
	})
}

// onTabClick dipanggil saat tab di-klik.
func (p *LicenseDetailPage) onTabClick(tab string) func(app.Context, app.Event) {
	return func(ctx app.Context, e app.Event) {
		p.activeTab = tab
		if tab == "activity" && len(p.auditLogs) == 0 {
			p.fetchAuditLogs(ctx)
		}
	}
}

// onActivate membuka credential dialog untuk aktivasi license.
func (p *LicenseDetailPage) onActivate(ctx app.Context, e app.Event) {
	p.credUsername = ""
	p.credPassword = ""
	p.credShowPassword = false
	p.credAction = "activate"
	p.showCredDialog = true
}

// onResetSuperuser membuka credential dialog untuk reset superuser.
func (p *LicenseDetailPage) onResetSuperuser(ctx app.Context, e app.Event) {
	p.credUsername = ""
	p.credPassword = ""
	p.credShowPassword = false
	p.credAction = "reset-superuser"
	p.showCredDialog = true
}

// onCredCancel menutup credential dialog.
func (p *LicenseDetailPage) onCredCancel(ctx app.Context, e app.Event) {
	p.showCredDialog = false
}

// onCredUsernameInput menangani input username di credential dialog.
func (p *LicenseDetailPage) onCredUsernameInput(ctx app.Context, e app.Event) {
	p.credUsername = ctx.JSSrc().Get("value").String()
}

// onCredPasswordInput menangani input password di credential dialog.
func (p *LicenseDetailPage) onCredPasswordInput(ctx app.Context, e app.Event) {
	p.credPassword = ctx.JSSrc().Get("value").String()
}

// onCredTogglePassword toggle show/hide password di credential dialog.
func (p *LicenseDetailPage) onCredTogglePassword(ctx app.Context, e app.Event) {
	p.credShowPassword = !p.credShowPassword
}

// onCredConfirm mengeksekusi aksi sesuai credAction.
func (p *LicenseDetailPage) onCredConfirm(ctx app.Context, e app.Event) {
	if p.credUsername == "" || p.credPassword == "" {
		p.errMsg = "Username dan password wajib diisi"
		return
	}
	p.credLoading = true
	token := p.authStore.GetToken()
	licenseID := p.licenseID
	username := p.credUsername
	password := p.credPassword
	action := p.credAction
	ctx.Async(func() {
		client := api.NewClient("", token)
		body := map[string]string{"username": username, "password": password}
		endpoint := "/api/internal/licenses/" + licenseID + "/" + action
		if err := client.Put(context.Background(), endpoint, body, nil); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.errMsg = err.Error()
				p.credLoading = false
			})
			return
		}
		var detail LicenseDetail
		if err := client.Get(context.Background(), "/api/internal/licenses/"+licenseID, &detail); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.errMsg = err.Error()
				p.credLoading = false
			})
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.errMsg = ""
			p.credLoading = false
			p.showCredDialog = false
			p.license = &detail
			switch action {
			case "activate":
				p.successMsg = "License berhasil diaktifkan"
			case "reset-superuser":
				p.successMsg = "Superuser berhasil direset"
			}
		})
		ctx.After(3*time.Second, func(ctx app.Context) { p.successMsg = "" })
	})
}

// renderPasswordInput merender input password dengan tombol show/hide.
func renderPasswordInput(placeholder string, showPassword bool, onInput, onToggle app.EventHandler) app.UI {
	inputType := "password"
	if showPassword {
		inputType = "text"
	}
	var icon app.UI
	if showPassword {
		icon = app.Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>`)
	} else {
		icon = app.Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/><line x1="1" y1="1" x2="23" y2="23"/></svg>`)
	}
	return app.Div().
		Style("position", "relative").
		Body(
			app.Input().
				Type(inputType).
				Placeholder(placeholder).
				Style("width", "100%").
				Style("padding", "10px 42px 10px 14px").
				Style("background", "rgba(77,41,117,0.15)").
				Style("border", "1px solid rgba(77,41,117,0.4)").
				Style("border-radius", "8px").
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				Style("outline", "none").
				Style("box-sizing", "border-box").
				OnInput(onInput),
			app.Button().
				Type("button").
				Style("position", "absolute").
				Style("right", "10px").
				Style("top", "50%").
				Style("transform", "translateY(-50%)").
				Style("background", "none").
				Style("border", "none").
				Style("padding", "4px").
				Style("cursor", "pointer").
				Style("color", "#9B8DB5").
				Style("display", "flex").
				Style("align-items", "center").
				Style("line-height", "0").
				OnClick(onToggle).
				Body(icon),
		)
}


// onSuspendRequest tampilkan dialog konfirmasi suspend.
func (p *LicenseDetailPage) onSuspendRequest(ctx app.Context, e app.Event) {
	p.showSuspendConfirm = true
}

// onSuspendCancel membatalkan konfirmasi suspend.
func (p *LicenseDetailPage) onSuspendCancel(ctx app.Context, e app.Event) {
	p.showSuspendConfirm = false
}

// onSuspendConfirm mengeksekusi suspend setelah konfirmasi.
func (p *LicenseDetailPage) onSuspendConfirm(ctx app.Context, e app.Event) {
	p.showSuspendConfirm = false
	token := p.authStore.GetToken()
	licenseID := p.licenseID
	ctx.Async(func() {
		client := api.NewClient("", token)
		if err := client.Put(context.Background(), "/api/internal/licenses/"+licenseID+"/suspend", nil, nil); err != nil {
			ctx.Dispatch(func(ctx app.Context) { p.errMsg = err.Error() })
			return
		}
		var detail LicenseDetail
		if err := client.Get(context.Background(), "/api/internal/licenses/"+licenseID, &detail); err != nil {
			ctx.Dispatch(func(ctx app.Context) { p.errMsg = err.Error() })
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.errMsg = ""
			p.license = &detail
		})
	})
}

// onRenew memperbarui license yang expired ke active.
func (p *LicenseDetailPage) onRenew(ctx app.Context, e app.Event) {
	token := p.authStore.GetToken()
	licenseID := p.licenseID
	ctx.Async(func() {
		client := api.NewClient("", token)
		if err := client.Put(context.Background(), "/api/internal/licenses/"+licenseID+"/renew", map[string]any{}, nil); err != nil {
			ctx.Dispatch(func(ctx app.Context) { p.errMsg = err.Error() })
			return
		}
		var detail LicenseDetail
		if err := client.Get(context.Background(), "/api/internal/licenses/"+licenseID, &detail); err != nil {
			ctx.Dispatch(func(ctx app.Context) { p.errMsg = err.Error() })
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.errMsg = ""
			p.license = &detail
		})
	})
}

// onBack navigasi ke daftar licenses.
func (p *LicenseDetailPage) onBack(ctx app.Context, e app.Event) {
	ctx.Navigate("/licenses")
}

// onSetStatus membuka dialog konfirmasi perubahan status (superuser only).
// Untuk status "active" dan "trial", pakai credential dialog (sama seperti Activate).
// Untuk status lainnya, pakai confirm dialog sederhana.
func (p *LicenseDetailPage) onSetStatus(status string) func(app.Context, app.Event) {
	return func(ctx app.Context, e app.Event) {
		p.showStatusDropdown = false
		if status == "active" || status == "trial" {
			p.credUsername = ""
			p.credPassword = ""
			p.credShowPassword = false
			p.credAction = "activate"
			p.showCredDialog = true
			return
		}
		p.pendingStatus = status
		p.setStatusLoading = false
		p.showSetStatusConfirm = true
	}
}

// onSetStatusConfirm mengeksekusi perubahan status setelah konfirmasi.
func (p *LicenseDetailPage) onSetStatusConfirm(ctx app.Context, e app.Event) {
	p.setStatusLoading = true
	token := p.authStore.GetToken()
	licenseID := p.licenseID
	status := p.pendingStatus
	ctx.Async(func() {
		client := api.NewClient("", token)
		body := map[string]string{"status": status}
		if err := client.Put(context.Background(), "/api/internal/licenses/"+licenseID+"/status", body, nil); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.errMsg = err.Error()
				p.setStatusLoading = false
				p.showSetStatusConfirm = false
			})
			return
		}
		var detail LicenseDetail
		if err := client.Get(context.Background(), "/api/internal/licenses/"+licenseID, &detail); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.errMsg = err.Error()
				p.setStatusLoading = false
				p.showSetStatusConfirm = false
			})
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.errMsg = ""
			p.setStatusLoading = false
			p.showSetStatusConfirm = false
			p.license = &detail
			label, _ := statusInfo(status)
			p.successMsg = "Status berhasil diubah ke " + label
		})
		ctx.After(3*time.Second, func(ctx app.Context) { p.successMsg = "" })
	})
}

// onSetStatusCancelConfirm menutup dialog konfirmasi set status.
func (p *LicenseDetailPage) onSetStatusCancelConfirm(ctx app.Context, e app.Event) {
	p.showSetStatusConfirm = false
	p.pendingStatus = ""
}

// onToggleStatusDropdown membuka/menutup dropdown status.
func (p *LicenseDetailPage) onToggleStatusDropdown(ctx app.Context, e app.Event) {
	p.showStatusDropdown = !p.showStatusDropdown
}

// Render menampilkan halaman detail license dalam Shell.
func (p *LicenseDetailPage) Render() app.UI {
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

// renderContent merender area utama detail page.
func (p *LicenseDetailPage) renderContent() app.UI {
	return app.Div().
		Style("padding", "32px").
		Body(
			// Back button + title
			app.Div().
				Style("display", "flex").
				Style("align-items", "center").
				Style("gap", "16px").
				Style("margin-bottom", "20px").
				Body(
					app.Button().
						Style("background", "rgba(77,41,117,0.2)").
						Style("border", "1px solid rgba(77,41,117,0.4)").
						Style("border-radius", "8px").
						Style("padding", "8px 16px").
						Style("color", "#9B8DB5").
						Style("font-size", "13px").
						Style("cursor", "pointer").
						OnClick(p.onBack).
						Text("← Kembali"),
					app.H1().
						Style("color", "#E2D9F3").
						Style("font-size", "22px").
						Style("font-weight", "700").
						Style("margin", "0").
						Text("Detail License"),
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

			// Success message
			app.If(p.successMsg != "",
				func() app.UI {
					return app.Div().
						Style("background", "rgba(34,197,94,0.12)").
						Style("border", "1px solid rgba(34,197,94,0.4)").
						Style("border-radius", "8px").
						Style("padding", "12px 16px").
						Style("color", "#22C55E").
						Style("font-size", "14px").
						Style("margin-bottom", "16px").
						Style("display", "flex").
						Style("align-items", "center").
						Style("gap", "8px").
						Body(
							app.Raw(`<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="20 6 9 17 4 12"/></svg>`),
							app.Span().Text(p.successMsg),
						)
				},
			),

			// Loading state
			app.If(p.loading,
				func() app.UI {
					return app.Div().
						Style("text-align", "center").
						Style("padding", "60px").
						Style("color", "#9B8DB5").
						Text("Memuat data license...")
				},
			),

			// Suspend confirm dialog
			app.If(p.showSuspendConfirm,
				func() app.UI {
					return p.renderSuspendConfirmDialog()
				},
			),

			// Set status confirm dialog (superuser)
			app.If(p.showSetStatusConfirm,
				func() app.UI {
					return p.renderSetStatusConfirmDialog()
				},
			),

			// Credential dialog (shared: activate, reset-superuser)
			app.If(p.showCredDialog,
				func() app.UI {
					return p.renderCredDialog()
				},
			),

			// Header card — license key + status + action (visible sebelum tabs)
			app.If(!p.loading && p.license != nil,
				func() app.UI {
					return p.renderLicenseHeader()
				},
			),

			// Content tabs (only when loaded)
			app.If(!p.loading && p.license != nil,
				func() app.UI {
					return p.renderTabs()
				},
			),
		)
}

// renderLicenseHeader merender header card dengan license key, status badge, dan action buttons.
// Selalu tampil di atas tabs sehingga user dapat langsung mengubah status.
func (p *LicenseDetailPage) renderLicenseHeader() app.UI {
	l := p.license
	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "20px 24px").
		Style("margin-bottom", "20px").
		Body(
			app.Div().
				Style("display", "flex").
				Style("align-items", "flex-start").
				Style("justify-content", "space-between").
				Style("flex-wrap", "wrap").
				Style("gap", "12px").
				Body(
					// License key + company
					app.Div().
						Body(
							app.Div().
								Style("font-size", "11px").
								Style("color", "#9B8DB5").
								Style("text-transform", "uppercase").
								Style("letter-spacing", "0.08em").
								Style("margin-bottom", "6px").
								Text("License Key"),
							app.Div().
								Style("font-family", "monospace").
								Style("font-size", "20px").
								Style("font-weight", "700").
								Style("color", "#E2D9F3").
								Style("letter-spacing", "0.04em").
								Text(l.LicenseKey),
							app.Div().
								Style("font-size", "13px").
								Style("color", "#9B8DB5").
								Style("margin-top", "4px").
								Text(l.CompanyName+" · "+l.ProjectName),
						),
					// Status + action buttons
					app.Div().
						Style("display", "flex").
						Style("align-items", "center").
						Style("gap", "12px").
						Style("flex-wrap", "wrap").
						Body(
							statusBadge(l.Status),
							app.If(p.authStore.HasRole("project_owner"),
								func() app.UI {
									return p.renderActionButtons()
								},
							),
							app.If(p.authStore.HasRole("superuser"),
								func() app.UI {
									return p.renderSuperuserStatusDropdown()
								},
							),
						),
				),
		)
}

// renderSuspendConfirmDialog merender dialog konfirmasi sebelum suspend.
func (p *LicenseDetailPage) renderSuspendConfirmDialog() app.UI {
	return app.Div().
		Style("position", "fixed").
		Style("inset", "0").
		Style("background", "rgba(0,0,0,0.6)").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "center").
		Style("z-index", "1000").
		Body(
			app.Div().
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(239,68,68,0.4)").
				Style("border-radius", "12px").
				Style("padding", "28px 32px").
				Style("max-width", "440px").
				Style("width", "100%").
				Body(
					app.H2().
						Style("color", "#E2D9F3").
						Style("font-size", "18px").
						Style("font-weight", "700").
						Style("margin", "0 0 12px").
						Text("Konfirmasi Suspend License"),
					app.P().
						Style("color", "#9B8DB5").
						Style("font-size", "14px").
						Style("margin", "0 0 24px").
						Text("License akan disuspend. Validate call berikutnya akan mengembalikan false dan client app akan memblokir akses user. Lanjutkan?"),
					app.Div().
						Style("display", "flex").
						Style("gap", "12px").
						Style("justify-content", "flex-end").
						Body(
							app.Button().
								Style("background", "rgba(77,41,117,0.2)").
								Style("border", "1px solid rgba(77,41,117,0.4)").
								Style("border-radius", "8px").
								Style("padding", "10px 20px").
								Style("color", "#9B8DB5").
								Style("font-size", "14px").
								Style("cursor", "pointer").
								OnClick(p.onSuspendCancel).
								Text("Batal"),
							app.Button().
								Style("background", "rgba(239,68,68,0.2)").
								Style("border", "1px solid rgba(239,68,68,0.5)").
								Style("border-radius", "8px").
								Style("padding", "10px 20px").
								Style("color", "#EF4444").
								Style("font-size", "14px").
								Style("font-weight", "600").
								Style("cursor", "pointer").
								OnClick(p.onSuspendConfirm).
								Text("Suspend"),
						),
				),
		)
}

// renderTabs merender tab bar dan konten tab aktif.
func (p *LicenseDetailPage) renderTabs() app.UI {
	tabs := []struct {
		key   string
		label string
	}{
		{"info", "Info"},
		{"registration", "Registration Status"},
		{"activity", "Activity"},
	}

	tabItems := make([]app.UI, 0, len(tabs))
	for _, t := range tabs {
		t := t
		isActive := p.activeTab == t.key
		bg := "transparent"
		color := "#9B8DB5"
		borderBottom := "2px solid transparent"
		if isActive {
			color = "#E2D9F3"
			borderBottom = "2px solid #4D2975"
		}
		tabItems = append(tabItems, app.Button().
			Style("background", bg).
			Style("border", "none").
			Style("border-bottom", borderBottom).
			Style("padding", "10px 20px").
			Style("color", color).
			Style("font-size", "14px").
			Style("font-weight", func() string {
				if isActive {
					return "600"
				}
				return "400"
			}()).
			Style("cursor", "pointer").
			Style("transition", "color 0.15s").
			OnClick(p.onTabClick(t.key)).
			Text(t.label),
		)
	}

	var tabContent app.UI
	switch p.activeTab {
	case "registration":
		tabContent = p.renderRegistrationTab()
	case "activity":
		tabContent = p.renderActivityTab()
	default:
		tabContent = p.renderInfoTab()
	}

	return app.Div().Body(
		// Tab bar
		app.Div().
			Style("display", "flex").
			Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
			Style("margin-bottom", "24px").
			Body(tabItems...),

		// Tab content
		tabContent,
	)
}

// renderInfoTab merender konten tab Info.
func (p *LicenseDetailPage) renderInfoTab() app.UI {
	l := p.license

	// Details grid
	detailGrid := app.Div().
		Style("display", "grid").
		Style("grid-template-columns", "repeat(2, 1fr)").
		Style("gap", "16px").
		Style("margin-bottom", "20px").
		Body(
			infoCard("Plan", l.Plan),
			infoCard("Product", l.ProductName),
			infoCard("Company", l.CompanyName),
			infoCard("Project", l.ProjectName),
			infoCardOptionalFloat("Contract Amount", l.ContractAmount),
			infoCardOptionalStr("Expires At", l.ExpiresAt),
		)

	// Constraints grid
	constraintsTitle := app.H3().
		Style("color", "#E2D9F3").
		Style("font-size", "15px").
		Style("font-weight", "600").
		Style("margin", "0 0 12px").
		Text("Constraints")

	constraintsGrid := app.Div().
		Style("display", "grid").
		Style("grid-template-columns", "repeat(3, 1fr)").
		Style("gap", "12px").
		Style("margin-bottom", "20px").
		Body(
			infoCardOptionalInt("Max Users", l.MaxUsers),
			infoCardOptionalInt("Max Trans/Bulan", l.MaxTransPerMonth),
			infoCardOptionalInt("Max Trans/Hari", l.MaxTransPerDay),
			infoCardOptionalInt("Max Items", l.MaxItems),
			infoCardOptionalInt("Max Customers", l.MaxCustomers),
			infoCardOptionalInt("Max Branches", l.MaxBranches),
			infoCardOptionalInt("Max Storage (MB)", l.MaxStorage),
		)

	// Modules and Apps chips
	modulesSection := app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "16px 20px").
		Style("margin-bottom", "16px").
		Body(
			app.Div().
				Style("font-size", "12px").
				Style("color", "#9B8DB5").
				Style("margin-bottom", "10px").
				Text("MODULES"),
			p.renderChips(l.Modules, "#4D2975", "rgba(77,41,117,0.2)"),
		)

	appsSection := app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "16px 20px").
		Body(
			app.Div().
				Style("font-size", "12px").
				Style("color", "#9B8DB5").
				Style("margin-bottom", "10px").
				Text("APPS"),
			p.renderChips(l.Apps, "#26B8B0", "rgba(38,184,176,0.15)"),
		)

	return app.Div().Body(
		detailGrid,
		constraintsTitle,
		constraintsGrid,
		modulesSection,
		appsSection,
	)
}

// renderActionButtons merender tombol aksi berdasarkan status license.
func (p *LicenseDetailPage) renderActionButtons() app.UI {
	l := p.license

	resetBtn := app.Button().
		Style("background", "rgba(245,158,11,0.15)").
		Style("border", "1px solid rgba(245,158,11,0.4)").
		Style("border-radius", "8px").
		Style("padding", "7px 16px").
		Style("color", "#F59E0B").
		Style("font-size", "13px").
		Style("font-weight", "600").
		Style("cursor", "pointer").
		OnClick(p.onResetSuperuser).
		Text("Reset Superuser")

	switch l.Status {
	case "active":
		return app.Div().
			Style("display", "flex").
			Style("gap", "8px").
			Body(
				app.Button().
					Style("background", "rgba(239,68,68,0.15)").
					Style("border", "1px solid rgba(239,68,68,0.4)").
					Style("border-radius", "8px").
					Style("padding", "7px 16px").
					Style("color", "#EF4444").
					Style("font-size", "13px").
					Style("font-weight", "600").
					Style("cursor", "pointer").
					OnClick(p.onSuspendRequest).
					Text("Suspend"),
				resetBtn,
			)
	case "trial":
		return app.Div().
			Style("display", "flex").
			Style("gap", "8px").
			Body(
				app.Button().
					Style("background", "rgba(34,197,94,0.15)").
					Style("border", "1px solid rgba(34,197,94,0.4)").
					Style("border-radius", "8px").
					Style("padding", "7px 16px").
					Style("color", "#22C55E").
					Style("font-size", "13px").
					Style("font-weight", "600").
					Style("cursor", "pointer").
					OnClick(p.onActivate).
					Text("Activate"),
				resetBtn,
			)
	case "suspended", "pending":
		return app.Div().
			Style("display", "flex").
			Style("gap", "8px").
			Body(
				app.Button().
					Style("background", "rgba(34,197,94,0.15)").
					Style("border", "1px solid rgba(34,197,94,0.4)").
					Style("border-radius", "8px").
					Style("padding", "7px 16px").
					Style("color", "#22C55E").
					Style("font-size", "13px").
					Style("font-weight", "600").
					Style("cursor", "pointer").
					OnClick(p.onActivate).
					Text("Activate"),
			)
	case "expired":
		return app.Div().
			Style("display", "flex").
			Style("gap", "8px").
			Body(
				app.Button().
					Style("background", "rgba(38,184,176,0.15)").
					Style("border", "1px solid rgba(38,184,176,0.4)").
					Style("border-radius", "8px").
					Style("padding", "7px 16px").
					Style("color", "#26B8B0").
					Style("font-size", "13px").
					Style("font-weight", "600").
					Style("cursor", "pointer").
					OnClick(p.onRenew).
					Text("Renew"),
			)
	}
	return app.Div()
}

// statusInfo mengembalikan label dan warna untuk sebuah status license.
func statusInfo(status string) (label, color string) {
	switch status {
	case "active":
		return "Active", "#22C55E"
	case "suspended":
		return "Suspended", "#EF4444"
	case "pending":
		return "Pending", "#F59E0B"
	case "trial":
		return "Trial", "#26B8B0"
	case "expired":
		return "Expired", "#9B8DB5"
	default:
		return status, "#9B8DB5"
	}
}

// renderSetStatusConfirmDialog merender dialog konfirmasi sebelum set status (superuser).
func (p *LicenseDetailPage) renderSetStatusConfirmDialog() app.UI {
	label, color := statusInfo(p.pendingStatus)
	btnText := "Ubah Status"
	if p.setStatusLoading {
		btnText = "Mengubah..."
	}
	return app.Div().
		Style("position", "fixed").
		Style("inset", "0").
		Style("background", "rgba(0,0,0,0.6)").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "center").
		Style("z-index", "1000").
		Body(
			app.Div().
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(77,41,117,0.5)").
				Style("border-radius", "12px").
				Style("padding", "28px 32px").
				Style("max-width", "400px").
				Style("width", "100%").
				Body(
					app.H2().
						Style("color", "#E2D9F3").
						Style("font-size", "18px").
						Style("font-weight", "700").
						Style("margin", "0 0 12px").
						Text("Konfirmasi Set Status"),
					app.P().
						Style("color", "#9B8DB5").
						Style("font-size", "14px").
						Style("margin", "0 0 20px").
						Text("Ubah status license menjadi:"),
					app.Div().
						Style("display", "inline-block").
						Style("background", fmt.Sprintf("rgba(%s,0.15)", colorToRGB(color))).
						Style("border", fmt.Sprintf("1px solid %s", color)).
						Style("border-radius", "8px").
						Style("padding", "6px 18px").
						Style("color", color).
						Style("font-size", "15px").
						Style("font-weight", "700").
						Style("margin-bottom", "24px").
						Text(label),
					app.Div().
						Style("display", "flex").
						Style("gap", "12px").
						Style("justify-content", "flex-end").
						Body(
							app.Button().
								Style("background", "rgba(77,41,117,0.2)").
								Style("border", "1px solid rgba(77,41,117,0.4)").
								Style("border-radius", "8px").
								Style("padding", "10px 20px").
								Style("color", "#9B8DB5").
								Style("font-size", "14px").
								Style("cursor", "pointer").
								Disabled(p.setStatusLoading).
								OnClick(p.onSetStatusCancelConfirm).
								Text("Batal"),
							app.Button().
								Style("background", fmt.Sprintf("rgba(%s,0.2)", colorToRGB(color))).
								Style("border", fmt.Sprintf("1px solid %s", color)).
								Style("border-radius", "8px").
								Style("padding", "10px 20px").
								Style("color", color).
								Style("font-size", "14px").
								Style("font-weight", "600").
								Style("cursor", "pointer").
								Disabled(p.setStatusLoading).
								OnClick(p.onSetStatusConfirm).
								Text(btnText),
						),
				),
		)
}

// colorToRGB mengkonversi hex color ke format "r,g,b" untuk rgba().
func colorToRGB(hex string) string {
	switch hex {
	case "#22C55E":
		return "34,197,94"
	case "#EF4444":
		return "239,68,68"
	case "#F59E0B":
		return "245,158,11"
	case "#26B8B0":
		return "38,184,176"
	case "#9B8DB5":
		return "155,141,181"
	default:
		return "77,41,117"
	}
}

// renderSuperuserStatusDropdown merender dropdown untuk superuser mengubah status secara langsung.
func (p *LicenseDetailPage) renderSuperuserStatusDropdown() app.UI {
	statuses := []struct {
		value string
		label string
		color string
	}{
		{"pending", "Pending", "#F59E0B"},
		{"active", "Active", "#22C55E"},
		{"trial", "Trial", "#26B8B0"},
		{"suspended", "Suspended", "#EF4444"},
		{"expired", "Expired", "#9B8DB5"},
	}

	items := make([]app.UI, 0, len(statuses))
	for _, s := range statuses {
		s := s
		items = append(items, app.Div().
			Style("padding", "8px 16px").
			Style("cursor", "pointer").
			Style("font-size", "13px").
			Style("color", s.color).
			Style("white-space", "nowrap").
			OnClick(p.onSetStatus(s.value)).
			Text(s.label),
		)
	}

	return app.Div().
		Style("position", "relative").
		Body(
			app.Button().
				Style("background", "rgba(77,41,117,0.2)").
				Style("border", "1px solid rgba(77,41,117,0.5)").
				Style("border-radius", "8px").
				Style("padding", "7px 16px").
				Style("color", "#9B8DB5").
				Style("font-size", "13px").
				Style("font-weight", "600").
				Style("cursor", "pointer").
				OnClick(p.onToggleStatusDropdown).
				Text("Set Status ▾"),
			app.If(p.showStatusDropdown,
				func() app.UI {
					return app.Div().
						Style("position", "absolute").
						Style("top", "calc(100% + 4px)").
						Style("right", "0").
						Style("background", "#1A1035").
						Style("border", "1px solid rgba(77,41,117,0.5)").
						Style("border-radius", "8px").
						Style("z-index", "100").
						Style("min-width", "140px").
						Style("overflow", "hidden").
						Body(items...)
				},
			),
		)
}

// renderCredDialog merender credential dialog yang digunakan untuk activate dan reset-superuser.
func (p *LicenseDetailPage) renderCredDialog() app.UI {
	var title, desc, btnText, borderColor, btnBg, btnBorder, btnColor string
	switch p.credAction {
	case "reset-superuser":
		title = "Reset Superuser"
		desc = "Masukkan username dan password baru untuk superuser di client app."
		btnText = "Reset Superuser"
		borderColor = "rgba(245,158,11,0.4)"
		btnBg = "rgba(245,158,11,0.2)"
		btnBorder = "rgba(245,158,11,0.5)"
		btnColor = "#F59E0B"
	default: // "activate"
		title = "Activate License"
		desc = "Masukkan username dan password untuk superuser di client app."
		btnText = "Activate"
		borderColor = "rgba(34,197,94,0.4)"
		btnBg = "rgba(34,197,94,0.2)"
		btnBorder = "rgba(34,197,94,0.5)"
		btnColor = "#22C55E"
	}
	if p.credLoading {
		btnText = "Memproses..."
	}
	return app.Div().
		Style("position", "fixed").
		Style("inset", "0").
		Style("background", "rgba(0,0,0,0.6)").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "center").
		Style("z-index", "1000").
		Body(
			app.Div().
				Style("background", "#1A1035").
				Style("border", "1px solid "+borderColor).
				Style("border-radius", "12px").
				Style("padding", "28px 32px").
				Style("max-width", "440px").
				Style("width", "100%").
				Body(
					app.H2().
						Style("color", "#E2D9F3").
						Style("font-size", "18px").
						Style("font-weight", "700").
						Style("margin", "0 0 8px").
						Text(title),
					app.P().
						Style("color", "#9B8DB5").
						Style("font-size", "14px").
						Style("margin", "0 0 20px").
						Text(desc),
					app.Div().
						Style("margin-bottom", "14px").
						Body(
							app.Label().
								Style("display", "block").
								Style("color", "#9B8DB5").
								Style("font-size", "12px").
								Style("margin-bottom", "6px").
								Text("Username"),
							app.Input().
								Type("text").
								Placeholder("superuser username").
								Style("width", "100%").
								Style("padding", "10px 14px").
								Style("background", "rgba(77,41,117,0.15)").
								Style("border", "1px solid rgba(77,41,117,0.4)").
								Style("border-radius", "8px").
								Style("color", "#E2D9F3").
								Style("font-size", "14px").
								Style("outline", "none").
								Style("box-sizing", "border-box").
								OnInput(p.onCredUsernameInput),
						),
					app.Div().
						Style("margin-bottom", "20px").
						Body(
							app.Label().
								Style("display", "block").
								Style("color", "#9B8DB5").
								Style("font-size", "12px").
								Style("margin-bottom", "6px").
								Text("Password"),
							renderPasswordInput("superuser password", p.credShowPassword, p.onCredPasswordInput, p.onCredTogglePassword),
						),
					app.Div().
						Style("display", "flex").
						Style("gap", "12px").
						Style("justify-content", "flex-end").
						Body(
							app.Button().
								Style("background", "rgba(77,41,117,0.2)").
								Style("border", "1px solid rgba(77,41,117,0.4)").
								Style("border-radius", "8px").
								Style("padding", "10px 20px").
								Style("color", "#9B8DB5").
								Style("font-size", "14px").
								Style("cursor", "pointer").
								Disabled(p.credLoading).
								OnClick(p.onCredCancel).
								Text("Batal"),
							app.Button().
								Style("background", btnBg).
								Style("border", "1px solid "+btnBorder).
								Style("border-radius", "8px").
								Style("padding", "10px 20px").
								Style("color", btnColor).
								Style("font-size", "14px").
								Style("font-weight", "600").
								Style("cursor", func() string {
									if p.credLoading {
										return "not-allowed"
									}
									return "pointer"
								}()).
								Disabled(p.credLoading).
								OnClick(p.onCredConfirm).
								Text(btnText),
						),
				),
		)
}



// renderRegistrationTab merender tab Registration Status.
func (p *LicenseDetailPage) renderRegistrationTab() app.UI {
	l := p.license

	regStatusLabel := "Belum Terdaftar"
	regStatusColor := "#F59E0B"
	regStatusBg := "rgba(245,158,11,0.15)"
	if l.IsRegistered {
		regStatusLabel = "Terdaftar"
		regStatusColor = "#22C55E"
		regStatusBg = "rgba(34,197,94,0.15)"
	}

	instanceURL := l.InstanceURL
	if instanceURL == "" {
		instanceURL = "—"
	}
	instanceName := l.InstanceName
	if instanceName == "" {
		instanceName = "—"
	}

	lastPullAt := "—"
	if l.LastPullAt != nil && *l.LastPullAt != "" {
		s := *l.LastPullAt
		if len(s) >= 19 {
			lastPullAt = s[:10] + " " + s[11:19]
		}
	}

	return app.Div().
		Style("display", "flex").
		Style("flex-direction", "column").
		Style("gap", "16px").
		Body(
			// Registration status badge
			app.Div().
				Style("display", "inline-flex").
				Style("align-items", "center").
				Style("gap", "8px").
				Style("background", regStatusBg).
				Style("border-radius", "8px").
				Style("padding", "10px 18px").
				Style("align-self", "flex-start").
				Body(
					app.Span().
						Style("display", "inline-block").
						Style("width", "8px").
						Style("height", "8px").
						Style("border-radius", "50%").
						Style("background", regStatusColor),
					app.Span().
						Style("color", regStatusColor).
						Style("font-size", "14px").
						Style("font-weight", "600").
						Text(regStatusLabel),
				),

			// Instance info grid
			app.Div().
				Style("display", "grid").
				Style("grid-template-columns", "repeat(2, 1fr)").
				Style("gap", "16px").
				Body(
					infoCard("Instance URL", instanceURL),
					infoCard("Instance Name", instanceName),
					infoCard("Client App IP", func() string {
						if l.ClientAppIP != "" {
							return l.ClientAppIP
						}
						return "—"
					}()),
					infoCard("Superuser", func() string {
						if l.SuperuserUsername != "" {
							return l.SuperuserUsername
						}
						return "—"
					}()),
					infoCard("Last Validate At", lastPullAt),
					infoCard("Check Interval", l.CheckInterval),
				),

		)
}

// renderActivityTab merender tab Activity (audit log timeline).
func (p *LicenseDetailPage) renderActivityTab() app.UI {
	// Show error message if failed to fetch
	if p.errMsg != "" {
		return app.Div().
			Style("background", "rgba(239,68,68,0.15)").
			Style("border", "1px solid rgba(239,68,68,0.4)").
			Style("border-radius", "8px").
			Style("padding", "20px 24px").
			Style("color", "#EF4444").
			Style("font-size", "14px").
			Style("text-align", "center").
			Text(p.errMsg)
	}
	if len(p.auditLogs) == 0 {
		return app.Div().
			Style("text-align", "center").
			Style("padding", "60px").
			Style("color", "#9B8DB5").
			Text("Belum ada aktivitas tercatat.")
	}

	items := make([]app.UI, 0, len(p.auditLogs))
	for _, log := range p.auditLogs {
		ts := log.CreatedAt
		if len(ts) >= 19 {
			ts = ts[:10] + " " + ts[11:19]
		}

		items = append(items, app.Div().
			Style("display", "flex").
			Style("gap", "16px").
			Style("margin-bottom", "16px").
			Body(
				// Timeline dot + line
				app.Div().
					Style("display", "flex").
					Style("flex-direction", "column").
					Style("align-items", "center").
					Style("gap", "0").
					Body(
						app.Div().
							Style("width", "10px").
							Style("height", "10px").
							Style("border-radius", "50%").
							Style("background", "#4D2975").
							Style("flex-shrink", "0").
							Style("margin-top", "4px"),
						app.Div().
							Style("width", "1px").
							Style("flex", "1").
							Style("background", "rgba(77,41,117,0.3)").
							Style("margin-top", "4px"),
					),
				// Log content
				app.Div().
					Style("flex", "1").
					Style("padding-bottom", "12px").
					Body(
						app.Div().
							Style("color", "#E2D9F3").
							Style("font-size", "14px").
							Style("font-weight", "600").
							Text(log.Action),
						app.Div().
							Style("color", "#9B8DB5").
							Style("font-size", "12px").
							Style("margin-top", "2px").
							Text("Oleh "+log.ActorName+" · "+ts),
					),
			),
		)
	}

	return app.Div().
		Style("display", "flex").
		Style("flex-direction", "column").
		Body(items...)
}

// renderChips merender slice string sebagai chip badges.
func (p *LicenseDetailPage) renderChips(items []string, color, bg string) app.UI {
	if len(items) == 0 {
		return app.Span().Style("color", "#9B8DB5").Style("font-size", "13px").Text("—")
	}
	chips := make([]app.UI, 0, len(items))
	for _, item := range items {
		chips = append(chips, app.Span().
			Style("display", "inline-block").
			Style("background", bg).
			Style("color", color).
			Style("border-radius", "6px").
			Style("padding", "4px 10px").
			Style("font-size", "12px").
			Style("font-weight", "500").
			Style("margin", "2px").
			Text(item),
		)
	}
	return app.Div().
		Style("display", "flex").
		Style("flex-wrap", "wrap").
		Style("gap", "6px").
		Body(chips...)
}

// infoCard merender card label+value standar.
func infoCard(label, value string) app.UI {
	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "10px").
		Style("padding", "14px 16px").
		Body(
			app.Div().
				Style("font-size", "11px").
				Style("color", "#9B8DB5").
				Style("text-transform", "uppercase").
				Style("letter-spacing", "0.08em").
				Style("margin-bottom", "6px").
				Text(label),
			app.Div().
				Style("font-size", "14px").
				Style("color", "#E2D9F3").
				Style("font-weight", "500").
				Text(value),
		)
}

// infoCardOptionalStr merender card dengan nilai *string (menampilkan "—" jika nil).
func infoCardOptionalStr(label string, value *string) app.UI {
	v := "—"
	if value != nil && *value != "" {
		s := *value
		if len(s) >= 10 {
			v = s[:10]
		} else {
			v = s
		}
	}
	return infoCard(label, v)
}

// infoCardOptionalInt merender card dengan nilai *int (menampilkan "—" jika nil).
func infoCardOptionalInt(label string, value *int) app.UI {
	v := "—"
	if value != nil {
		v = fmt.Sprintf("%d", *value)
	}
	return infoCard(label, v)
}

// infoCardOptionalFloat merender card dengan nilai *float64.
func infoCardOptionalFloat(label string, value *float64) app.UI {
	v := "—"
	if value != nil {
		v = fmt.Sprintf("%.2f", *value)
	}
	return infoCard(label, v)
}
