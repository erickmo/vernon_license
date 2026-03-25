//go:build wasm

package pages

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/components"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ProductItem merepresentasikan satu product dalam daftar.
type ProductItem struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	Description      *string         `json:"description"`
	AvailableModules json.RawMessage `json:"available_modules"`
	AvailablePlans   []string        `json:"available_plans"`
	IsActive         bool            `json:"is_active"`
}

// productsListResponse adalah response dari GET /api/internal/products.
type productsListResponse struct {
	Data []ProductItem `json:"data"`
}

// ProductsListPage menampilkan daftar products — hanya untuk superuser.
type ProductsListPage struct {
	app.Compo
	products  []ProductItem
	loading   bool
	showForm  bool
	editID    string
	saving    bool
	deleting  string
	errMsg    string
	formErr   string
	authStore store.AuthStore

	// Form fields
	formName        string
	formSlug        string
	formDescription string
	formIsActive    bool
	formModules     string // JSON array string
	formPlans       []string
}

// OnNav dipanggil saat halaman ini di-navigasi.
// Redirect ke / jika bukan superuser.
func (p *ProductsListPage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}
	if !p.authStore.HasRole("superuser") {
		ctx.Navigate("/")
		return
	}
	p.loadProducts(ctx)
}

// loadProducts mengambil daftar products dari API.
func (p *ProductsListPage) loadProducts(ctx app.Context) {
	p.loading = true
	p.errMsg = ""

	token := p.authStore.GetToken()

	ctx.Async(func() {
		client := api.NewClient("", token)
		var resp productsListResponse
		err := client.Get(ctx, "/api/internal/products", &resp)

		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			if err != nil {
				p.errMsg = "Gagal memuat products."
				return
			}
			p.products = resp.Data
		})
	})
}

// onViewClick navigates ke product detail page.
func (p *ProductsListPage) onViewClick(id string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		ctx.Navigate("/products/" + id)
	}
}

// onOpenCreate membuka form tambah product.
func (p *ProductsListPage) onOpenCreate(ctx app.Context, e app.Event) {
	p.showForm = true
	p.editID = ""
	p.formName = ""
	p.formSlug = ""
	p.formDescription = ""
	p.formIsActive = true
	p.formModules = "[]"
	p.formPlans = []string{}
	p.formErr = ""
}

// onOpenEdit membuka form edit product.
func (p *ProductsListPage) onOpenEdit(id string) func(app.Context, app.Event) {
	return func(ctx app.Context, e app.Event) {
		for _, prod := range p.products {
			if prod.ID == id {
				desc := ""
				if prod.Description != nil {
					desc = *prod.Description
				}
				modStr := "[]"
				if len(prod.AvailableModules) > 0 {
					modStr = string(prod.AvailableModules)
				}
				p.showForm = true
				p.editID = id
				p.formName = prod.Name
				p.formSlug = prod.Slug
				p.formDescription = desc
				p.formIsActive = prod.IsActive
				p.formModules = modStr
				p.formPlans = prod.AvailablePlans
				p.formErr = ""
				return
			}
		}
	}
}

// onCloseForm menutup form.
func (p *ProductsListPage) onCloseForm(ctx app.Context, e app.Event) {
	p.showForm = false
	p.formErr = ""
}

