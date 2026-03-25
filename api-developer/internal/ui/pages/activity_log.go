//go:build wasm

// Package pages menyediakan semua halaman aplikasi Vernon License.
package pages

import (
	"context"
	"fmt"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/components"
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
		type dashboardResp struct {
			RecentActivity []activityItem `json:"recent_activity"`
		}
		var stats dashboardResp
		err := client.Get(context.Background(), "/api/internal/dashboard", &stats)

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

// Render merender halaman activity log menggunakan Shell standar.
func (p *ActivityLogPage) Render() app.UI {
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

// renderContent merender konten utama halaman activity log.
func (p *ActivityLogPage) renderContent() app.UI {
	return app.Div().
		Style("padding", "32px").
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
						Text("Activity Log"),
					app.P().
						Style("color", "#9B8DB5").
						Style("font-size", "14px").
						Style("margin", "0").
						Text("10 aktivitas terakhir dalam sistem"),
				),

			// Loading state
			app.If(p.loading, func() app.UI {
				return app.Div().
					Style("text-align", "center").
					Style("padding", "60px").
					Style("color", "#9B8DB5").
					Text("Memuat activity log...")
			}),

			// Error state
			app.If(p.errMsg != "", func() app.UI {
				return app.Div().
					Style("background", "rgba(239,68,68,0.15)").
					Style("border", "1px solid rgba(239,68,68,0.4)").
					Style("border-radius", "8px").
					Style("padding", "12px 16px").
					Style("color", "#EF4444").
					Style("font-size", "14px").
					Style("margin-bottom", "20px").
					Text(p.errMsg)
			}),

			// Content
			app.If(!p.loading && p.errMsg == "", func() app.UI {
				if len(p.activities) == 0 {
					return app.Div().
						Style("background", "#1A1035").
						Style("border", "1px solid rgba(77,41,117,0.3)").
						Style("border-radius", "12px").
						Style("padding", "48px").
						Style("text-align", "center").
						Style("color", "#9B8DB5").
						Style("font-size", "14px").
						Text("Tidak ada aktivitas tercatat.")
				}
				return app.Div().
					Style("background", "#1A1035").
					Style("border", "1px solid rgba(77,41,117,0.3)").
					Style("border-radius", "12px").
					Style("padding", "20px 24px").
					Body(p.renderActivityItems())
			}),
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
			Style("padding", "14px 0").
			Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
			Body(
				// Dot indicator
				app.Div().
					Style("width", "8px").
					Style("height", "8px").
					Style("border-radius", "50%").
					Style("background", "#26B8B0").
					Style("flex-shrink", "0").
					Style("margin-top", "4px"),
				app.Div().
					Style("flex", "1").
					Style("min-width", "0").
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
							Style("margin-top", "3px").
							Text(fmt.Sprintf("oleh %s · %s", act.ActorName, act.CreatedAt)),
					),
			))
	}

	return app.Div().Body(items...)
}
