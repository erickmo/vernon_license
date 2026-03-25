//go:build wasm

package pages

import (
	"context"
	"fmt"
	"strings"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/components"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// CompanyOption adalah representasi ringkas company untuk dropdown.
type CompanyOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ProjectOption adalah representasi ringkas project untuk dropdown.
type ProjectOption struct {
	ID        string `json:"id"`
	CompanyID string `json:"company_id"`
	Name      string `json:"name"`
}

// LicenseCreatePage adalah wizard 3 step untuk PO/superuser membuat license langsung.
// Step 1: Basic info, Step 2: Constraints, Step 3: Review + submit.
type LicenseCreatePage struct {
	app.Compo
	step int // 1, 2, 3

	// Step 1 fields
	companyID string
	projectID string
	productID string
	plan      string

	// Step 2 fields
	maxUsers         string
	maxTransPerMonth string
	maxTransPerDay   string
	maxItems         string
	maxCustomers     string
	maxBranches      string
	maxStorage       string
	expiresAt        string
	checkInterval    string

	// Dropdown data
	products  []ProductOption
	companies []CompanyOption
	projects  []ProjectOption // filtered by selected company

	// Post-creation state
	createdLicenseID              string
	createdOTP string
	createdLicenseKey             string

	loading    bool
	submitting bool
	errMsg     string
	authStore  store.AuthStore
}

// licenseCreateRequest adalah payload untuk POST /api/internal/licenses.
type licenseCreateRequest struct {
	CompanyID        string   `json:"company_id"`
	ProjectID        string   `json:"project_id"`
	ProductID        string   `json:"product_id"`
	Plan             string   `json:"plan"`
	Modules          []string `json:"modules"`
	Apps             []string `json:"apps"`
	MaxUsers         *int     `json:"max_users,omitempty"`
	MaxTransPerMonth *int     `json:"max_trans_per_month,omitempty"`
	MaxTransPerDay   *int     `json:"max_trans_per_day,omitempty"`
	MaxItems         *int     `json:"max_items,omitempty"`
	MaxCustomers     *int     `json:"max_customers,omitempty"`
	MaxBranches      *int     `json:"max_branches,omitempty"`
	MaxStorage       *int     `json:"max_storage,omitempty"`
	ExpiresAt        *string  `json:"expires_at,omitempty"`
	CheckInterval    string   `json:"check_interval"`
}

// licenseCreateResponse adalah response dari POST /api/internal/licenses.
type licenseCreateResponse struct {
	ID                     string `json:"id"`
	LicenseKey             string `json:"license_key"`
	OTP string `json:"otp"`
}

// OnNav dipanggil saat halaman ini di-navigasi.
// Redirect ke /login jika belum login, redirect ke /licenses jika bukan project_owner.
func (p *LicenseCreatePage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}
	if !p.authStore.HasRole("project_owner") {
		ctx.Navigate("/licenses")
		return
	}
	p.step = 1
	p.createdLicenseID = ""
	p.createdOTP = ""
	p.createdLicenseKey = ""
	p.errMsg = ""
	p.checkInterval = "6h"
	p.fetchDropdownData(ctx)
}

// fetchDropdownData mengambil data companies dan products untuk dropdown.
func (p *LicenseCreatePage) fetchDropdownData(ctx app.Context) {
	p.loading = true

	token := p.authStore.GetToken()
	ctx.Async(func() {
		client := api.NewClient("", token)

		var companies []CompanyOption
		_ = client.Get(context.Background(), "/api/internal/companies", &companies)

		var rawProducts []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		_ = client.Get(context.Background(), "/api/internal/products", &rawProducts)

		products := make([]ProductOption, 0, len(rawProducts))
		for _, rp := range rawProducts {
			products = append(products, ProductOption{ID: rp.ID, Name: rp.Name})
		}

		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			p.companies = companies
			p.products = products
		})
	})
}

// fetchProjects mengambil projects berdasarkan company yang dipilih.
func (p *LicenseCreatePage) fetchProjects(ctx app.Context, companyID string) {
	if companyID == "" {
		p.projects = nil
		return
	}
	token := p.authStore.GetToken()
	ctx.Async(func() {
		client := api.NewClient("", token)
		var projects []ProjectOption
		_ = client.Get(context.Background(), "/api/internal/companies/"+companyID+"/projects", &projects)
		ctx.Dispatch(func(ctx app.Context) {
			p.projects = projects
			p.projectID = ""
		})
	})
}