// onSaveProduct menyimpan product (create atau update).
func (p *ProductsListPage) onSaveProduct(ctx app.Context, e app.Event) {
	e.PreventDefault()
	if p.saving {
		return
	}

	if p.formName == "" || p.formSlug == "" {
		p.formErr = "Name dan Slug wajib diisi"
		return
	}

	p.saving = true
	p.formErr = ""

	token := p.authStore.GetToken()
	isEdit := p.editID != ""
	editID := p.editID

	name := p.formName
	slug := p.formSlug
	desc := p.formDescription
	isActive := p.formIsActive
	modulesStr := p.formModules
	plans := p.formPlans

	ctx.Async(func() {
		client := api.NewClient("", token)

		// Parse modules
		var modulesRaw json.RawMessage
		if modulesStr == "" {
			modulesRaw = json.RawMessage("[]")
		} else {
			modulesRaw = json.RawMessage(modulesStr)
		}

		var descPtr *string
		if desc != "" {
			descPtr = &desc
		}

		body := map[string]any{
			"name":              name,
			"slug":              slug,
			"description":       descPtr,
			"is_active":         isActive,
			"available_modules": modulesRaw,
			"available_plans":   plans,
		}

		var updated ProductItem
		var err error
		if isEdit {
			err = client.Put(ctx, "/api/internal/products/"+editID, body, &updated)
		} else {
			err = client.Post(ctx, "/api/internal/products", body, &updated)
		}

		ctx.Dispatch(func(ctx app.Context) {
			p.saving = false
			if err != nil {
				p.formErr = fmt.Sprintf("Gagal menyimpan: %v", err)
				return
			}
			p.showForm = false
			p.editID = ""
			if isEdit {
				// Update in-place agar UI langsung reflect perubahan tanpa re-fetch.
				for i, prod := range p.products {
					if prod.ID == editID {
						p.products[i] = updated
						break
					}
				}
			} else {
				// Tambahkan product baru ke awal list.
				p.products = append([]ProductItem{updated}, p.products...)
			}
		})
	})
}

// onDelete menghapus product.
func (p *ProductsListPage) onDelete(id string) func(app.Context, app.Event) {
	return func(ctx app.Context, e app.Event) {
		p.deleting = id

		token := p.authStore.GetToken()

		ctx.Async(func() {
			client := api.NewClient("", token)
			err := client.Delete(ctx, "/api/internal/products/"+id)

			ctx.Dispatch(func(ctx app.Context) {
				p.deleting = ""
				if err != nil {
					p.errMsg = fmt.Sprintf("Gagal menghapus product: %v", err)
					return
				}
				p.loadProducts(ctx)
			})
		})
	}
}

// onTogglePlan toggle pilihan plan di form.
func (p *ProductsListPage) onTogglePlan(plan string) func(app.Context, app.Event) {
	return func(ctx app.Context, e app.Event) {
		for i, pl := range p.formPlans {
			if pl == plan {
				p.formPlans = append(p.formPlans[:i], p.formPlans[i+1:]...)
				return
			}
		}
		p.formPlans = append(p.formPlans, plan)
	}
}

// hasPlan cek apakah plan dipilih di form.
func (p *ProductsListPage) hasPlan(plan string) bool {
	for _, pl := range p.formPlans {
		if pl == plan {
			return true
		}
	}
	return false
}

// Render menampilkan halaman products.
func (p *ProductsListPage) Render() app.UI {
	if !p.authStore.IsLoggedIn() || !p.authStore.HasRole("superuser") {
		return app.Div()
	}

	return app.Elem("x-shell").
		Body(
			&components.Shell{
				Content: p.renderContent(),
			},
		)
}

// renderContent merender area konten products list.
func (p *ProductsListPage) renderContent() app.UI {
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
					app.Div().
						Body(
							app.H1().
								Style("color", "#E2D9F3").
								Style("font-size", "24px").
								Style("font-weight", "700").
								Style("margin", "0 0 4px").
								Text("Products"),
							app.P().
								Style("color", "#9B8DB5").
								Style("font-size", "14px").
								Style("margin", "0").
								Text("Kelola produk yang dapat dilisensikan"),
						),
					app.Button().
						Style("background", "#4D2975").
						Style("color", "#E2D9F3").
						Style("border", "none").
						Style("border-radius", "8px").
						Style("padding", "10px 18px").
						Style("font-size", "14px").
						Style("font-weight", "600").
						Style("cursor", "pointer").
						OnClick(p.onOpenCreate).
						Text("+ Tambah Product"),
				),

			// Error
			app.If(p.errMsg != "",
				func() app.UI {
					return app.Div().
						Style("background", "rgba(239,68,68,0.1)").
						Style("border", "1px solid #EF4444").
						Style("border-radius", "8px").
						Style("padding", "12px 16px").
						Style("color", "#EF4444").
						Style("font-size", "14px").
						Style("margin-bottom", "20px").
						Text(p.errMsg)
				},
			),

			// Loading
			app.If(p.loading,
				func() app.UI {
					return app.Div().
						Style("color", "#9B8DB5").
						Style("text-align", "center").
						Style("padding", "48px 0").
						Text("Memuat products...")
				},
			),

			// Table
			app.If(!p.loading,
				func() app.UI {
					return p.renderTable()
				},
			),

			// Modal form
			app.If(p.showForm,
				func() app.UI {
					return p.renderModal()
				},
			),
		)
}

