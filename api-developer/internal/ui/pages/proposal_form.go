//go:build wasm

package pages

import (
	"context"
	"fmt"
	"strings"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ProductOption adalah representasi ringkas product untuk dropdown.
type ProductOption struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	AvailablePlans   []string       `json:"available_plans"`
	AvailableModules []ModuleOption `json:"available_modules"`
	AvailableApps    []AppOption    `json:"available_apps"`
}

// ModuleOption adalah satu modul yang tersedia di suatu product.
type ModuleOption struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// AppOption adalah satu aplikasi yang tersedia di suatu product.
type AppOption struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// ProposalFormPage adalah form untuk create atau edit proposal.
type ProposalFormPage struct {
	app.Compo
	proposalID string // kosong = create mode
	projectID  string
	companyID  string

	// form fields
	productID        string
	plan             string
	selectedModules  map[string]bool
	selectedApps     map[string]bool
	maxUsers         string
	maxTransPerMonth string
	maxTransPerDay   string
	maxItems         string
	maxCustomers     string
	maxBranches      string
	maxStorage       string
	contractAmount   string
	expiresAt        string
	notes            string
	ownerNotes       string // hanya PO/superuser

	// state
	products   []ProductOption
	loading    bool
	saving     bool
	errMsg     string
	isEditMode bool
	authStore  store.AuthStore

	// live diff: versi tersimpan (untuk edit mode)
	savedVersion *ProposalDetail
}

// OnNav dipanggil saat navigasi ke halaman ini.
func (p *ProposalFormPage) OnNav(ctx app.Context) {
	urlPath := ctx.Page().URL().Path
	// /proposals/{id}/edit → edit mode
	// /proposals/create?project_id=...&company_id=... → create mode
	parts := strings.Split(strings.Trim(urlPath, "/"), "/")

	p.selectedModules = make(map[string]bool)
	p.selectedApps = make(map[string]bool)

	if len(parts) >= 3 && parts[2] == "edit" {
		p.proposalID = parts[1]
		p.isEditMode = true
	} else {
		p.proposalID = ""
		p.isEditMode = false
		// Baca query params
		q := ctx.Page().URL().Query()
		p.projectID = q.Get("project_id")
		p.companyID = q.Get("company_id")
	}

	p.loading = true

	go p.loadInitialData(ctx)
}

// loadInitialData memuat daftar products dan (jika edit) data proposal.
func (p *ProposalFormPage) loadInitialData(ctx app.Context) {
	user := p.authStore.GetUser()
	if user == nil {
		ctx.Dispatch(func(ctx app.Context) { ctx.Navigate("/login") })
		return
	}

	client := api.NewClient("", user.Token)

	// Load products
	var productsResp struct {
		Data []ProductOption `json:"data"`
	}
	if err := client.Get(context.Background(), "/api/internal/products", &productsResp); err != nil {
		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			p.errMsg = "Gagal memuat daftar produk: " + err.Error()
		})
		return
	}

	ctx.Dispatch(func(ctx app.Context) {
		p.products = productsResp.Data
	})

	// Jika edit mode, load data proposal
	if p.isEditMode && p.proposalID != "" {
		var proposal ProposalDetail
		if err := client.Get(context.Background(), "/api/internal/proposals/"+p.proposalID, &proposal); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.loading = false
				p.errMsg = "Gagal memuat data proposal: " + err.Error()
			})
			return
		}

		ctx.Dispatch(func(ctx app.Context) {
			p.savedVersion = &proposal
			p.projectID = proposal.ProjectID
			p.companyID = proposal.CompanyID
			p.productID = proposal.ProductID
			p.plan = proposal.Plan

			// Populate selected modules & apps
			p.selectedModules = make(map[string]bool)
			for _, m := range proposal.Modules {
				p.selectedModules[m] = true
			}
			p.selectedApps = make(map[string]bool)
			for _, a := range proposal.Apps {
				p.selectedApps[a] = true
			}

			// Constraints
			if proposal.MaxUsers != nil {
				p.maxUsers = fmt.Sprintf("%d", *proposal.MaxUsers)
			}
			if proposal.MaxTransPerMonth != nil {
				p.maxTransPerMonth = fmt.Sprintf("%d", *proposal.MaxTransPerMonth)
			}
			if proposal.MaxTransPerDay != nil {
				p.maxTransPerDay = fmt.Sprintf("%d", *proposal.MaxTransPerDay)
			}
			if proposal.MaxItems != nil {
				p.maxItems = fmt.Sprintf("%d", *proposal.MaxItems)
			}
			if proposal.MaxCustomers != nil {
				p.maxCustomers = fmt.Sprintf("%d", *proposal.MaxCustomers)
			}
			if proposal.MaxBranches != nil {
				p.maxBranches = fmt.Sprintf("%d", *proposal.MaxBranches)
			}
			if proposal.MaxStorage != nil {
				p.maxStorage = fmt.Sprintf("%d", *proposal.MaxStorage)
			}
			if proposal.ContractAmount != nil {
				p.contractAmount = fmt.Sprintf("%.2f", *proposal.ContractAmount)
			}
			if proposal.ExpiresAt != nil && len(*proposal.ExpiresAt) >= 10 {
				p.expiresAt = (*proposal.ExpiresAt)[:10]
			}
			p.notes = proposal.Notes
			p.ownerNotes = proposal.OwnerNotes

			p.loading = false
		})
		return
	}

	ctx.Dispatch(func(ctx app.Context) {
		p.loading = false
	})
}