// Field change handlers

func (p *LicenseCreatePage) onCompanyChange(ctx app.Context, e app.Event) {
	p.companyID = ctx.JSSrc().Get("value").String()
	p.fetchProjects(ctx, p.companyID)
}

func (p *LicenseCreatePage) onProjectChange(ctx app.Context, e app.Event) {
	p.projectID = ctx.JSSrc().Get("value").String()
}

func (p *LicenseCreatePage) onProductChange(ctx app.Context, e app.Event) {
	p.productID = ctx.JSSrc().Get("value").String()
}

func (p *LicenseCreatePage) onPlanChange(plan string) func(app.Context, app.Event) {
	return func(ctx app.Context, e app.Event) {
		p.plan = plan
	}
}

func (p *LicenseCreatePage) onMaxUsersChange(ctx app.Context, e app.Event) {
	p.maxUsers = ctx.JSSrc().Get("value").String()
}
func (p *LicenseCreatePage) onMaxTransPerMonthChange(ctx app.Context, e app.Event) {
	p.maxTransPerMonth = ctx.JSSrc().Get("value").String()
}
func (p *LicenseCreatePage) onMaxTransPerDayChange(ctx app.Context, e app.Event) {
	p.maxTransPerDay = ctx.JSSrc().Get("value").String()
}
func (p *LicenseCreatePage) onMaxItemsChange(ctx app.Context, e app.Event) {
	p.maxItems = ctx.JSSrc().Get("value").String()
}
func (p *LicenseCreatePage) onMaxCustomersChange(ctx app.Context, e app.Event) {
	p.maxCustomers = ctx.JSSrc().Get("value").String()
}
func (p *LicenseCreatePage) onMaxBranchesChange(ctx app.Context, e app.Event) {
	p.maxBranches = ctx.JSSrc().Get("value").String()
}
func (p *LicenseCreatePage) onMaxStorageChange(ctx app.Context, e app.Event) {
	p.maxStorage = ctx.JSSrc().Get("value").String()
}
func (p *LicenseCreatePage) onExpiresAtChange(ctx app.Context, e app.Event) {
	p.expiresAt = ctx.JSSrc().Get("value").String()
}
func (p *LicenseCreatePage) onCheckIntervalChange(ctx app.Context, e app.Event) {
	p.checkInterval = ctx.JSSrc().Get("value").String()
}

// Navigation between steps

func (p *LicenseCreatePage) onNextStep(ctx app.Context, e app.Event) {
	if p.step == 1 {
		if p.companyID == "" || p.projectID == "" || p.productID == "" || p.plan == "" {
			p.errMsg = "Company, Project, Product, dan Plan wajib diisi."
			return
		}
	}
	p.errMsg = ""
	p.step++
}

func (p *LicenseCreatePage) onPrevStep(ctx app.Context, e app.Event) {
	if p.step > 1 {
		p.step--
		p.errMsg = ""
	}
}

func (p *LicenseCreatePage) onCancel(ctx app.Context, e app.Event) {
	ctx.Navigate("/licenses")
}

// onSubmit mengirim request buat license ke API.
func (p *LicenseCreatePage) onSubmit(ctx app.Context, e app.Event) {
	p.submitting = true
	p.errMsg = ""

	req := licenseCreateRequest{
		CompanyID:     p.companyID,
		ProjectID:     p.projectID,
		ProductID:     p.productID,
		Plan:          p.plan,
		Modules:       []string{},
		Apps:          []string{},
		CheckInterval: p.checkInterval,
	}

	// Parse optional integer fields
	if v := parseIntStr(p.maxUsers); v != nil {
		req.MaxUsers = v
	}
	if v := parseIntStr(p.maxTransPerMonth); v != nil {
		req.MaxTransPerMonth = v
	}
	if v := parseIntStr(p.maxTransPerDay); v != nil {
		req.MaxTransPerDay = v
	}
	if v := parseIntStr(p.maxItems); v != nil {
		req.MaxItems = v
	}
	if v := parseIntStr(p.maxCustomers); v != nil {
		req.MaxCustomers = v
	}
	if v := parseIntStr(p.maxBranches); v != nil {
		req.MaxBranches = v
	}
	if v := parseIntStr(p.maxStorage); v != nil {
		req.MaxStorage = v
	}
	if p.expiresAt != "" {
		// API accepts YYYY-MM-DD format
		ea := p.expiresAt
		req.ExpiresAt = &ea
	}

	if req.CheckInterval == "" {
		req.CheckInterval = "6h"
	}

	token := p.authStore.GetToken()
	ctx.Async(func() {
		client := api.NewClient("", token)
		var resp licenseCreateResponse
		if err := client.Post(context.Background(), "/api/internal/licenses", req, &resp); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.submitting = false
				p.errMsg = err.Error()
			})
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.submitting = false
			p.createdLicenseID = resp.ID
			p.createdOTP = resp.OTP
			p.createdLicenseKey = resp.LicenseKey
		})
	})
}