// renderTable merender tabel products.
func (p *ProductsListPage) renderTable() app.UI {
	if len(p.products) == 0 {
		return app.Div().
			Style("text-align", "center").
			Style("color", "#9B8DB5").
			Style("padding", "48px 0").
			Text("Belum ada product. Klik 'Tambah Product' untuk membuat.")
	}

	rows := make([]app.UI, 0, len(p.products))
	for _, prod := range p.products {
		prod := prod
		statusBg := "rgba(34,197,94,0.15)"
		statusColor := "#22C55E"
		statusText := "Active"
		if !prod.IsActive {
			statusBg = "rgba(155,141,181,0.15)"
			statusColor = "#9B8DB5"
			statusText = "Inactive"
		}

		rows = append(rows, app.Tr().
			Style("border-bottom", "1px solid rgba(77,41,117,0.2)").
			Style("transition", "background 0.15s").
			Style("cursor", "pointer").
			OnClick(p.onViewClick(prod.ID)).
			Body(
				app.Td().Style("padding", "12px 16px").Style("color", "#E2D9F3").Style("font-size", "14px").Style("font-weight", "500").Text(prod.Name),
				app.Td().Style("padding", "12px 16px").Style("color", "#9B8DB5").Style("font-size", "13px").Style("font-family", "monospace").Text(prod.Slug),
				app.Td().Style("padding", "12px 16px").Body(
					app.Span().
						Style("display", "inline-block").
						Style("padding", "3px 10px").
						Style("border-radius", "20px").
						Style("background", statusBg).
						Style("color", statusColor).
						Style("font-size", "12px").
						Style("font-weight", "600").
						Text(statusText),
				),
				app.Td().Style("padding", "12px 16px").Body(
					app.Span().Style("color", "#9B8DB5").Style("font-size", "13px").Text(strings.Join(prod.AvailablePlans, ", ")),
				),
				app.Td().Style("padding", "12px 16px").Body(
					app.Div().
						Style("display", "flex").
						Style("gap", "8px").
						Body(
							app.Button().
								Style("background", "rgba(77,41,117,0.3)").
								Style("color", "#E2D9F3").
								Style("border", "1px solid rgba(77,41,117,0.5)").
								Style("border-radius", "6px").
								Style("padding", "5px 12px").
								Style("font-size", "12px").
								Style("cursor", "pointer").
								OnClick(p.onOpenEdit(prod.ID)).
								Text("Edit"),
							app.Button().
								Disabled(p.deleting == prod.ID).
								Style("background", "rgba(239,68,68,0.15)").
								Style("color", "#EF4444").
								Style("border", "1px solid rgba(239,68,68,0.3)").
								Style("border-radius", "6px").
								Style("padding", "5px 12px").
								Style("font-size", "12px").
								Style("cursor", func() string {
									if p.deleting == prod.ID {
										return "not-allowed"
									}
									return "pointer"
								}()).
								OnClick(p.onDelete(prod.ID)).
								Text(func() string {
									if p.deleting == prod.ID {
										return "Menghapus..."
									}
									return "Hapus"
								}()),
						),
				),
			))
	}

	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("overflow", "hidden").
		Body(
			app.Table().
				Style("width", "100%").
				Style("border-collapse", "collapse").
				Body(
					app.THead().
						Style("background", "rgba(77,41,117,0.15)").
						Body(
							app.Tr().
								Body(
									app.Th().Style("padding", "12px 16px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Name"),
									app.Th().Style("padding", "12px 16px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Slug"),
									app.Th().Style("padding", "12px 16px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Status"),
									app.Th().Style("padding", "12px 16px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Plans"),
									app.Th().Style("padding", "12px 16px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Actions"),
								),
						),
					app.TBody().Body(rows...),
				),
		)
}

