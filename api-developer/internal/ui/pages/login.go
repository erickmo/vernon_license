//go:build wasm

// Package pages berisi semua halaman UI untuk Vernon App.
package pages

import (
	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// loginResponse adalah response dari POST /api/internal/auth/login.
type loginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Role string `json:"role"`
	} `json:"user"`
}

// LoginPage adalah halaman login Vernon App.
// Menampilkan form email + password dengan tema gelap.
type LoginPage struct {
	app.Compo
	email     string
	password  string
	loading   bool
	errMsg    string
	authStore store.AuthStore
}

// OnNav dipanggil saat halaman ini di-navigasi.
// Redirect ke / jika sudah login.
func (p *LoginPage) OnNav(ctx app.Context) {
	if p.authStore.IsLoggedIn() {
		ctx.Navigate("/")
	}
}

// Render menampilkan form login dengan tema gelap.
// Design: full-screen #0F0A1A, centered card dengan gradient border #4D2975.
func (p *LoginPage) Render() app.UI {
	return app.Div().
		Style("min-height", "100vh").
		Style("background", "#0F0A1A").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "center").
		Style("font-family", "'Inter', system-ui, sans-serif").
		Body(
			app.Div().
				Style("width", "100%").
				Style("max-width", "420px").
				Style("padding", "0 16px").
				Body(
					// Card
					app.Div().
						Style("background", "#1A1035").
						Style("border", "1px solid #4D2975").
						Style("border-radius", "16px").
						Style("padding", "40px 32px").
						Style("box-shadow", "0 0 40px rgba(77,41,117,0.3)").
						Body(
							// Logo + Title
							app.Div().
								Style("text-align", "center").
								Style("margin-bottom", "32px").
								Body(
									app.Div().
										Style("width", "48px").
										Style("height", "48px").
										Style("background", "linear-gradient(135deg, #4D2975, #26B8B0)").
										Style("border-radius", "12px").
										Style("margin", "0 auto 16px").
										Style("display", "flex").
										Style("align-items", "center").
										Style("justify-content", "center").
										Body(
											app.Text("V"),
										),
									app.H1().
										Style("color", "#E2D9F3").
										Style("font-size", "24px").
										Style("font-weight", "700").
										Style("margin", "0 0 4px").
										Text("Vernon License"),
									app.P().
										Style("color", "#9B8DB5").
										Style("font-size", "14px").
										Style("margin", "0").
										Text("License Management System"),
								),

							// Error message
							app.If(p.errMsg != "",
								func() app.UI {
									return app.Div().
										Style("background", "rgba(239,68,68,0.1)").
										Style("border", "1px solid #EF4444").
										Style("border-radius", "8px").
										Style("padding", "12px 16px").
										Style("margin-bottom", "20px").
										Style("color", "#EF4444").
										Style("font-size", "14px").
										Text(p.errMsg)
								},
							),

							// Form
							app.Form().
								OnSubmit(p.onSubmit).
								Body(
									// Email field
									app.Div().
										Style("margin-bottom", "16px").
										Body(
											app.Label().
												Style("display", "block").
												Style("color", "#9B8DB5").
												Style("font-size", "13px").
												Style("font-weight", "500").
												Style("margin-bottom", "6px").
												For("email").
												Text("Email"),
											app.Input().
												ID("email").
												Type("email").
												Placeholder("email@example.com").
												Value(p.email).
												Required(true).
												Style("width", "100%").
												Style("background", "#0F0A1A").
												Style("border", "1px solid #4D2975").
												Style("border-radius", "8px").
												Style("padding", "10px 14px").
												Style("color", "#E2D9F3").
												Style("font-size", "15px").
												Style("box-sizing", "border-box").
												Style("outline", "none").
												OnChange(p.onEmailChange),
										),

									// Password field
									app.Div().
										Style("margin-bottom", "24px").
										Body(
											app.Label().
												Style("display", "block").
												Style("color", "#9B8DB5").
												Style("font-size", "13px").
												Style("font-weight", "500").
												Style("margin-bottom", "6px").
												For("password").
												Text("Password"),
											app.Input().
												ID("password").
												Type("password").
												Placeholder("••••••••").
												Value(p.password).
												Required(true).
												Style("width", "100%").
												Style("background", "#0F0A1A").
												Style("border", "1px solid #4D2975").
												Style("border-radius", "8px").
												Style("padding", "10px 14px").
												Style("color", "#E2D9F3").
												Style("font-size", "15px").
												Style("box-sizing", "border-box").
												Style("outline", "none").
												OnChange(p.onPasswordChange),
										),

									// Submit button
									app.Button().
										Type("submit").
										Disabled(p.loading).
										Style("width", "100%").
										Style("background", func() string {
											if p.loading {
												return "#3D1F5E"
											}
											return "#4D2975"
										}()).
										Style("color", "#E2D9F3").
										Style("border", "none").
										Style("border-radius", "8px").
										Style("padding", "12px").
										Style("font-size", "15px").
										Style("font-weight", "600").
										Style("cursor", func() string {
											if p.loading {
												return "not-allowed"
											}
											return "pointer"
										}()).
										Style("transition", "background 0.2s").
										Body(
											app.If(p.loading,
												func() app.UI {
													return app.Text("Memproses...")
												},
											).Else(
												func() app.UI {
													return app.Text("Masuk")
												},
											),
										),
								),
						),
				),
		)
}

// onEmailChange menyimpan nilai email dari input.
func (p *LoginPage) onEmailChange(ctx app.Context, e app.Event) {
	p.email = ctx.JSSrc().Get("value").String()
}

// onPasswordChange menyimpan nilai password dari input.
func (p *LoginPage) onPasswordChange(ctx app.Context, e app.Event) {
	p.password = ctx.JSSrc().Get("value").String()
}

// onSubmit handle form submit.
// POST ke /api/internal/auth/login.
// Simpan ke authStore jika berhasil, redirect ke /.
func (p *LoginPage) onSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	if p.loading {
		return
	}

	p.loading = true
	p.errMsg = ""

	email := p.email
	password := p.password

	ctx.Async(func() {
		client := api.NewClient("", "")

		var resp loginResponse
		err := client.Post(ctx, "/api/internal/auth/login", map[string]string{
			"email":    email,
			"password": password,
		}, &resp)

		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false

			if err != nil {
				p.errMsg = "Email atau password salah."
				return
			}

			authUser := store.AuthUser{
				ID:    resp.User.ID,
				Name:  resp.User.Name,
				Role:  resp.User.Role,
				Token: resp.Token,
			}
			if saveErr := p.authStore.Save(authUser); saveErr != nil {
				p.errMsg = "Gagal menyimpan sesi. Coba lagi."
				return
			}

			ctx.Navigate("/")
		})
	})
}