// selectedProduct mengembalikan ProductOption yang dipilih, atau nil.
func (p *ProposalFormPage) selectedProduct() *ProductOption {
	for i := range p.products {
		if p.products[i].ID == p.productID {
			return &p.products[i]
		}
	}
	return nil
}

// changedFields mengembalikan daftar field yang berubah dari saved version (edit mode).
func (p *ProposalFormPage) changedFields() []string {
	if p.savedVersion == nil {
		return nil
	}
	sv := p.savedVersion
	var changed []string

	if p.plan != sv.Plan {
		changed = append(changed, "plan")
	}
	if p.contractAmount != formatAmount(sv.ContractAmount) {
		if !(p.contractAmount == "" && sv.ContractAmount == nil) {
			changed = append(changed, "contract_amount")
		}
	}
	if p.notes != sv.Notes {
		changed = append(changed, "notes")
	}
	if p.ownerNotes != sv.OwnerNotes {
		changed = append(changed, "owner_notes")
	}

	// Constraints
	if p.maxUsers != "" || sv.MaxUsers != nil {
		expected := ""
		if sv.MaxUsers != nil {
			expected = fmt.Sprintf("%d", *sv.MaxUsers)
		}
		if p.maxUsers != expected {
			changed = append(changed, "max_users")
		}
	}

	return changed
}

// onProductChange menangani pergantian product.
func (p *ProposalFormPage) onProductChange(ctx app.Context, e app.Event) {
	p.productID = ctx.JSSrc().Get("value").String()
	p.plan = ""
	p.selectedModules = make(map[string]bool)
	p.selectedApps = make(map[string]bool)
}

// onPlanChange menangani pergantian plan.
func (p *ProposalFormPage) onPlanChange(plan string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		p.plan = plan
	}
}

// onModuleToggle menangani toggle modul.
func (p *ProposalFormPage) onModuleToggle(key string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		if p.selectedModules == nil {
			p.selectedModules = make(map[string]bool)
		}
		p.selectedModules[key] = !p.selectedModules[key]
	}
}

// onAppToggle menangani toggle aplikasi.
func (p *ProposalFormPage) onAppToggle(key string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		if p.selectedApps == nil {
			p.selectedApps = make(map[string]bool)
		}
		p.selectedApps[key] = !p.selectedApps[key]
	}
}

