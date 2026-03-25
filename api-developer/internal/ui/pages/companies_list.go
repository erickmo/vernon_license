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

// CompanyItem adalah representasi company untuk tampilan list.
type CompanyItem struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Address   *string `json:"address"`
	PICName   *string `json:"pic_name"`
	PICEmail  *string `json:"pic_email"`
	PICPhone  *string `json:"pic_phone"`
	Notes     *string `json:"notes"`
	CreatedAt string  `json:"created_at"`
}

// companiesResponse adalah format response dari GET /api/internal/companies.
type companiesResponse struct {
	Data []CompanyItem `json:"data"`
}

// CompaniesListPage menampilkan daftar semua companies dengan form tambah/edit.
type CompaniesListPage struct {
	app.Compo
	authStore store.AuthStore

	companies []CompanyItem
	loading   bool
	errMsg    string

	// Form state
	showForm   bool
	editID     string // kosong = create baru
	formName   string
	formEmail  string
	formPhone  string
	formAddr   string
	formPIC    string
	formPICEm  string
	formPICPh  string
	formNotes  string
	formErr    string
	formSaving bool

	// Delete confirm
	deleteID      string
	deleteConfirm bool
}

// OnNav dipanggil saat halaman ini di-navigasi.
// Redirect ke /login jika belum login, lalu fetch data.
func (p *CompaniesListPage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}
	p.fetchCompanies(ctx)
}

// OnMount dipanggil saat halaman di-mount untuk menangani browser back button.
func (p *CompaniesListPage) OnMount(ctx app.Context) {
	if p.authStore.IsLoggedIn() {
		p.fetchCompanies(ctx)
	}
}

// fetchCompanies mengambil daftar companies dari API.
func (p *CompaniesListPage) fetchCompanies(ctx app.Context) {
	p.loading = true
	p.errMsg = ""

	token := p.authStore.GetToken()
	ctx.Async(func() {
		client := api.NewClient("", token)
		var resp companiesResponse
		err := client.Get(context.Background(), "/api/internal/companies", &resp)
		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			if err != nil {
				p.errMsg = fmt.Sprintf("Gagal memuat companies: %s", err.Error())
				return
			}
			p.companies = resp.Data
		})
	})
}

// onShowCreateForm membuka form tambah company baru.
func (p *CompaniesListPage) onShowCreateForm(ctx app.Context, e app.Event) {
	p.editID = ""
	p.formName = ""
	p.formEmail = ""
	p.formPhone = ""
	p.formAddr = ""
	p.formPIC = ""
	p.formPICEm = ""
	p.formPICPh = ""
	p.formNotes = ""
	p.formErr = ""
	p.showForm = true
}

// onViewClick navigates ke company detail page.
func (p *CompaniesListPage) onViewClick(id string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		ctx.Navigate("/companies/" + id)
	}
}

// onShowEditForm membuka form edit company.
func (p *CompaniesListPage) onShowEditForm(id string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		for _, c := range p.companies {
			if c.ID == id {
				p.editID = c.ID
				p.formName = c.Name
				p.formEmail = safeDeref(c.Email)
				p.formPhone = safeDeref(c.Phone)
				p.formAddr = safeDeref(c.Address)
				p.formPIC = safeDeref(c.PICName)
				p.formPICEm = safeDeref(c.PICEmail)
				p.formPICPh = safeDeref(c.PICPhone)
				p.formNotes = safeDeref(c.Notes)
				p.formErr = ""
				p.showForm = true
				return
			}
		}
	}
}

// onHideForm menutup form.
func (p *CompaniesListPage) onHideForm(ctx app.Context, e app.Event) {
	p.showForm = false
	p.formErr = ""
}

// onFormNameChange menangani input nama.
func (p *CompaniesListPage) onFormNameChange(ctx app.Context, e app.Event) {
	p.formName = ctx.JSSrc().Get("value").String()
}

// onFormEmailChange menangani input email.
func (p *CompaniesListPage) onFormEmailChange(ctx app.Context, e app.Event) {
	p.formEmail = ctx.JSSrc().Get("value").String()
}

// onFormPhoneChange menangani input phone.
func (p *CompaniesListPage) onFormPhoneChange(ctx app.Context, e app.Event) {
	p.formPhone = ctx.JSSrc().Get("value").String()
}

