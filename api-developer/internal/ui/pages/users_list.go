//go:build wasm

package pages

import (
	"fmt"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/components"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// UserItem merepresentasikan satu user dalam daftar.
type UserItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

// usersListResponse adalah response dari GET /api/internal/users.
type usersListResponse struct {
	Data []UserItem `json:"data"`
}

// UsersListPage menampilkan daftar users — hanya untuk superuser.
type UsersListPage struct {
	app.Compo
	users     []UserItem
	loading   bool
	showForm  bool
	saving    bool
	toggling  string
	errMsg    string
	formErr   string
	authStore store.AuthStore

	// Form fields
	formName     string
	formEmail    string
	formPassword string
	formRole     string
}

// OnNav dipanggil saat halaman ini di-navigasi.
// Redirect ke / jika bukan superuser.
func (p *UsersListPage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}
	if !p.authStore.HasRole("superuser") {
		ctx.Navigate("/")
		return
	}
	p.loadUsers(ctx)
}

// loadUsers mengambil daftar users dari API.
func (p *UsersListPage) loadUsers(ctx app.Context) {
	p.loading = true
	p.errMsg = ""

	token := p.authStore.GetToken()

	ctx.Async(func() {
		client := api.NewClient("", token)
		var resp usersListResponse
		err := client.Get(ctx, "/api/internal/users", &resp)

		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			if err != nil {
				p.errMsg = "Gagal memuat users."
				return
			}
			p.users = resp.Data
		})
	})
}

// onOpenCreate membuka form tambah user.
func (p *UsersListPage) onOpenCreate(ctx app.Context, e app.Event) {
	p.showForm = true
	p.formName = ""
	p.formEmail = ""
	p.formPassword = ""
	p.formRole = "project_owner"
	p.formErr = ""
}

// onCloseForm menutup form.
func (p *UsersListPage) onCloseForm(ctx app.Context, e app.Event) {
	p.showForm = false
	p.formErr = ""
}

// onSaveUser menyimpan user baru.
func (p *UsersListPage) onSaveUser(ctx app.Context, e app.Event) {
	e.PreventDefault()
	if p.saving {
		return
	}

	if p.formName == "" || p.formEmail == "" || p.formPassword == "" || p.formRole == "" {
		p.formErr = "Semua field wajib diisi"
		return
	}

	p.saving = true
	p.formErr = ""

	token := p.authStore.GetToken()
	name := p.formName
	email := p.formEmail
	password := p.formPassword
	role := p.formRole

	ctx.Async(func() {
		client := api.NewClient("", token)
		body := map[string]string{
			"name":     name,
			"email":    email,
			"password": password,
			"role":     role,
		}
		err := client.Post(ctx, "/api/internal/users", body, nil)

		ctx.Dispatch(func(ctx app.Context) {
			p.saving = false
			if err != nil {
				p.formErr = fmt.Sprintf("Gagal membuat user: %v", err)
				return
			}
			p.showForm = false
			p.loadUsers(ctx)
		})
	})
}

// onToggleActive mengaktifkan atau menonaktifkan user.
func (p *UsersListPage) onToggleActive(id string, currentActive bool) func(app.Context, app.Event) {
	return func(ctx app.Context, e app.Event) {
		p.toggling = id

		token := p.authStore.GetToken()
		newActive := !currentActive

		ctx.Async(func() {
			client := api.NewClient("", token)
			body := map[string]bool{"is_active": newActive}
			err := client.Put(ctx, "/api/internal/users/"+id+"/active", body, nil)

			ctx.Dispatch(func(ctx app.Context) {
				p.toggling = ""
				if err != nil {
					p.errMsg = fmt.Sprintf("Gagal mengubah status user: %v", err)
					return
				}
				p.loadUsers(ctx)
			})
		})
	}
}