// buildSavePayload membangun payload JSON untuk create/update.
func (p *ProposalFormPage) buildSavePayload() map[string]any {
	payload := map[string]any{
		"plan": p.plan,
	}

	if !p.isEditMode {
		payload["project_id"] = p.projectID
		payload["company_id"] = p.companyID
		payload["product_id"] = p.productID
	}

	// Modules & Apps
	modules := []string{}
	for k, v := range p.selectedModules {
		if v {
			modules = append(modules, k)
		}
	}
	payload["modules"] = modules

	apps := []string{}
	for k, v := range p.selectedApps {
		if v {
			apps = append(apps, k)
		}
	}
	payload["apps"] = apps

	// Notes
	if p.notes != "" {
		payload["notes"] = p.notes
	}
	if p.ownerNotes != "" {
		payload["owner_notes"] = p.ownerNotes
	}

	// Constraints (parseInt, skip jika kosong)
	setIntField := func(field, val string) {
		if val != "" {
			var n int
			if _, err := fmt.Sscanf(val, "%d", &n); err == nil {
				payload[field] = n
			}
		}
	}
	setIntField("max_users", p.maxUsers)
	setIntField("max_trans_per_month", p.maxTransPerMonth)
	setIntField("max_trans_per_day", p.maxTransPerDay)
	setIntField("max_items", p.maxItems)
	setIntField("max_customers", p.maxCustomers)
	setIntField("max_branches", p.maxBranches)
	setIntField("max_storage", p.maxStorage)

	// Contract amount
	if p.contractAmount != "" {
		var f float64
		if _, err := fmt.Sscanf(p.contractAmount, "%f", &f); err == nil {
			payload["contract_amount"] = f
		}
	}

	// Expires at
	if p.expiresAt != "" {
		payload["expires_at"] = p.expiresAt + "T00:00:00Z"
	}

	return payload
}

// onSave menyimpan proposal (create atau update).
func (p *ProposalFormPage) onSave(ctx app.Context, e app.Event) {
	if p.plan == "" {
		p.errMsg = "Plan wajib dipilih."
		return
	}
	if !p.isEditMode && p.productID == "" {
		p.errMsg = "Produk wajib dipilih."
		return
	}

	p.saving = true
	p.errMsg = ""

	payload := p.buildSavePayload()

	go func() {
		user := p.authStore.GetUser()
		if user == nil {
			ctx.Dispatch(func(ctx app.Context) { ctx.Navigate("/login") })
			return
		}

		client := api.NewClient("", user.Token)
		var result map[string]any

		var err error
		if p.isEditMode {
			err = client.Put(context.Background(), "/api/internal/proposals/"+p.proposalID, payload, &result)
		} else {
			err = client.Post(context.Background(), "/api/internal/proposals", payload, &result)
		}

		if err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.saving = false
				p.errMsg = "Gagal menyimpan: " + err.Error()
			})
			return
		}

		ctx.Dispatch(func(ctx app.Context) {
			p.saving = false
			if p.isEditMode {
				ctx.Navigate("/proposals/" + p.proposalID)
			} else {
				// Navigasi ke halaman proposals list atau project detail
				if p.projectID != "" {
					ctx.Navigate("/projects/" + p.projectID)
				} else {
					ctx.Navigate("/proposals")
				}
			}
		})
	}()
}

