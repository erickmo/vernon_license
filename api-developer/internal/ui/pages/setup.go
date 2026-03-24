//go:build wasm

package pages

import (
	"encoding/json"
	"fmt"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// SetupPage adalah halaman first-run setup untuk membuat superuser pertama.
// Ditampilkan secara otomatis jika belum ada user di database.
type SetupPage struct {
	app.Compo
	name      string
	email     string
	password  string
	password2 string
	loading   bool
	errMsg    string
	authStore store.AuthStore
}

// OnNav memeriksa apakah setup masih diperlukan. Jika sudah setup, redirect ke login.
func (p *SetupPage) OnNav(ctx app.Context) {
	ctx.Async(func() {
		client := api.NewClient("", "")
		var resp struct {
			IsSetup bool `json:"is_setup"`
		}
		if err := client.Get(ctx, "/api/internal/setup/status", &resp); err == nil && resp.IsSetup {
			ctx.Dispatch(func(ctx app.Context) {
				ctx.Navigate("/login")
			})
		}
	})
}

// Render menampilkan form setup superuser dengan tema gelap.
func (p *SetupPage) Render() app.UI {
	return app.Div().
		Style("min-height", "100vh").
		Style("background", "#0F0A1A").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "center").
		Style("padding", "24px").
		Body(
			app.Div().
				Style("width", "100%").
				Style("max-width", "440px").
				Body(
					p.renderHeader(),
					p.renderCard(),
				),
		)
}

func (p *SetupPage) renderHeader() app.UI {
	return app.Div().
		Style("text-align", "center").
		Style("margin-bottom", "32px").
		Body(
			app.Div().
				Style("width", "64px").
				Style("height", "64px").
				Style("background", "linear-gradient(135deg, #4D2975, #26B8B0)").
				Style("border-radius", "16px").
				Style("display", "flex").
				Style("align-items", "center").
				Style("justify-content", "center").
				Style("font-size", "28px").
				Style("margin", "0 auto 16px").
				Text("🔑"),
			app.H1().
				Style("font-size", "24px").
				Style("font-weight", "700").
				Style("color", "#E2D9F3").
				Style("letter-spacing", "-0.5px").
				Text("Selamat Datang di Vernon"),
			app.P().
				Style("font-size", "14px").
				Style("color", "#9B8DB5").
				Style("margin-top", "8px").
				Text("Buat akun superuser untuk memulai"),
		)
}

func (p *SetupPage) renderCard() app.UI {
	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.4)").
		Style("border-radius", "16px").
		Style("padding", "32px").
		Body(
			// Error message
			app.If(p.errMsg != "", func() app.UI {
				return app.Div().
					Style("background", "rgba(239,68,68,0.1)").
					Style("border", "1px solid rgba(239,68,68,0.3)").
					Style("border-radius", "8px").
					Style("padding", "12px 16px").
					Style("margin-bottom", "20px").
					Style("font-size", "13px").
					Style("color", "#EF4444").
					Text(p.errMsg)
			}),

			p.renderInput("Nama Lengkap", "text", "Nama superuser", p.name, "name"),
			p.renderInput("Email", "email", "admin@company.com", p.email, "email"),
			p.renderInput("Password", "password", "Minimal 8 karakter", p.password, "password"),
			p.renderInput("Konfirmasi Password", "password", "Ulangi password", p.password2, "password2"),

			app.Button().
				Style("width", "100%").
				Style("padding", "14px").
				Style("background", func() string {
					if p.loading {
						return "#3D1F5E"
					}
					return "linear-gradient(135deg, #4D2975, #3D1F5E)"
				}()).
				Style("border", "none").
				Style("border-radius", "10px").
				Style("color", "#E2D9F3").
				Style("font-size", "15px").
				Style("font-weight", "600").
				Style("cursor", func() string {
					if p.loading {
						return "not-allowed"
					}
					return "pointer"
				}()).
				Style("margin-top", "8px").
				Disabled(p.loading).
				OnClick(p.onSubmit).
				Text(func() string {
					if p.loading {
						return "Menyiapkan sistem..."
					}
					return "Mulai Setup"
				}()),
		)
}

func (p *SetupPage) renderInput(label, inputType, placeholder, value, field string) app.UI {
	return app.Div().
		Style("margin-bottom", "16px").
		Body(
			app.Label().
				Style("display", "block").
				Style("font-size", "13px").
				Style("font-weight", "500").
				Style("color", "#9B8DB5").
				Style("margin-bottom", "6px").
				Text(label),
			app.Input().
				Style("width", "100%").
				Style("padding", "11px 14px").
				Style("background", "#0F0A1A").
				Style("border", "1px solid rgba(77,41,117,0.5)").
				Style("border-radius", "8px").
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				Style("outline", "none").
				Type(inputType).
				Placeholder(placeholder).
				Value(value).
				OnChange(p.onFieldChange(field)),
		)
}

func (p *SetupPage) onFieldChange(field string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		val := ctx.JSSrc().Get("value").String()
		switch field {
		case "name":
			p.name = val
		case "email":
			p.email = val
		case "password":
			p.password = val
		case "password2":
			p.password2 = val
		}
	}
}

func (p *SetupPage) onSubmit(ctx app.Context, e app.Event) {
	p.errMsg = ""

	// Validasi
	if p.name == "" || p.email == "" || p.password == "" {
		p.errMsg = "Semua field wajib diisi"
		return
	}
	if len(p.password) < 8 {
		p.errMsg = "Password minimal 8 karakter"
		return
	}
	if p.password != p.password2 {
		p.errMsg = "Password tidak cocok"
		return
	}

	p.loading = true

	ctx.Async(func() {
		client := api.NewClient("", "")
		body := map[string]string{
			"name":     p.name,
			"email":    p.email,
			"password": p.password,
		}

		var resp struct {
			Token string `json:"token"`
			User  struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Role string `json:"role"`
			} `json:"user"`
		}

		bodyBytes, _ := json.Marshal(body)
		fmt.Println("DEBUG setup install", string(bodyBytes))

		err := client.Post(ctx, "/api/internal/setup/install", body, &resp)
		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			if err != nil {
				p.errMsg = "Setup gagal: " + err.Error()
				return
			}
			// Simpan auth info dan redirect ke dashboard
			p.authStore.Save(store.AuthUser{
				ID:    resp.User.ID,
				Name:  resp.User.Name,
				Role:  resp.User.Role,
				Token: resp.Token,
			})
			ctx.Navigate("/")
		})
	})
}