// onCopyCreatedKey menyalin OTP hasil create ke clipboard.
func (p *LicenseCreatePage) onCopyCreatedKey(ctx app.Context, e app.Event) {
	app.Window().Get("navigator").Get("clipboard").Call("writeText", p.createdOTP)
}

// onGoToDetail navigasi ke halaman detail license yang baru dibuat.
func (p *LicenseCreatePage) onGoToDetail(ctx app.Context, e app.Event) {
	ctx.Navigate("/licenses/" + p.createdLicenseID)
}

// Render menampilkan halaman buat license dalam Shell.
func (p *LicenseCreatePage) Render() app.UI {
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

// renderContent merender area utama halaman create license.
func (p *LicenseCreatePage) renderContent() app.UI {
	// Show success state after creation
	if p.createdLicenseID != "" {
		return p.renderSuccess()
	}

	return app.Div().
		Style("padding", "32px").
		Style("max-width", "720px").
		Body(
			// Header
			app.H1().
				Style("color", "#E2D9F3").
				Style("font-size", "22px").
				Style("font-weight", "700").
				Style("margin", "0 0 4px").
				Text("Buat License Baru"),
			app.P().
				Style("color", "#9B8DB5").
				Style("font-size", "14px").
				Style("margin", "0 0 24px").
				Text("Buat license langsung tanpa melalui proposal."),

			// Step indicator
			p.renderStepIndicator(),

			// Error
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

			// Loading
			app.If(p.loading,
				func() app.UI {
					return app.Div().
						Style("text-align", "center").
						Style("padding", "40px").
						Style("color", "#9B8DB5").
						Text("Memuat data...")
				},
			),

			// Step content
			app.If(!p.loading,
				func() app.UI {
					switch p.step {
					case 2:
						return p.renderStep2()
					case 3:
						return p.renderStep3()
					default:
						return p.renderStep1()
					}
				},
			),
		)
}

// renderStepIndicator merender indikator step wizard.
func (p *LicenseCreatePage) renderStepIndicator() app.UI {
	steps := []string{"Basic", "Constraints", "Review"}
	items := make([]app.UI, 0, len(steps)*2-1)
	for i, s := range steps {
		num := i + 1
		isActive := num == p.step
		isDone := num < p.step
		color := "#9B8DB5"
		bg := "rgba(155,141,181,0.15)"
		if isActive {
			color = "#E2D9F3"
			bg = "#4D2975"
		} else if isDone {
			color = "#22C55E"
			bg = "rgba(34,197,94,0.2)"
		}
		items = append(items, app.Div().
			Style("display", "flex").
			Style("align-items", "center").
			Style("gap", "8px").
			Body(
				app.Div().
					Style("width", "28px").
					Style("height", "28px").
					Style("border-radius", "50%").
					Style("background", bg).
					Style("color", color).
					Style("display", "flex").
					Style("align-items", "center").
					Style("justify-content", "center").
					Style("font-size", "13px").
					Style("font-weight", "700").
					Text(func() string {
						if isDone {
							return "✓"
						}
						return fmt.Sprintf("%d", num)
					}()),
				app.Span().
					Style("color", color).
					Style("font-size", "13px").
					Style("font-weight", func() string {
						if isActive {
							return "600"
						}
						return "400"
					}()).
					Text(s),
			),
		)
		if i < len(steps)-1 {
			items = append(items, app.Div().
				Style("flex", "1").
				Style("height", "1px").
				Style("background", "rgba(77,41,117,0.3)").
				Style("margin", "0 8px"),
			)
		}
	}

	return app.Div().
		Style("display", "flex").
		Style("align-items", "center").
		Style("margin-bottom", "28px").
		Style("padding", "16px 20px").
		Style("background", "#1A1035").
		Style("border-radius", "10px").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Body(items...)
}

// renderStep1 merender form step 1: Basic Info.
func (p *LicenseCreatePage) renderStep1() app.UI {
	// Build company options
	companyOpts := make([]app.UI, 0, len(p.companies)+1)
	companyOpts = append(companyOpts, app.Option().Value("").Text("Pilih Company..."))
	for _, c := range p.companies {
		companyOpts = append(companyOpts, app.Option().Value(c.ID).Selected(c.ID == p.companyID).Text(c.Name))
	}

	// Build project options (filtered by company)
	projectOpts := make([]app.UI, 0, len(p.projects)+1)
	projectOpts = append(projectOpts, app.Option().Value("").Text("Pilih Project..."))
	for _, pr := range p.projects {
		projectOpts = append(projectOpts, app.Option().Value(pr.ID).Selected(pr.ID == p.projectID).Text(pr.Name))
	}

	// Build product options
	productOpts := make([]app.UI, 0, len(p.products)+1)
	productOpts = append(productOpts, app.Option().Value("").Text("Pilih Product..."))
	for _, prod := range p.products {
		productOpts = append(productOpts, app.Option().Value(prod.ID).Selected(prod.ID == p.productID).Text(prod.Name))
	}

	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "24px").
		Body(
			app.H2().
				Style("color", "#E2D9F3").
				Style("font-size", "16px").
				Style("font-weight", "600").
				Style("margin", "0 0 20px").
				Text("Step 1 — Basic Info"),

			formField("Company", app.Select().
				Style("width", "100%").
				Style("background", "#0F0A1A").
				Style("border", "1px solid rgba(77,41,117,0.4)").
				Style("border-radius", "8px").
				Style("padding", "10px 14px").
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				OnChange(p.onCompanyChange).
				Body(companyOpts...),
			),

			formField("Project", app.Select().
				Style("width", "100%").
				Style("background", "#0F0A1A").
				Style("border", "1px solid rgba(77,41,117,0.4)").
				Style("border-radius", "8px").
				Style("padding", "10px 14px").
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				OnChange(p.onProjectChange).
				Body(projectOpts...),
			),

			formField("Product", app.Select().
				Style("width", "100%").
				Style("background", "#0F0A1A").
				Style("border", "1px solid rgba(77,41,117,0.4)").
				Style("border-radius", "8px").
				Style("padding", "10px 14px").
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				OnChange(p.onProductChange).
				Body(productOpts...),
			),

			// Plan radio buttons
			app.Div().
				Style("margin-bottom", "20px").
				Body(
					app.Label().
						Style("display", "block").
						Style("font-size", "12px").
						Style("color", "#9B8DB5").
						Style("text-transform", "uppercase").
						Style("letter-spacing", "0.06em").
						Style("margin-bottom", "10px").
						Text("Plan"),
					app.Div().
						Style("display", "flex").
						Style("gap", "12px").
						Body(
							planRadio("SaaS", p.plan == "SaaS" || p.plan == "saas", p.onPlanChange("SaaS")),
							planRadio("Dedicated", p.plan == "Dedicated" || p.plan == "dedicated", p.onPlanChange("Dedicated")),
						),
				),

			// Navigation buttons
			p.renderNavButtons(false),
		)
}