// onSaveAndApprove menyimpan lalu langsung menyetujui proposal (PO).
func (p *ProposalFormPage) onSaveAndApprove(ctx app.Context, e app.Event) {
	if p.plan == "" {
		p.errMsg = "Plan wajib dipilih."
		return
	}

	p.saving = true
	p.errMsg = ""

	payload := p.buildSavePayload()
	proposalID := p.proposalID

	go func() {
		user := p.authStore.GetUser()
		if user == nil {
			ctx.Dispatch(func(ctx app.Context) { ctx.Navigate("/login") })
			return
		}

		client := api.NewClient("", user.Token)

		// Update proposal
		if err := client.Put(context.Background(), "/api/internal/proposals/"+proposalID, payload, nil); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.saving = false
				p.errMsg = "Gagal menyimpan: " + err.Error()
			})
			return
		}

		// Submit (draft → submitted)
		if err := client.Put(context.Background(), "/api/internal/proposals/"+proposalID+"/submit", nil, nil); err != nil {
			// Mungkin sudah submitted, lanjut approve
		}

		// Approve
		type approveReq struct{}
		if err := client.Put(context.Background(), "/api/internal/proposals/"+proposalID+"/approve", approveReq{}, nil); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.saving = false
				p.errMsg = "Gagal menyetujui: " + err.Error()
			})
			return
		}

		ctx.Dispatch(func(ctx app.Context) {
			p.saving = false
			ctx.Navigate("/proposals/" + proposalID)
		})
	}()
}

// Render menampilkan form proposal.
func (p *ProposalFormPage) Render() app.UI {
	if p.loading {
		return app.Div().
			Style("display", "flex").
			Style("align-items", "center").
			Style("justify-content", "center").
			Style("min-height", "200px").
			Style("color", "#9B8DB5").
			Text("Memuat...")
	}

	role := p.authStore.GetRole()
	isPO := role == "project_owner" || role == "superuser"

	title := "Buat Proposal"
	if p.isEditMode && p.savedVersion != nil {
		title = fmt.Sprintf("Edit Proposal v%d", p.savedVersion.Version)
	}

	return app.Div().
		Style("padding", "32px").
		Style("max-width", "960px").
		Style("margin", "0 auto").
		Body(
			// Header
			app.Div().
				Style("margin-bottom", "24px").
				Body(
					app.H1().
						Style("color", "#E2D9F3").
						Style("font-size", "24px").
						Style("font-weight", "700").
						Style("margin", "0 0 4px").
						Text(title),
					app.If(p.isEditMode && p.savedVersion != nil, func() app.UI {
						return app.Span().
							Style("color", "#9B8DB5").
							Style("font-size", "13px").
							Text(p.savedVersion.CompanyName + " · " + p.savedVersion.ProjectName)
					}),
				),

			// Error message
			app.If(p.errMsg != "", func() app.UI {
				return app.Div().
					Style("background", "rgba(239,68,68,0.1)").
					Style("border", "1px solid rgba(239,68,68,0.4)").
					Style("border-radius", "8px").
					Style("padding", "12px 16px").
					Style("margin-bottom", "16px").
					Style("color", "#EF4444").
					Style("font-size", "14px").
					Text(p.errMsg)
			}),

			// Layout: form (kiri) + diff panel (kanan, hanya edit mode)
			app.Div().
				Style("display", "grid").
				Style("grid-template-columns", func() string {
					if p.isEditMode {
						return "1fr 280px"
					}
					return "1fr"
				}()).
				Style("gap", "24px").
				Style("align-items", "start").
				Body(
					p.renderForm(isPO),
					app.If(p.isEditMode, func() app.UI { return p.renderDiffPanel() }),
				),
		)
}