// onFormAddrChange menangani input address.
func (p *CompaniesListPage) onFormAddrChange(ctx app.Context, e app.Event) {
	p.formAddr = ctx.JSSrc().Get("value").String()
}

// onFormPICChange menangani input PIC name.
func (p *CompaniesListPage) onFormPICChange(ctx app.Context, e app.Event) {
	p.formPIC = ctx.JSSrc().Get("value").String()
}

// onFormPICEmChange menangani input PIC email.
func (p *CompaniesListPage) onFormPICEmChange(ctx app.Context, e app.Event) {
	p.formPICEm = ctx.JSSrc().Get("value").String()
}

// onFormPICPhChange menangani input PIC phone.
func (p *CompaniesListPage) onFormPICPhChange(ctx app.Context, e app.Event) {
	p.formPICPh = ctx.JSSrc().Get("value").String()
}

// onFormNotesChange menangani input notes.
func (p *CompaniesListPage) onFormNotesChange(ctx app.Context, e app.Event) {
	p.formNotes = ctx.JSSrc().Get("value").String()
}

// onFormSubmit mengirim form create atau update.
func (p *CompaniesListPage) onFormSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()
	if strings.TrimSpace(p.formName) == "" {
		p.formErr = "Nama company wajib diisi"
		return
	}

	p.formErr = ""
	p.formSaving = true

	token := p.authStore.GetToken()
	editID := p.editID
	body := buildCompanyBody(p)

	ctx.Async(func() {
		client := api.NewClient("", token)
		var err error
		if editID == "" {
			var resp CompanyItem
			err = client.Post(context.Background(), "/api/internal/companies", body, &resp)
		} else {
			var resp CompanyItem
			err = client.Put(context.Background(), "/api/internal/companies/"+editID, body, &resp)
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.formSaving = false
			if err != nil {
				p.formErr = fmt.Sprintf("Gagal menyimpan: %s", err.Error())
				return
			}
			p.showForm = false
			p.fetchCompanies(ctx)
		})
	})
}

// onDeleteClick menampilkan konfirmasi delete.
func (p *CompaniesListPage) onDeleteClick(id string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		p.deleteID = id
		p.deleteConfirm = true
	}
}

// onDeleteCancel membatalkan delete.
func (p *CompaniesListPage) onDeleteCancel(ctx app.Context, e app.Event) {
	p.deleteID = ""
	p.deleteConfirm = false
}

// onDeleteConfirm melakukan delete company.
func (p *CompaniesListPage) onDeleteConfirm(ctx app.Context, e app.Event) {
	deleteID := p.deleteID
	p.deleteConfirm = false
	p.deleteID = ""

	token := p.authStore.GetToken()
	ctx.Async(func() {
		client := api.NewClient("", token)
		err := client.Delete(context.Background(), "/api/internal/companies/"+deleteID)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				p.errMsg = fmt.Sprintf("Gagal menghapus: %s", err.Error())
				return
			}
			p.fetchCompanies(ctx)
		})
	})
}

// Render menampilkan halaman daftar companies.
func (p *CompaniesListPage) Render() app.UI {
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

// renderContent merender area konten utama.
func (p *CompaniesListPage) renderContent() app.UI {
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
								Text("Companies"),
							app.P().
								Style("color", "#9B8DB5").
								Style("font-size", "14px").
								Style("margin", "0").
								Text("Kelola perusahaan klien Vernon"),
						),
					app.Button().
						Style("background", "#4D2975").
						Style("color", "#E2D9F3").
						Style("border", "none").
						Style("border-radius", "8px").
						Style("padding", "10px 20px").
						Style("font-size", "14px").
						Style("font-weight", "600").
						Style("cursor", "pointer").
						OnClick(p.onShowCreateForm).
						Body(
							app.Raw(`<svg style="width:16px;height:16px;margin-right:6px;vertical-align:middle" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/></svg>`),
							app.Text("Tambah Company"),
						),
				),

			// Error global
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

			// Loading state
			app.If(p.loading,
				func() app.UI {
					return app.Div().
						Style("display", "flex").
						Style("align-items", "center").
						Style("justify-content", "center").
						Style("padding", "60px").
						Body(
							app.Div().
								Style("color", "#9B8DB5").
								Style("font-size", "14px").
								Text("Memuat data..."),
						)
				},
			),

			// Empty state
			app.If(!p.loading && len(p.companies) == 0 && p.errMsg == "",
				func() app.UI {
					return app.Div().
						Style("background", "#1A1035").
						Style("border-radius", "12px").
						Style("padding", "60px 32px").
						Style("text-align", "center").
						Body(
							app.Raw(`<svg style="width:48px;height:48px;color:#9B8DB5;margin:0 auto 16px;display:block" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/></svg>`),
							app.P().
								Style("color", "#9B8DB5").
								Style("font-size", "15px").
								Style("margin", "0 0 8px").
								Text("Belum ada company"),
							app.P().
								Style("color", "#6B5E8A").
								Style("font-size", "13px").
								Style("margin", "0").
								Text("Klik tombol Tambah Company untuk mulai"),
						)
				},
			),

			// Table
			app.If(!p.loading && len(p.companies) > 0,
				func() app.UI {
					return p.renderTable()
				},
			),

			// Modal form
			app.If(p.showForm,
				func() app.UI {
					return p.renderForm()
				},
			),

			// Delete confirm modal
			app.If(p.deleteConfirm,
				func() app.UI {
					return p.renderDeleteConfirm()
				},
			),
		)
}