// renderStep2 merender form step 2: Constraints.
func (p *LicenseCreatePage) renderStep2() app.UI {
	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "24px").
		Body(
			app.H2().
				Style("color", "#E2D9F3").
				Style("font-size", "16px").
				Style("font-weight", "600").
				Style("margin", "0 0 20px").
				Text("Step 2 — Constraints"),

			app.Div().
				Style("display", "grid").
				Style("grid-template-columns", "repeat(2, 1fr)").
				Style("gap", "16px").
				Body(
					formField("Max Users", numberInput("Tanpa batas jika kosong", p.maxUsers, p.onMaxUsersChange)),
					formField("Max Trans/Bulan", numberInput("Tanpa batas jika kosong", p.maxTransPerMonth, p.onMaxTransPerMonthChange)),
					formField("Max Trans/Hari", numberInput("Tanpa batas jika kosong", p.maxTransPerDay, p.onMaxTransPerDayChange)),
					formField("Max Items", numberInput("Tanpa batas jika kosong", p.maxItems, p.onMaxItemsChange)),
					formField("Max Customers", numberInput("Tanpa batas jika kosong", p.maxCustomers, p.onMaxCustomersChange)),
					formField("Max Branches", numberInput("Tanpa batas jika kosong", p.maxBranches, p.onMaxBranchesChange)),
					formField("Max Storage (MB)", numberInput("Tanpa batas jika kosong", p.maxStorage, p.onMaxStorageChange)),
					formField("Expires At", app.Input().
						Type("date").
						Value(p.expiresAt).
						Style("width", "100%").
						Style("background", "#0F0A1A").
						Style("border", "1px solid rgba(77,41,117,0.4)").
						Style("border-radius", "8px").
						Style("padding", "10px 14px").
						Style("color", "#E2D9F3").
						Style("font-size", "14px").
						OnChange(p.onExpiresAtChange),
					),
				),

			formField("Check Interval", app.Select().
				Style("width", "100%").
				Style("background", "#0F0A1A").
				Style("border", "1px solid rgba(77,41,117,0.4)").
				Style("border-radius", "8px").
				Style("padding", "10px 14px").
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				OnChange(p.onCheckIntervalChange).
				Body(
					app.Option().Value("1h").Selected(p.checkInterval == "1h").Text("Setiap 1 jam"),
					app.Option().Value("6h").Selected(p.checkInterval == "6h" || p.checkInterval == "").Text("Setiap 6 jam (default)"),
					app.Option().Value("24h").Selected(p.checkInterval == "24h").Text("Setiap 24 jam"),
				),
			),

			p.renderNavButtons(true),
		)
}