// renderForm menampilkan form utama.
func (p *ProposalFormPage) renderForm(isPO bool) app.UI {
	role := p.authStore.GetRole()
	isPOOrSuperuser := role == "project_owner" || role == "superuser"

	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "24px").
		Body(
			// Product (hanya create mode)
			app.If(!p.isEditMode, func() app.UI {
				return p.renderFieldGroup("Produk", p.renderProductDropdown())
			}),

			// Plan
			p.renderFieldGroup("Plan", p.renderPlanRadios()),

			// Modules (jika product dipilih)
			app.If(p.productID != "" && p.selectedProduct() != nil, func() app.UI {
				prod := p.selectedProduct()
				if prod == nil || len(prod.AvailableModules) == 0 {
					return app.Text("")
				}
				return p.renderFieldGroup("Modules", p.renderCheckboxGroup(prod.AvailableModules, p.selectedModules, p.onModuleToggle))
			}),

			// Apps (jika product dipilih)
			app.If(p.productID != "" && p.selectedProduct() != nil, func() app.UI {
				prod := p.selectedProduct()
				if prod == nil || len(prod.AvailableApps) == 0 {
					return app.Text("")
				}
				return p.renderFieldGroup("Apps", p.renderAppCheckboxGroup(prod.AvailableApps, p.selectedApps, p.onAppToggle))
			}),

			// Constraints
			p.renderFieldGroup("Constraints", p.renderConstraintInputs()),

			// Kontrak
			p.renderFieldGroup("Kontrak", p.renderContractInputs()),

			// Notes
			p.renderFieldGroup("Catatan untuk PO", p.renderTextarea("notes", p.notes, "Tulis catatan untuk Project Owner...", func(ctx app.Context, e app.Event) {
				p.notes = ctx.JSSrc().Get("value").String()
			})),

			// Owner notes (PO/superuser saja)
			app.If(isPOOrSuperuser, func() app.UI {
				return p.renderFieldGroup("Catatan PO", p.renderTextarea("owner_notes", p.ownerNotes, "Tulis catatan internal PO...", func(ctx app.Context, e app.Event) {
					p.ownerNotes = ctx.JSSrc().Get("value").String()
				}))
			}),

			// Buttons
			app.Div().
				Style("display", "flex").
				Style("gap", "12px").
				Style("margin-top", "24px").
				Body(
					app.Button().
						Style("padding", "10px 24px").
						Style("background", "#4D2975").
						Style("border", "none").
						Style("border-radius", "8px").
						Style("color", "#E2D9F3").
						Style("font-size", "14px").
						Style("font-weight", "600").
						Style("cursor", "pointer").
						Disabled(p.saving).
						OnClick(p.onSave).
						Text(func() string {
							if p.saving {
								return "Menyimpan..."
							}
							return "Simpan Draft"
						}()),

					// Simpan & Setujui (hanya edit mode + PO/superuser)
					app.If(p.isEditMode && isPO, func() app.UI {
						return app.Button().
							Style("padding", "10px 24px").
							Style("background", "rgba(34,197,94,0.15)").
							Style("border", "1px solid #22C55E").
							Style("border-radius", "8px").
							Style("color", "#22C55E").
							Style("font-size", "14px").
							Style("font-weight", "600").
							Style("cursor", "pointer").
							Disabled(p.saving).
							OnClick(p.onSaveAndApprove).
							Text("Simpan & Setujui")
					}),

					// Batal
					app.Button().
						Style("padding", "10px 24px").
						Style("background", "transparent").
						Style("border", "1px solid rgba(155,141,181,0.3)").
						Style("border-radius", "8px").
						Style("color", "#9B8DB5").
						Style("font-size", "14px").
						Style("cursor", "pointer").
						OnClick(func(ctx app.Context, e app.Event) {
							if p.isEditMode {
								ctx.Navigate("/proposals/" + p.proposalID)
							} else if p.projectID != "" {
								ctx.Navigate("/projects/" + p.projectID)
							} else {
								ctx.Navigate("/proposals")
							}
						}).
						Text("Batal"),
				),
		)
}

// renderFieldGroup membungkus field dengan label.
func (p *ProposalFormPage) renderFieldGroup(label string, content app.UI) app.UI {
	return app.Div().
		Style("margin-bottom", "20px").
		Body(
			app.Label().
				Style("display", "block").
				Style("color", "#9B8DB5").
				Style("font-size", "12px").
				Style("font-weight", "600").
				Style("text-transform", "uppercase").
				Style("letter-spacing", "0.06em").
				Style("margin-bottom", "8px").
				Text(label),
			content,
		)
}