// renderModal merender modal form create/edit product.
func (p *ProductsListPage) renderModal() app.UI {
	title := "Tambah Product"
	if p.editID != "" {
		title = "Edit Product"
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
				Style("border", "1px solid #4D2975").
				Style("border-radius", "16px").
				Style("padding", "32px").
				Style("width", "100%").
				Style("max-width", "560px").
				Style("max-height", "90vh").
				Style("overflow-y", "auto").
				Style("font-family", "'Inter', system-ui, sans-serif").
				Body(
					// Modal header
					app.Div().
						Style("display", "flex").
						Style("align-items", "center").
						Style("justify-content", "space-between").
						Style("margin-bottom", "24px").
						Body(
							app.H2().
								Style("color", "#E2D9F3").
								Style("font-size", "18px").
								Style("font-weight", "700").
								Style("margin", "0").
								Text(title),
							app.Button().
								Style("background", "none").
								Style("border", "none").
								Style("color", "#9B8DB5").
								Style("font-size", "20px").
								Style("cursor", "pointer").
								Style("line-height", "1").
								OnClick(p.onCloseForm).
								Text("×"),
						),

					// Error
					app.If(p.formErr != "",
						func() app.UI {
							return app.Div().
								Style("background", "rgba(239,68,68,0.1)").
								Style("border", "1px solid #EF4444").
								Style("border-radius", "8px").
								Style("padding", "10px 14px").
								Style("color", "#EF4444").
								Style("font-size", "13px").
								Style("margin-bottom", "16px").
								Text(p.formErr)
						},
					),

					// Form
					app.Form().
						OnSubmit(p.onSaveProduct).
						Body(
							// Name
							p.formField("Name", "text", "nama product", p.formName, func(ctx app.Context, e app.Event) {
								p.formName = ctx.JSSrc().Get("value").String()
								// Auto-generate slug dari name jika create baru
								if p.editID == "" {
									p.formSlug = slugify(p.formName)
								}
							}),

							// Slug
							p.formField("Slug", "text", "product-slug", p.formSlug, func(ctx app.Context, e app.Event) {
								p.formSlug = ctx.JSSrc().Get("value").String()
							}),

							// Description
							app.Div().
								Style("margin-bottom", "16px").
								Body(
									app.Label().
										Style("display", "block").
										Style("color", "#9B8DB5").
										Style("font-size", "13px").
										Style("font-weight", "500").
										Style("margin-bottom", "6px").
										Text("Description"),
									app.Textarea().
										Style("width", "100%").
										Style("background", "#0F0A1A").
										Style("border", "1px solid #4D2975").
										Style("border-radius", "8px").
										Style("padding", "10px 14px").
										Style("color", "#E2D9F3").
										Style("font-size", "14px").
										Style("box-sizing", "border-box").
										Style("outline", "none").
										Style("resize", "vertical").
										Style("min-height", "80px").
										Placeholder("Deskripsi product (opsional)").
										Text(p.formDescription).
										OnChange(func(ctx app.Context, e app.Event) {
											p.formDescription = ctx.JSSrc().Get("value").String()
										}),
								),

							// Available Modules (JSON)
							app.Div().
								Style("margin-bottom", "16px").
								Body(
									app.Label().
										Style("display", "block").
										Style("color", "#9B8DB5").
										Style("font-size", "13px").
										Style("font-weight", "500").
										Style("margin-bottom", "6px").
										Text("Available Modules (JSON)"),
									app.Textarea().
										Style("width", "100%").
										Style("background", "#0F0A1A").
										Style("border", "1px solid #4D2975").
										Style("border-radius", "8px").
										Style("padding", "10px 14px").
										Style("color", "#E2D9F3").
										Style("font-size", "13px").
										Style("font-family", "monospace").
										Style("box-sizing", "border-box").
										Style("outline", "none").
										Style("resize", "vertical").
										Style("min-height", "80px").
										Placeholder(`["inventory","accounting","pos"]`).
										Text(p.formModules).
										OnChange(func(ctx app.Context, e app.Event) {
											p.formModules = ctx.JSSrc().Get("value").String()
										}),
								),

							// Available Plans
							app.Div().
								Style("margin-bottom", "16px").
								Body(
									app.Label().
										Style("display", "block").
										Style("color", "#9B8DB5").
										Style("font-size", "13px").
										Style("font-weight", "500").
										Style("margin-bottom", "8px").
										Text("Available Plans"),
									app.Div().
										Style("display", "flex").
										Style("gap", "16px").
										Body(
											p.planCheckbox("saas"),
											p.planCheckbox("dedicated"),
										),
								),

							// Is Active toggle
							app.Div().
								Style("margin-bottom", "24px").
								Style("display", "flex").
								Style("align-items", "center").
								Style("gap", "12px").
								Body(
									app.Input().
										Type("checkbox").
										ID("prod-is-active").
										Checked(p.formIsActive).
										Style("width", "16px").
										Style("height", "16px").
										Style("cursor", "pointer").
										Style("accent-color", "#4D2975").
										OnChange(func(ctx app.Context, e app.Event) {
											p.formIsActive = ctx.JSSrc().Get("checked").Bool()
										}),
									app.Label().
										For("prod-is-active").
										Style("color", "#E2D9F3").
										Style("font-size", "14px").
										Style("cursor", "pointer").
										Text("Product Aktif"),
								),

							// Submit button
							app.Div().
								Style("display", "flex").
								Style("gap", "12px").
								Body(
									app.Button().
										Type("submit").
										Disabled(p.saving).
										Style("flex", "1").
										Style("background", func() string {
											if p.saving {
												return "#3D1F5E"
											}
											return "#4D2975"
										}()).
										Style("color", "#E2D9F3").
										Style("border", "none").
										Style("border-radius", "8px").
										Style("padding", "12px").
										Style("font-size", "14px").
										Style("font-weight", "600").
										Style("cursor", func() string {
											if p.saving {
												return "not-allowed"
											}
											return "pointer"
										}()).
										Text(func() string {
											if p.saving {
												return "Menyimpan..."
											}
											if p.editID != "" {
												return "Update"
											}
											return "Simpan"
										}()),
									app.Button().
										Type("button").
										Style("flex", "1").
										Style("background", "transparent").
										Style("color", "#9B8DB5").
										Style("border", "1px solid rgba(155,141,181,0.3)").
										Style("border-radius", "8px").
										Style("padding", "12px").
										Style("font-size", "14px").
										Style("cursor", "pointer").
										OnClick(p.onCloseForm).
										Text("Batal"),
								),
						),
				),
		)
}