// Render menampilkan halaman users.
func (p *UsersListPage) Render() app.UI {
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

// renderContent merender area konten users list.
func (p *UsersListPage) renderContent() app.UI {
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
								Text("Users"),
							app.P().
								Style("color", "#9B8DB5").
								Style("font-size", "14px").
								Style("margin", "0").
								Text("Kelola akun pengguna Vernon App"),
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
						Text("+ Tambah User"),
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
						Text("Memuat users...")
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

// renderTable merender tabel users.
func (p *UsersListPage) renderTable() app.UI {
	if len(p.users) == 0 {
		return app.Div().
			Style("text-align", "center").
			Style("color", "#9B8DB5").
			Style("padding", "48px 0").
			Text("Belum ada user selain superuser.")
	}

	rows := make([]app.UI, 0, len(p.users))
	for _, user := range p.users {
		user := user

		// Role badge
		roleBg, roleColor, roleText := roleBadgeStyle(user.Role)

		// Status badge
		statusBg := "rgba(34,197,94,0.15)"
		statusColor := "#22C55E"
		statusText := "Active"
		if !user.IsActive {
			statusBg = "rgba(155,141,181,0.15)"
			statusColor = "#9B8DB5"
			statusText = "Inactive"
		}

		// Toggle button — superuser tidak bisa di-deactivate.
		var toggleUI app.UI
		if user.Role == "superuser" {
			toggleUI = app.Span().
				Style("color", "#9B8DB5").
				Style("font-size", "12px").
				Style("font-style", "italic").
				Text("—")
		} else {
			toggleLabel := "Deactivate"
			toggleBg := "rgba(239,68,68,0.15)"
			toggleColor := "#EF4444"
			toggleBorder := "rgba(239,68,68,0.3)"
			if !user.IsActive {
				toggleLabel = "Activate"
				toggleBg = "rgba(34,197,94,0.15)"
				toggleColor = "#22C55E"
				toggleBorder = "rgba(34,197,94,0.3)"
			}
			if p.toggling == user.ID {
				toggleLabel = "..."
			}
			toggleUI = app.Button().
				Disabled(p.toggling == user.ID).
				Style("background", toggleBg).
				Style("color", toggleColor).
				Style("border", "1px solid "+toggleBorder).
				Style("border-radius", "6px").
				Style("padding", "5px 12px").
				Style("font-size", "12px").
				Style("cursor", func() string {
					if p.toggling == user.ID {
						return "not-allowed"
					}
					return "pointer"
				}()).
				OnClick(p.onToggleActive(user.ID, user.IsActive)).
				Text(toggleLabel)
		}

		rows = append(rows, app.Tr().
			Style("border-bottom", "1px solid rgba(77,41,117,0.2)").
			Body(
				app.Td().Style("padding", "12px 16px").Body(
					app.Div().Style("color", "#E2D9F3").Style("font-size", "14px").Style("font-weight", "500").Text(user.Name),
					app.Div().Style("color", "#9B8DB5").Style("font-size", "12px").Style("margin-top", "2px").Text(user.Email),
				),
				app.Td().Style("padding", "12px 16px").Body(
					app.Span().
						Style("display", "inline-block").
						Style("padding", "3px 10px").
						Style("border-radius", "20px").
						Style("background", roleBg).
						Style("color", roleColor).
						Style("font-size", "12px").
						Style("font-weight", "600").
						Text(roleText),
				),
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
				app.Td().Style("padding", "12px 16px").Body(toggleUI),
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
									app.Th().Style("padding", "12px 16px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("User"),
									app.Th().Style("padding", "12px 16px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Role"),
									app.Th().Style("padding", "12px 16px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Status"),
									app.Th().Style("padding", "12px 16px").Style("text-align", "left").Style("color", "#9B8DB5").Style("font-size", "12px").Style("font-weight", "500").Style("text-transform", "uppercase").Text("Actions"),
								),
						),
					app.TBody().Body(rows...),
				),
		)
}

// renderModal merender modal form create user.
func (p *UsersListPage) renderModal() app.UI {
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
				Style("max-width", "480px").
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
								Text("Tambah User"),
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
						OnSubmit(p.onSaveUser).
						Body(
							// Name
							p.userFormField("Nama", "text", "nama lengkap", p.formName, func(ctx app.Context, e app.Event) {
								p.formName = ctx.JSSrc().Get("value").String()
							}),

							// Email
							p.userFormField("Email", "email", "email@example.com", p.formEmail, func(ctx app.Context, e app.Event) {
								p.formEmail = ctx.JSSrc().Get("value").String()
							}),

							// Password
							p.userFormField("Password", "password", "••••••••", p.formPassword, func(ctx app.Context, e app.Event) {
								p.formPassword = ctx.JSSrc().Get("value").String()
							}),

							// Role dropdown
							app.Div().
								Style("margin-bottom", "24px").
								Body(
									app.Label().
										Style("display", "block").
										Style("color", "#9B8DB5").
										Style("font-size", "13px").
										Style("font-weight", "500").
										Style("margin-bottom", "6px").
										Text("Role"),
									app.Select().
										Style("width", "100%").
										Style("background", "#0F0A1A").
										Style("border", "1px solid #4D2975").
										Style("border-radius", "8px").
										Style("padding", "10px 14px").
										Style("color", "#E2D9F3").
										Style("font-size", "14px").
										Style("box-sizing", "border-box").
										Style("outline", "none").
										OnChange(func(ctx app.Context, e app.Event) {
											p.formRole = ctx.JSSrc().Get("value").String()
										}).
										Body(
											app.Option().Value("project_owner").
												Selected(p.formRole == "project_owner").
												Text("Project Owner"),
											app.Option().Value("sales").
												Selected(p.formRole == "sales").
												Text("Sales"),
										),
								),

							// Submit
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
											return "Buat User"
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

// userFormField adalah helper untuk merender satu input field di user form.
func (p *UsersListPage) userFormField(label, inputType, placeholder, value string, onChange func(app.Context, app.Event)) app.UI {
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

// roleBadgeStyle mengembalikan background, color, dan label untuk role badge.
func roleBadgeStyle(role string) (bg, color, label string) {
	switch role {
	case "superuser":
		return "rgba(77,41,117,0.3)", "#4D2975", "Superuser"
	case "project_owner":
		return "rgba(38,184,176,0.15)", "#26B8B0", "Project Owner"
	case "sales":
		return "rgba(233,168,0,0.15)", "#E9A800", "Sales"
	default:
		return "rgba(155,141,181,0.15)", "#9B8DB5", role
	}
}