// renderProductDropdown menampilkan dropdown untuk memilih produk.
func (p *ProposalFormPage) renderProductDropdown() app.UI {
	opts := []app.UI{
		app.Option().Value("").Text("— Pilih Produk —").Selected(p.productID == ""),
	}
	for _, prod := range p.products {
		opts = append(opts, app.Option().
			Value(prod.ID).
			Text(prod.Name).
			Selected(p.productID == prod.ID))
	}
	return app.Select().
		Style("width", "100%").
		Style("background", "rgba(77,41,117,0.15)").
		Style("border", "1px solid rgba(77,41,117,0.4)").
		Style("border-radius", "8px").
		Style("padding", "10px 12px").
		Style("color", "#E2D9F3").
		Style("font-size", "14px").
		OnChange(p.onProductChange).
		Body(opts...)
}

// renderPlanRadios menampilkan radio buttons untuk pilihan plan.
func (p *ProposalFormPage) renderPlanRadios() app.UI {
	var plans []string
	if prod := p.selectedProduct(); prod != nil {
		plans = prod.AvailablePlans
	}
	if len(plans) == 0 {
		plans = []string{"basic", "professional", "enterprise"}
	}

	radios := make([]app.UI, len(plans))
	for i, plan := range plans {
		isActive := p.plan == plan
		bg := "rgba(77,41,117,0.1)"
		border := "1px solid rgba(77,41,117,0.3)"
		color := "#9B8DB5"
		if isActive {
			bg = "rgba(77,41,117,0.3)"
			border = "1px solid #4D2975"
			color = "#E2D9F3"
		}
		planCopy := plan
		radios[i] = app.Label().
			Style("display", "flex").
			Style("align-items", "center").
			Style("gap", "8px").
			Style("padding", "10px 16px").
			Style("background", bg).
			Style("border", border).
			Style("border-radius", "8px").
			Style("cursor", "pointer").
			Style("color", color).
			Style("font-size", "14px").
			Style("transition", "background 0.15s").
			Body(
				app.Input().
					Type("radio").
					Name("plan").
					Value(planCopy).
					Checked(isActive).
					Style("display", "none").
					OnChange(p.onPlanChange(planCopy)),
				app.Span().Text(strings.ToUpper(planCopy)),
			).
			OnClick(p.onPlanChange(planCopy))
	}

	return app.Div().
		Style("display", "flex").
		Style("flex-wrap", "wrap").
		Style("gap", "10px").
		Body(radios...)
}

// renderCheckboxGroup menampilkan checkboxes untuk modules.
func (p *ProposalFormPage) renderCheckboxGroup(modules []ModuleOption, selected map[string]bool, toggleFn func(string) func(app.Context, app.Event)) app.UI {
	items := make([]app.UI, len(modules))
	for i, m := range modules {
		isChecked := selected[m.Key]
		bg := "rgba(77,41,117,0.1)"
		border := "1px solid rgba(77,41,117,0.3)"
		color := "#9B8DB5"
		if isChecked {
			bg = "rgba(77,41,117,0.3)"
			border = "1px solid #4D2975"
			color = "#E2D9F3"
		}
		key := m.Key
		name := m.Name
		items[i] = app.Label().
			Style("display", "flex").
			Style("align-items", "center").
			Style("gap", "8px").
			Style("padding", "8px 14px").
			Style("background", bg).
			Style("border", border).
			Style("border-radius", "8px").
			Style("cursor", "pointer").
			Style("color", color).
			Style("font-size", "13px").
			OnClick(toggleFn(key)).
			Body(
				app.Input().
					Type("checkbox").
					Checked(isChecked).
					Style("margin", "0").
					OnChange(toggleFn(key)),
				app.Span().Text(name),
			)
	}
	return app.Div().
		Style("display", "flex").
		Style("flex-wrap", "wrap").
		Style("gap", "8px").
		Body(items...)
}