// formField adalah helper untuk merender satu input field.
func (p *ProductsListPage) formField(label, inputType, placeholder, value string, onChange func(app.Context, app.Event)) app.UI {
	return app.Div().
		Style("margin-bottom", "16px").
		Body(
			app.Label().
				Style("display", "block").
				Style("color", "#9B8DB5").
				Style("font-size", "13px").
				Style("font-weight", "500").
				Style("margin-bottom", "6px").
				Text(label),
			app.Input().
				Type(inputType).
				Placeholder(placeholder).
				Value(value).
				Style("width", "100%").
				Style("background", "#0F0A1A").
				Style("border", "1px solid #4D2975").
				Style("border-radius", "8px").
				Style("padding", "10px 14px").
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				Style("box-sizing", "border-box").
				Style("outline", "none").
				OnChange(onChange),
		)
}

// planCheckbox merender satu plan checkbox.
func (p *ProductsListPage) planCheckbox(plan string) app.UI {
	checked := p.hasPlan(plan)
	return app.Div().
		Style("display", "flex").
		Style("align-items", "center").
		Style("gap", "8px").
		Body(
			app.Input().
				Type("checkbox").
				ID("plan-"+plan).
				Checked(checked).
				Style("width", "16px").
				Style("height", "16px").
				Style("cursor", "pointer").
				Style("accent-color", "#4D2975").
				OnChange(p.onTogglePlan(plan)),
			app.Label().
				For("plan-"+plan).
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				Style("cursor", "pointer").
				Style("text-transform", "capitalize").
				Text(plan),
		)
}

// slugify mengubah string menjadi slug sederhana (lowercase, spaces → dashes).
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// Hapus karakter non-alphanumeric kecuali dash
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' {
			result = append(result, c)
		}
	}
	return string(result)
}
