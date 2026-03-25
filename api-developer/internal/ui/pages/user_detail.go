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

// UserDetailPage menampilkan detail user.
type UserDetailPage struct {
	app.Compo
	userID    string
	user      *UserItem
	loading   bool
	errMsg    string
	authStore store.AuthStore
}

// OnNav dipanggil saat halaman ini di-navigasi.
func (p *UserDetailPage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}
	// Extract user ID from URL path: /users/{id}
	path := ctx.Page().URL().Path
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) >= 2 {
		p.userID = parts[1]
	}
	p.fetchUser(ctx)
}

// fetchUser mengambil detail user dari API.
func (p *UserDetailPage) fetchUser(ctx app.Context) {
	if p.userID == "" {
		return
	}
	p.loading = true
	p.errMsg = ""

	token := p.authStore.GetToken()
	userID := p.userID
	ctx.Async(func() {
		client := api.NewClient("", token)
		var detail UserItem
		if err := client.Get(context.Background(), "/api/internal/users/"+userID, &detail); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.loading = false
				p.errMsg = err.Error()
			})
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			p.user = &detail
		})
	})
}

// onBackClick navigasi kembali ke users list.
func (p *UserDetailPage) onBackClick(ctx app.Context, e app.Event) {
	ctx.Navigate("/users")
}

// Render menampilkan halaman detail user.
func (p *UserDetailPage) Render() app.UI {
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

// renderContent merender area konten user detail.
func (p *UserDetailPage) renderContent() app.UI {
	return app.Div().
		Style("padding", "32px").
		Body(
			// Back button
			app.Button().
				Style("background", "none").
				Style("border", "none").
				Style("color", "#26B8B0").
				Style("font-size", "14px").
				Style("cursor", "pointer").
				Style("padding", "0 0 16px").
				OnClick(p.onBackClick).
				Body(
					app.Raw(`<svg style="width:16px;height:16px;vertical-align:middle;margin-right:6px" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7"/></svg>`),
					app.Text("Kembali ke Users"),
				),

			// Loading
			app.If(p.loading,
				func() app.UI {
					return app.Div().
						Style("text-align", "center").
						Style("padding", "60px").
						Style("color", "#9B8DB5").
						Text("Memuat...")
				},
			),

			// Error
			app.If(!p.loading && p.errMsg != "",
				func() app.UI {
					return app.Div().
						Style("background", "rgba(239,68,68,0.1)").
						Style("border", "1px solid #EF4444").
						Style("border-radius", "8px").
						Style("padding", "12px 16px").
						Style("color", "#EF4444").
						Style("font-size", "14px").
						Text(p.errMsg)
				},
			),

			// Content
			app.If(!p.loading && p.user != nil,
				func() app.UI {
					roleBg, roleColor, roleText := roleBadgeStyle(p.user.Role)
					statusBg := "rgba(34,197,94,0.15)"
					statusColor := "#22C55E"
					statusText := "Active"
					if !p.user.IsActive {
						statusBg = "rgba(155,141,181,0.15)"
						statusColor = "#9B8DB5"
						statusText = "Inactive"
					}

					return app.Div().
						Style("background", "#1A1035").
						Style("border", "1px solid rgba(77,41,117,0.3)").
						Style("border-radius", "12px").
						Style("padding", "32px").
						Body(
							app.H2().
								Style("color", "#E2D9F3").
								Style("font-size", "20px").
								Style("font-weight", "700").
								Style("margin", "0 0 24px").
								Text(p.user.Name),

							// Role and Status badges
							app.Div().
								Style("display", "flex").
								Style("gap", "12px").
								Style("margin-bottom", "24px").
								Body(
									app.Span().
										Style("display", "inline-block").
										Style("padding", "3px 10px").
										Style("border-radius", "20px").
										Style("background", roleBg).
										Style("color", roleColor).
										Style("font-size", "12px").
										Style("font-weight", "600").
										Text(roleText),
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

							// Fields
							p.renderField("Email", p.user.Email),
						)
				},
			),
		)
}

// renderField merender satu field dalam detail view.
func (p *UserDetailPage) renderField(label, value string) app.UI {
	if value == "" {
		return app.Div()
	}
	return app.Div().
		Style("margin-bottom", "16px").
		Body(
			app.Div().
				Style("color", "#9B8DB5").
				Style("font-size", "12px").
				Style("text-transform", "uppercase").
				Style("margin-bottom", "6px").
				Text(label),
			app.Div().
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				Text(value),
		)
}