// renderAppCheckboxGroup menampilkan checkboxes untuk apps.
func (p *ProposalFormPage) renderAppCheckboxGroup(apps []AppOption, selected map[string]bool, toggleFn func(string) func(app.Context, app.Event)) app.UI {
	items := make([]app.UI, len(apps))
	for i, a := range apps {
		isChecked := selected[a.Key]
		bg := "rgba(38,184,176,0.1)"
		border := "1px solid rgba(38,184,176,0.3)"
		color := "#9B8DB5"
		if isChecked {
			bg = "rgba(38,184,176,0.25)"
			border = "1px solid #26B8B0"
			color = "#26B8B0"
		}
		key := a.Key
		name := a.Name
		items[i] = app.Label().
			Style("display", "flex").
			Style("align-items", "center").
			Style("gap", "8px").
			Style("padding", "8px 14px").
			Style("background", bg).
			Style("border", border).
			Style("border-radius", "8px").
			Style("cursor", "pointer").
			Style("color", color).
			Style("font-size", "13px").
			OnClick(toggleFn(key)).
			Body(
				app.Input().
					Type("checkbox").
					Checked(isChecked).
					Style("margin", "0").
					OnChange(toggleFn(key)),
				app.Span().Text(name),
			)
	}
	return app.Div().
		Style("display", "flex").
		Style("flex-wrap", "wrap").
		Style("gap", "8px").
		Body(items...)
}

// renderConstraintInputs menampilkan input number untuk constraints.
func (p *ProposalFormPage) renderConstraintInputs() app.UI {
	type constraintField struct {
		label   string
		valPtr  *string
		onInput func(app.Context, app.Event)
	}
	fields := []constraintField{
		{"Max Users", &p.maxUsers, func(ctx app.Context, e app.Event) {
			p.maxUsers = ctx.JSSrc().Get("value").String()
		}},
		{"Max Trans/Bulan", &p.maxTransPerMonth, func(ctx app.Context, e app.Event) {
			p.maxTransPerMonth = ctx.JSSrc().Get("value").String()
		}},
		{"Max Trans/Hari", &p.maxTransPerDay, func(ctx app.Context, e app.Event) {
			p.maxTransPerDay = ctx.JSSrc().Get("value").String()
		}},
		{"Max Items", &p.maxItems, func(ctx app.Context, e app.Event) {
			p.maxItems = ctx.JSSrc().Get("value").String()
		}},
		{"Max Customers", &p.maxCustomers, func(ctx app.Context, e app.Event) {
			p.maxCustomers = ctx.JSSrc().Get("value").String()
		}},
		{"Max Branches", &p.maxBranches, func(ctx app.Context, e app.Event) {
			p.maxBranches = ctx.JSSrc().Get("value").String()
		}},
		{"Max Storage (GB)", &p.maxStorage, func(ctx app.Context, e app.Event) {
			p.maxStorage = ctx.JSSrc().Get("value").String()
		}},
	}

	inputs := make([]app.UI, len(fields))
	for i, f := range fields {
		val := *f.valPtr
		onInput := f.onInput
		inputs[i] = app.Div().
			Style("display", "flex").
			Style("flex-direction", "column").
			Style("gap", "4px").
			Body(
				app.Label().
					Style("color", "#9B8DB5").
					Style("font-size", "12px").
					Text(f.label),
				app.Input().
					Type("number").
					Value(val).
					Min("0").
					Style("background", "rgba(77,41,117,0.15)").
					Style("border", "1px solid rgba(77,41,117,0.4)").
					Style("border-radius", "8px").
					Style("padding", "8px 12px").
					Style("color", "#E2D9F3").
					Style("font-size", "14px").
					Placeholder("—").
					OnInput(onInput),
			)
	}

	return app.Div().
		Style("display", "grid").
		Style("grid-template-columns", "repeat(auto-fill, minmax(160px, 1fr))").
		Style("gap", "14px").
		Body(inputs...)
}