// renderStep3 merender review sebelum submit.
func (p *LicenseCreatePage) renderStep3() app.UI {
	// Find display names
	companyName := p.companyID
	for _, c := range p.companies {
		if c.ID == p.companyID {
			companyName = c.Name
			break
		}
	}
	projectName := p.projectID
	for _, pr := range p.projects {
		if pr.ID == p.projectID {
			projectName = pr.Name
			break
		}
	}
	productName := p.productID
	for _, prod := range p.products {
		if prod.ID == p.productID {
			productName = prod.Name
			break
		}
	}

	expiresAtDisplay := "—"
	if p.expiresAt != "" {
		expiresAtDisplay = p.expiresAt
	}

	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "24px").
		Body(
			app.H2().
				Style("color", "#E2D9F3").
				Style("font-size", "16px").
				Style("font-weight", "600").
				Style("margin", "0 0 20px").
				Text("Step 3 — Review"),

			// Summary section
			app.Div().
				Style("display", "grid").
				Style("grid-template-columns", "repeat(2, 1fr)").
				Style("gap", "12px").
				Style("margin-bottom", "20px").
				Body(
					infoCard("Company", companyName),
					infoCard("Project", projectName),
					infoCard("Product", productName),
					infoCard("Plan", p.plan),
					infoCard("Check Interval", p.checkInterval),
					infoCard("Expires At", expiresAtDisplay),
				),

			// Constraints summary
			app.H3().
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				Style("font-weight", "600").
				Style("margin", "0 0 12px").
				Text("Constraints"),

			app.Div().
				Style("display", "grid").
				Style("grid-template-columns", "repeat(3, 1fr)").
				Style("gap", "10px").
				Style("margin-bottom", "24px").
				Body(
					infoCard("Max Users", orDash(p.maxUsers)),
					infoCard("Max Trans/Bulan", orDash(p.maxTransPerMonth)),
					infoCard("Max Trans/Hari", orDash(p.maxTransPerDay)),
					infoCard("Max Items", orDash(p.maxItems)),
					infoCard("Max Customers", orDash(p.maxCustomers)),
					infoCard("Max Branches", orDash(p.maxBranches)),
					infoCard("Max Storage (MB)", orDash(p.maxStorage)),
				),

			// Action buttons
			app.Div().
				Style("display", "flex").
				Style("gap", "12px").
				Style("justify-content", "space-between").
				Body(
					app.Button().
						Style("background", "rgba(77,41,117,0.2)").
						Style("border", "1px solid rgba(77,41,117,0.4)").
						Style("border-radius", "8px").
						Style("padding", "10px 20px").
						Style("color", "#9B8DB5").
						Style("font-size", "14px").
						Style("cursor", "pointer").
						OnClick(p.onPrevStep).
						Text("← Kembali"),
					app.Button().
						Style("background", "linear-gradient(135deg, #4D2975, #26B8B0)").
						Style("border", "none").
						Style("border-radius", "8px").
						Style("padding", "10px 28px").
						Style("color", "#E2D9F3").
						Style("font-size", "14px").
						Style("font-weight", "600").
						Style("cursor", "pointer").
						Disabled(p.submitting).
						OnClick(p.onSubmit).
						Text(func() string {
							if p.submitting {
								return "Membuat..."
							}
							return "Buat License"
						}()),
				),
		)
}

