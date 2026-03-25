//go:build wasm

// Package pages menyediakan semua halaman aplikasi Vernon License.
package pages

import (
	"fmt"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ActivityLogPage menampilkan log aktivitas sistem.
type ActivityLogPage struct {
	app.Compo
	authStore  store.AuthStore
	activities []activityItem
	loading    bool
	errMsg     string
}

// OnNav dipanggil saat halaman ini di-navigasi.
func (p *ActivityLogPage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}
	p.loadActivities(ctx)
}

// loadActivities mengambil activity log dari dashboard stats API.
func (p *ActivityLogPage) loadActivities(ctx app.Context) {
	p.loading = true
	p.errMsg = ""
	token := p.authStore.GetToken()

	ctx.Async(func() {
		client := api.NewClient("", token)
		type dashboardStats struct {
			RecentActivity []activityItem `json:"recent_activity"`
		}
		var stats dashboardStats
		err := client.Get(ctx, "/api/internal/dashboard/stats", &stats)

		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			if err != nil {
				p.errMsg = fmt.Sprintf("Gagal memuat activity log: %v", err)
				return
			}
			p.activities = stats.RecentActivity
		})
	})
}

// Render merender halaman activity log.
func (p *ActivityLogPage) Render() app.UI {
	return app.Div().
		Style("background", "#0F0620").
		Style("min-height", "100vh").
		Style("color", "#E2D9F3").
		Style("font-family", "-apple-system, BlinkMacSystemFont, Segoe UI, Roboto").
		Body(
			// Header
			app.Div().
				Style("background", "linear-gradient(135deg, #1A1035 0%, #2D1B4E 100%)").
				Style("padding", "20px 24px").
				Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
				Body(
					app.H1().
						Style("margin", "0").
						Style("font-size", "24px").
						Style("font-weight", "700").
						Text("Activity Log"),
				),

			// Content
			app.Div().
				Style("padding", "24px").
				Body(
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
								Style("margin-bottom", "20px").
								Text(p.errMsg)
						},
					),

					// Activity list
					app.If(len(p.activities) > 0,
						func() app.UI {
							return app.Div().
								Style("background", "#1A1035").
								Style("border", "1px solid rgba(77,41,117,0.3)").
								Style("border-radius", "12px").
								Style("padding", "20px").
								Body(
									p.renderActivityItems(),
								)
						},
					),
					app.If(len(p.activities) == 0,
						func() app.UI {
							return app.Div().
								Style("text-align", "center").
								Style("padding", "40px 20px").
								Style("color", "#9B8DB5").
								Body(
									app.Div().
										Style("font-size", "16px").
										Text("Tidak ada activity log."),
								)
						},
					),
				),
		)
}

// renderActivityItems merender daftar aktivitas.
func (p *ActivityLogPage) renderActivityItems() app.UI {
	items := make([]app.UI, 0, len(p.activities))
	for _, act := range p.activities {
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

	return app.Div().Body(items...)
}