// renderTable merender tabel companies.
func (p *CompaniesListPage) renderTable() app.UI {
	rows := make([]app.UI, 0, len(p.companies))
	for _, c := range p.companies {
		c := c // capture
		rows = append(rows, p.renderRow(c))
	}

	return app.Div().
		Style("background", "#1A1035").
		Style("border-radius", "12px").
		Style("overflow", "hidden").
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
									app.Th().
										Style("text-align", "left").
										Style("padding", "12px 16px").
										Style("color", "#9B8DB5").
										Style("font-size", "12px").
										Style("font-weight", "600").
										Style("text-transform", "uppercase").
										Style("letter-spacing", "0.05em").
										Text("Nama"),
									app.Th().
										Style("text-align", "left").
										Style("padding", "12px 16px").
										Style("color", "#9B8DB5").
										Style("font-size", "12px").
										Style("font-weight", "600").
										Style("text-transform", "uppercase").
										Style("letter-spacing", "0.05em").
										Text("Email"),
									app.Th().
										Style("text-align", "left").
										Style("padding", "12px 16px").
										Style("color", "#9B8DB5").
										Style("font-size", "12px").
										Style("font-weight", "600").
										Style("text-transform", "uppercase").
										Style("letter-spacing", "0.05em").
										Text("PIC"),
									app.Th().
										Style("text-align", "left").
										Style("padding", "12px 16px").
										Style("color", "#9B8DB5").
										Style("font-size", "12px").
										Style("font-weight", "600").
										Style("text-transform", "uppercase").
										Style("letter-spacing", "0.05em").
										Text("Telepon"),
									app.Th().
										Style("text-align", "right").
										Style("padding", "12px 16px").
										Style("color", "#9B8DB5").
										Style("font-size", "12px").
										Style("font-weight", "600").
										Style("text-transform", "uppercase").
										Style("letter-spacing", "0.05em").
										Text("Aksi"),
								),
						),
					app.TBody().Body(rows...),
				),
		)
}