// renderNavButtons merender tombol navigasi antar step.
func (p *LicenseCreatePage) renderNavButtons(showBack bool) app.UI {
	return app.Div().
		Style("display", "flex").
		Style("gap", "12px").
		Style("justify-content", "space-between").
		Style("margin-top", "24px").
		Body(
			app.Div().
				Style("display", "flex").
				Style("gap", "12px").
				Body(
					app.If(showBack,
						func() app.UI {
							return app.Button().
								Style("background", "rgba(77,41,117,0.2)").
								Style("border", "1px solid rgba(77,41,117,0.4)").
								Style("border-radius", "8px").
								Style("padding", "10px 20px").
								Style("color", "#9B8DB5").
								Style("font-size", "14px").
								Style("cursor", "pointer").
								OnClick(p.onPrevStep).
								Text("← Kembali")
						},
					),
					app.Button().
						Style("background", "rgba(155,141,181,0.1)").
						Style("border", "1px solid rgba(155,141,181,0.3)").
						Style("border-radius", "8px").
						Style("padding", "10px 20px").
						Style("color", "#9B8DB5").
						Style("font-size", "14px").
						Style("cursor", "pointer").
						OnClick(p.onCancel).
						Text("Batal"),
				),
			app.Button().
				Style("background", "linear-gradient(135deg, #4D2975, #26B8B0)").
				Style("border", "none").
				Style("border-radius", "8px").
				Style("padding", "10px 24px").
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				Style("font-weight", "600").
				Style("cursor", "pointer").
				OnClick(p.onNextStep).
				Text("Lanjut →"),
		)
}