// renderContractInputs menampilkan input untuk contract amount dan expires at.
func (p *ProposalFormPage) renderContractInputs() app.UI {
	return app.Div().
		Style("display", "grid").
		Style("grid-template-columns", "1fr 1fr").
		Style("gap", "14px").
		Body(
			app.Div().
				Style("display", "flex").
				Style("flex-direction", "column").
				Style("gap", "4px").
				Body(
					app.Label().Style("color", "#9B8DB5").Style("font-size", "12px").Text("Contract Amount"),
					app.Input().
						Type("number").
						Value(p.contractAmount).
						Min("0").
						Style("background", "rgba(77,41,117,0.15)").
						Style("border", "1px solid rgba(77,41,117,0.4)").
						Style("border-radius", "8px").
						Style("padding", "8px 12px").
						Style("color", "#E2D9F3").
						Style("font-size", "14px").
						Placeholder("0.00").
						OnInput(func(ctx app.Context, e app.Event) {
							p.contractAmount = ctx.JSSrc().Get("value").String()
						}),
				),
			app.Div().
				Style("display", "flex").
				Style("flex-direction", "column").
				Style("gap", "4px").
				Body(
					app.Label().Style("color", "#9B8DB5").Style("font-size", "12px").Text("Expires At"),
					app.Input().
						Type("date").
						Value(p.expiresAt).
						Style("background", "rgba(77,41,117,0.15)").
						Style("border", "1px solid rgba(77,41,117,0.4)").
						Style("border-radius", "8px").
						Style("padding", "8px 12px").
						Style("color", "#E2D9F3").
						Style("font-size", "14px").
						OnInput(func(ctx app.Context, e app.Event) {
							p.expiresAt = ctx.JSSrc().Get("value").String()
						}),
				),
		)
}

// renderTextarea menampilkan textarea untuk notes.
func (p *ProposalFormPage) renderTextarea(name, value, placeholder string, onInput func(app.Context, app.Event)) app.UI {
	return app.Textarea().
		Style("width", "100%").
		Style("background", "rgba(77,41,117,0.15)").
		Style("border", "1px solid rgba(77,41,117,0.4)").
		Style("border-radius", "8px").
		Style("padding", "10px 12px").
		Style("color", "#E2D9F3").
		Style("font-size", "14px").
		Style("min-height", "90px").
		Style("resize", "vertical").
		Style("box-sizing", "border-box").
		Name(name).
		Placeholder(placeholder).
		Text(value).
		OnInput(onInput)
}

// renderDiffPanel menampilkan panel live diff di sebelah kanan form.
func (p *ProposalFormPage) renderDiffPanel() app.UI {
	changed := p.changedFields()

	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "18px").
		Style("position", "sticky").
		Style("top", "24px").
		Body(
			app.Div().
				Style("color", "#9B8DB5").
				Style("font-size", "12px").
				Style("font-weight", "600").
				Style("text-transform", "uppercase").
				Style("letter-spacing", "0.06em").
				Style("margin-bottom", "12px").
				Text("Perubahan"),

			app.If(len(changed) == 0, func() app.UI {
				return app.Div().
					Style("color", "#9B8DB5").
					Style("font-size", "13px").
					Text("Belum ada perubahan dari versi tersimpan.")
			}),

			app.If(len(changed) > 0, func() app.UI {
				items := make([]app.UI, len(changed))
				for i, field := range changed {
					items[i] = app.Div().
						Style("display", "flex").
						Style("align-items", "center").
						Style("gap", "8px").
						Style("padding", "6px 10px").
						Style("background", "rgba(233,168,0,0.1)").
						Style("border", "1px solid rgba(233,168,0,0.3)").
						Style("border-radius", "6px").
						Style("margin-bottom", "6px").
						Body(
							app.Span().
								Style("width", "8px").
								Style("height", "8px").
								Style("border-radius", "50%").
								Style("background", "#E9A800").
								Style("flex-shrink", "0"),
							app.Span().
								Style("color", "#E9A800").
								Style("font-size", "13px").
								Text(formatFieldName(field)),
						)
				}
				return app.Div().Body(
					app.Div().
						Style("color", "#E2D9F3").
						Style("font-size", "13px").
						Style("margin-bottom", "10px").
						Text(fmt.Sprintf("%d field berubah:", len(changed))),
					app.Div().Body(items...),
				)
			}),
		)
}