// renderRow merender satu baris tabel company.
func (p *CompaniesListPage) renderRow(c CompanyItem) app.UI {
	return app.Tr().
		Style("border-bottom", "1px solid rgba(77,41,117,0.2)").
		Style("transition", "background 0.15s").
		Style("cursor", "pointer").
		OnClick(p.onViewClick(c.ID)).
		Body(
			app.Td().
				Style("padding", "14px 16px").
				Body(
					app.Div().
						Style("color", "#E2D9F3").
						Style("font-size", "14px").
						Style("font-weight", "500").
						Style("overflow", "hidden").
						Style("text-overflow", "ellipsis").
						Style("white-space", "nowrap").
						Text(c.Name),
					app.If(c.Address != nil && *c.Address != "",
						func() app.UI {
							return app.Div().
								Style("color", "#9B8DB5").
								Style("font-size", "12px").
								Style("margin-top", "2px").
								Style("overflow", "hidden").
								Style("text-overflow", "ellipsis").
								Style("white-space", "nowrap").
								Text(*c.Address)
						},
					),
				),
			app.Td().
				Style("padding", "14px 16px").
				Style("color", "#9B8DB5").
				Style("font-size", "13px").
				Text(safeDeref(c.Email)),
			app.Td().
				Style("padding", "14px 16px").
				Style("color", "#9B8DB5").
				Style("font-size", "13px").
				Text(safeDeref(c.PICName)),
			app.Td().
				Style("padding", "14px 16px").
				Style("color", "#9B8DB5").
				Style("font-size", "13px").
				Text(safeDeref(c.Phone)),
			app.Td().
				Style("padding", "14px 16px").
				Style("text-align", "right").
				Body(
					app.Button().
						Style("background", "rgba(38,184,176,0.15)").
						Style("color", "#26B8B0").
						Style("border", "1px solid rgba(38,184,176,0.3)").
						Style("border-radius", "6px").
						Style("padding", "6px 12px").
						Style("font-size", "12px").
						Style("cursor", "pointer").
						Style("margin-right", "6px").
						OnClick(p.onShowEditForm(c.ID)).
						Text("Edit"),
					app.Button().
						Style("background", "rgba(239,68,68,0.12)").
						Style("color", "#EF4444").
						Style("border", "1px solid rgba(239,68,68,0.3)").
						Style("border-radius", "6px").
						Style("padding", "6px 12px").
						Style("font-size", "12px").
						Style("cursor", "pointer").
						OnClick(p.onDeleteClick(c.ID)).
						Text("Hapus"),
				),
		)
}

// renderForm merender modal form create/edit.
func (p *CompaniesListPage) renderForm() app.UI {
	title := "Tambah Company"
	if p.editID != "" {
		title = "Edit Company"
	}

	return app.Div().
		Style("position", "fixed").
		Style("inset", "0").
		Style("background", "rgba(0,0,0,0.7)").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "center").
		Style("z-index", "100").
		Body(
			app.Div().
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(77,41,117,0.4)").
				Style("border-radius", "16px").
				Style("padding", "32px").
				Style("width", "100%").
				Style("max-width", "540px").
				Style("max-height", "90vh").
				Style("overflow-y", "auto").
				Body(
					// Title
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
								Style("cursor", "pointer").
								Style("font-size", "20px").
								Style("padding", "0").
								OnClick(p.onHideForm).
								Text("×"),
						),

					// Form error
					app.If(p.formErr != "",
						func() app.UI {
							return app.Div().
								Style("background", "rgba(239,68,68,0.15)").
								Style("border", "1px solid rgba(239,68,68,0.4)").
								Style("border-radius", "8px").
								Style("padding", "10px 14px").
								Style("color", "#EF4444").
								Style("font-size", "13px").
								Style("margin-bottom", "16px").
								Text(p.formErr)
						},
					),

					// Form fields
					app.Form().
						OnSubmit(p.onFormSubmit).
						Body(
							renderFormField("Nama *", "text", p.formName, "Nama company", p.onFormNameChange),
							renderFormField("Email", "email", p.formEmail, "email@company.com", p.onFormEmailChange),
							renderFormField("Telepon", "text", p.formPhone, "+62-xxx", p.onFormPhoneChange),
							renderFormField("Alamat", "text", p.formAddr, "Alamat lengkap", p.onFormAddrChange),
							renderFormField("PIC Name", "text", p.formPIC, "Nama PIC", p.onFormPICChange),
							renderFormField("PIC Email", "email", p.formPICEm, "email@pic.com", p.onFormPICEmChange),
							renderFormField("PIC Telepon", "text", p.formPICPh, "+62-xxx", p.onFormPICPhChange),
							renderTextareaField("Notes", p.formNotes, "Catatan tambahan", p.onFormNotesChange),

							// Buttons
							app.Div().
								Style("display", "flex").
								Style("gap", "12px").
								Style("margin-top", "24px").
								Body(
									app.Button().
										Type("submit").
										Style("flex", "1").
										Style("background", "#4D2975").
										Style("color", "#E2D9F3").
										Style("border", "none").
										Style("border-radius", "8px").
										Style("padding", "12px").
										Style("font-size", "14px").
										Style("font-weight", "600").
										Style("cursor", "pointer").
										Text(func() string {
											if p.formSaving {
												return "Menyimpan..."
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
										OnClick(p.onHideForm).
										Text("Batal"),
								),
						),
				),
		)
}

// renderDeleteConfirm merender modal konfirmasi delete.
func (p *CompaniesListPage) renderDeleteConfirm() app.UI {
	return app.Div().
		Style("position", "fixed").
		Style("inset", "0").
		Style("background", "rgba(0,0,0,0.7)").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "center").
		Style("z-index", "100").
		Body(
			app.Div().
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(239,68,68,0.3)").
				Style("border-radius", "16px").
				Style("padding", "32px").
				Style("width", "100%").
				Style("max-width", "400px").
				Body(
					app.H3().
						Style("color", "#E2D9F3").
						Style("font-size", "16px").
						Style("font-weight", "700").
						Style("margin", "0 0 8px").
						Text("Hapus Company"),
					app.P().
						Style("color", "#9B8DB5").
						Style("font-size", "14px").
						Style("margin", "0 0 24px").
						Text("Yakin ingin menghapus company ini? Semua projects yang terkait akan terpengaruh."),
					app.Div().
						Style("display", "flex").
						Style("gap", "12px").
						Body(
							app.Button().
								Style("flex", "1").
								Style("background", "#EF4444").
								Style("color", "#fff").
								Style("border", "none").
								Style("border-radius", "8px").
								Style("padding", "12px").
								Style("font-size", "14px").
								Style("font-weight", "600").
								Style("cursor", "pointer").
								OnClick(p.onDeleteConfirm).
								Text("Hapus"),
							app.Button().
								Style("flex", "1").
								Style("background", "transparent").
								Style("color", "#9B8DB5").
								Style("border", "1px solid rgba(155,141,181,0.3)").
								Style("border-radius", "8px").
								Style("padding", "12px").
								Style("font-size", "14px").
								Style("cursor", "pointer").
								OnClick(p.onDeleteCancel).
								Text("Batal"),
						),
				),
		)
}