// renderSuccess merender halaman setelah license berhasil dibuat.
func (p *LicenseCreatePage) renderSuccess() app.UI {
	return app.Div().
		Style("padding", "32px").
		Style("max-width", "600px").
		Body(
			app.Div().
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(34,197,94,0.4)").
				Style("border-radius", "16px").
				Style("padding", "32px").
				Body(
					app.Div().
						Style("text-align", "center").
						Style("margin-bottom", "24px").
						Body(
							app.Div().
								Style("font-size", "40px").
								Style("margin-bottom", "8px").
								Text("✓"),
							app.H2().
								Style("color", "#22C55E").
								Style("font-size", "20px").
								Style("font-weight", "700").
								Style("margin", "0 0 8px").
								Text("License Berhasil Dibuat!"),
							app.P().
								Style("color", "#9B8DB5").
								Style("font-size", "14px").
								Style("margin", "0").
								Text("Salin OTP dan berikan ke tim client untuk melakukan registrasi."),
						),

					// License key display
					app.Div().
						Style("background", "rgba(77,41,117,0.2)").
						Style("border-radius", "8px").
						Style("padding", "14px 18px").
						Style("margin-bottom", "16px").
						Body(
							app.Div().
								Style("font-size", "11px").
								Style("color", "#9B8DB5").
								Style("text-transform", "uppercase").
								Style("margin-bottom", "6px").
								Text("License Key"),
							app.Div().
								Style("font-family", "monospace").
								Style("font-size", "18px").
								Style("font-weight", "700").
								Style("color", "#E2D9F3").
								Text(p.createdLicenseKey),
						),

					// OTP
					app.Div().
						Style("background", "rgba(38,184,176,0.1)").
						Style("border", "1px solid rgba(38,184,176,0.3)").
						Style("border-radius", "8px").
						Style("padding", "14px 18px").
						Style("margin-bottom", "24px").
						Body(
							app.Div().
								Style("font-size", "11px").
								Style("color", "#9B8DB5").
								Style("text-transform", "uppercase").
								Style("margin-bottom", "6px").
								Text("OTP"),
							app.Div().
								Style("display", "flex").
								Style("align-items", "center").
								Style("gap", "12px").
								Body(
									app.Span().
										Style("font-family", "monospace").
										Style("font-size", "14px").
										Style("color", "#26B8B0").
										Style("word-break", "break-all").
										Text(p.createdOTP),
									app.Button().
										Style("background", "rgba(38,184,176,0.15)").
										Style("border", "1px solid rgba(38,184,176,0.4)").
										Style("border-radius", "6px").
										Style("padding", "6px 14px").
										Style("color", "#26B8B0").
										Style("font-size", "12px").
										Style("font-weight", "600").
										Style("cursor", "pointer").
										Style("flex-shrink", "0").
										OnClick(p.onCopyCreatedKey).
										Text("Copy"),
								),
						),

					// Action buttons
					app.Div().
						Style("display", "flex").
						Style("gap", "12px").
						Style("justify-content", "center").
						Body(
							app.Button().
								Style("background", "rgba(77,41,117,0.2)").
								Style("border", "1px solid rgba(77,41,117,0.4)").
								Style("border-radius", "8px").
								Style("padding", "10px 20px").
								Style("color", "#9B8DB5").
								Style("font-size", "14px").
								Style("cursor", "pointer").
								OnClick(p.onCancel).
								Text("Kembali ke Daftar"),
							app.Button().
								Style("background", "linear-gradient(135deg, #4D2975, #26B8B0)").
								Style("border", "none").
								Style("border-radius", "8px").
								Style("padding", "10px 20px").
								Style("color", "#E2D9F3").
								Style("font-size", "14px").
								Style("font-weight", "600").
								Style("cursor", "pointer").
								OnClick(p.onGoToDetail).
								Text("Lihat Detail License"),
						),
				),
		)
}

// Helper UI functions

// formField merender label + input dalam satu group.
func formField(label string, input app.UI) app.UI {
	return app.Div().
		Style("margin-bottom", "16px").
		Body(
			app.Label().
				Style("display", "block").
				Style("font-size", "12px").
				Style("color", "#9B8DB5").
				Style("text-transform", "uppercase").
				Style("letter-spacing", "0.06em").
				Style("margin-bottom", "8px").
				Text(label),
			input,
		)
}

// numberInput merender input number dengan placeholder.
func numberInput(placeholder, value string, handler func(app.Context, app.Event)) app.UI {
	return app.Input().
		Type("number").
		Min("0").
		Placeholder(placeholder).
		Value(value).
		Style("width", "100%").
		Style("background", "#0F0A1A").
		Style("border", "1px solid rgba(77,41,117,0.4)").
		Style("border-radius", "8px").
		Style("padding", "10px 14px").
		Style("color", "#E2D9F3").
		Style("font-size", "14px").
		OnInput(handler)
}

// planRadio merender tombol radio style untuk pilihan plan.
func planRadio(label string, selected bool, handler func(app.Context, app.Event)) app.UI {
	bg := "rgba(77,41,117,0.15)"
	border := "1px solid rgba(77,41,117,0.3)"
	color := "#9B8DB5"
	if selected {
		bg = "rgba(77,41,117,0.4)"
		border = "1px solid #4D2975"
		color = "#E2D9F3"
	}
	return app.Button().
		Style("background", bg).
		Style("border", border).
		Style("border-radius", "8px").
		Style("padding", "10px 24px").
		Style("color", color).
		Style("font-size", "14px").
		Style("font-weight", "600").
		Style("cursor", "pointer").
		OnClick(handler).
		Text(label)
}

// orDash mengembalikan s atau "—" jika kosong.
func orDash(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "—"
	}
	return s
}

// parseIntStr mem-parse string ke *int. Returns nil jika tidak valid atau kosong.
func parseIntStr(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return nil
		}
		n = n*10 + int(c-'0')
	}
	return &n
}