// renderFormField merender satu field input dalam form.
func renderFormField(label, inputType, value, placeholder string, onChange func(app.Context, app.Event)) app.UI {
	return app.Div().
		Style("margin-bottom", "16px").
		Body(
			app.Label().
				Style("display", "block").
				Style("color", "#9B8DB5").
				Style("font-size", "12px").
				Style("font-weight", "600").
				Style("margin-bottom", "6px").
				Text(label),
			app.Input().
				Type(inputType).
				Value(value).
				Placeholder(placeholder).
				Style("width", "100%").
				Style("background", "#0F0A1A").
				Style("border", "1px solid rgba(77,41,117,0.4)").
				Style("border-radius", "8px").
				Style("padding", "10px 14px").
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				Style("box-sizing", "border-box").
				OnInput(onChange),
		)
}

// renderTextareaField merender field textarea dalam form.
func renderTextareaField(label, value, placeholder string, onChange func(app.Context, app.Event)) app.UI {
	return app.Div().
		Style("margin-bottom", "16px").
		Body(
			app.Label().
				Style("display", "block").
				Style("color", "#9B8DB5").
				Style("font-size", "12px").
				Style("font-weight", "600").
				Style("margin-bottom", "6px").
				Text(label),
			app.Textarea().
				Placeholder(placeholder).
				Style("width", "100%").
				Style("background", "#0F0A1A").
				Style("border", "1px solid rgba(77,41,117,0.4)").
				Style("border-radius", "8px").
				Style("padding", "10px 14px").
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				Style("box-sizing", "border-box").
				Style("min-height", "80px").
				Style("resize", "vertical").
				OnInput(onChange).
				Text(value),
		)
}

// safeDeref adalah helper untuk dereferensi pointer string safely.
func safeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// strPtr mengkonversi string ke *string. Kosong → nil.
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// buildCompanyBody membangun request body untuk create/update company.
func buildCompanyBody(p *CompaniesListPage) map[string]any {
	return map[string]any{
		"name":      p.formName,
		"email":     strPtr(p.formEmail),
		"phone":     strPtr(p.formPhone),
		"address":   strPtr(p.formAddr),
		"pic_name":  strPtr(p.formPIC),
		"pic_email": strPtr(p.formPICEm),
		"pic_phone": strPtr(p.formPICPh),
		"notes":     strPtr(p.formNotes),
	}
}
