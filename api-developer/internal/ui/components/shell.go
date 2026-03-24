//go:build wasm

// Package components berisi komponen UI yang dapat digunakan ulang.
package components

import (
	"time"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// unreadCountResponse adalah response dari GET /api/internal/notifications/unread-count.
type unreadCountResponse struct {
	Count int `json:"count"`
}

// navItem merepresentasikan satu item navigasi di sidebar.
type navItem struct {
	Path  string
	Label string
	Icon  string // SVG path data
}

// Shell adalah layout utama dengan sidebar navigasi yang role-aware.
// Menampilkan sidebar kiri dengan nav items berdasarkan role user,
// dan content area di kanan untuk child component.
type Shell struct {
	app.Compo
	// Content adalah child component yang ditampilkan di area utama.
	Content     app.UI
	activeRoute string
	authStore   store.AuthStore
	notifCount  int
}

// OnNav dipanggil saat route berubah. Update activeRoute.
func (s *Shell) OnNav(ctx app.Context) {
	s.activeRoute = ctx.Page().URL().Path
}

// OnMount dipanggil saat komponen pertama kali dimount.
// Langsung poll notif count dan jadwalkan polling berikutnya.
func (s *Shell) OnMount(ctx app.Context) {
	s.pollNotifCount(ctx)
	s.scheduleNotifPoll(ctx)
}

// scheduleNotifPoll menjadwalkan polling notifikasi berikutnya setelah 30 detik.
func (s *Shell) scheduleNotifPoll(ctx app.Context) {
	ctx.After(30*time.Second, func(ctx app.Context) {
		s.pollNotifCount(ctx)
		s.scheduleNotifPoll(ctx)
	})
}

// pollNotifCount mengambil jumlah notifikasi yang belum dibaca dari API.
func (s *Shell) pollNotifCount(ctx app.Context) {
	token := s.authStore.GetToken()
	if token == "" {
		return
	}

	ctx.Async(func() {
		client := api.NewClient("", token)
		var resp unreadCountResponse
		if err := client.Get(ctx, "/api/internal/notifications/unread-count", &resp); err != nil {
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			s.notifCount = resp.Count
		})
	})
}

// navItems mengembalikan daftar nav items berdasarkan role user.
func (s *Shell) navItems() []navItem {
	items := []navItem{
		{
			Path:  "/",
			Label: "Dashboard",
			Icon:  "M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6",
		},
		{
			Path:  "/companies",
			Label: "Companies",
			Icon:  "M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4",
		},
		{
			Path:  "/licenses",
			Label: "Licenses",
			Icon:  "M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z",
		},
		{
			Path:  "/proposals",
			Label: "Proposals",
			Icon:  "M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z",
		},
		{
			Path:  "/notifications",
			Label: "Notifications",
			Icon:  "M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9",
		},
	}

	// Tambahkan Products dan Users hanya untuk superuser
	if s.authStore.HasRole("superuser") {
		items = append(items,
			navItem{
				Path:  "/products",
				Label: "Products",
				Icon:  "M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4",
			},
			navItem{
				Path:  "/users",
				Label: "Users",
				Icon:  "M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z",
			},
		)
	}

	return items
}

// Render menampilkan shell dengan sidebar dan content area.
func (s *Shell) Render() app.UI {
	user := s.authStore.GetUser()
	userName := "User"
	userRole := ""
	if user != nil {
		userName = user.Name
		userRole = user.Role
	}

	return app.Div().
		Style("display", "flex").
		Style("min-height", "100vh").
		Style("background", "#0F0A1A").
		Style("font-family", "'Inter', system-ui, sans-serif").
		Body(
			// Sidebar
			app.Nav().
				Style("width", "240px").
				Style("min-height", "100vh").
				Style("background", "#1A1035").
				Style("border-right", "1px solid rgba(77,41,117,0.4)").
				Style("display", "flex").
				Style("flex-direction", "column").
				Style("padding", "0").
				Style("flex-shrink", "0").
				Body(
					// Logo area
					app.Div().
						Style("padding", "20px 20px 16px").
						Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
						Body(
							app.Div().
								Style("display", "flex").
								Style("align-items", "center").
								Style("gap", "10px").
								Body(
									app.Div().
										Style("width", "32px").
										Style("height", "32px").
										Style("background", "linear-gradient(135deg, #4D2975, #26B8B0)").
										Style("border-radius", "8px").
										Style("display", "flex").
										Style("align-items", "center").
										Style("justify-content", "center").
										Style("color", "#E2D9F3").
										Style("font-weight", "700").
										Style("font-size", "16px").
										Text("V"),
									app.Span().
										Style("color", "#E2D9F3").
										Style("font-weight", "700").
										Style("font-size", "15px").
										Text("Vernon"),
								),
						),

					// Nav items
					app.Div().
						Style("flex", "1").
						Style("padding", "12px 8px").
						Style("overflow-y", "auto").
						Body(s.renderNavItems()...),

					// User info + logout
					app.Div().
						Style("padding", "12px 16px").
						Style("border-top", "1px solid rgba(77,41,117,0.3)").
						Body(
							app.Div().
								Style("display", "flex").
								Style("align-items", "center").
								Style("justify-content", "space-between").
								Body(
									app.Div().
										Body(
											app.Div().
												Style("color", "#E2D9F3").
												Style("font-size", "13px").
												Style("font-weight", "600").
												Text(userName),
											app.Div().
												Style("color", "#9B8DB5").
												Style("font-size", "11px").
												Style("text-transform", "capitalize").
												Text(userRole),
										),
									app.Button().
										Style("background", "none").
										Style("border", "1px solid rgba(155,141,181,0.3)").
										Style("border-radius", "6px").
										Style("padding", "6px 10px").
										Style("color", "#9B8DB5").
										Style("font-size", "12px").
										Style("cursor", "pointer").
										OnClick(s.onLogout).
										Text("Keluar"),
								),
						),
				),

			// Main content area
			app.Main().
				Style("flex", "1").
				Style("overflow-y", "auto").
				Style("min-height", "100vh").
				Body(s.Content),
		)
}

// renderNavItems mengembalikan slice UI untuk semua nav items.
func (s *Shell) renderNavItems() []app.UI {
	items := s.navItems()
	uis := make([]app.UI, 0, len(items))
	for _, item := range items {
		isActive := s.activeRoute == item.Path
		uis = append(uis, s.renderNavItem(item, isActive))
	}
	return uis
}

// renderNavItem merender satu nav item dengan styling active state.
func (s *Shell) renderNavItem(item navItem, isActive bool) app.UI {
	bgStyle := "transparent"
	colorStyle := "#9B8DB5"
	if isActive {
		bgStyle = "rgba(77,41,117,0.3)"
		colorStyle = "#E2D9F3"
	}

	return app.A().
		Href(item.Path).
		Style("display", "flex").
		Style("align-items", "center").
		Style("gap", "10px").
		Style("padding", "9px 12px").
		Style("border-radius", "8px").
		Style("margin-bottom", "2px").
		Style("background", bgStyle).
		Style("color", colorStyle).
		Style("text-decoration", "none").
		Style("font-size", "14px").
		Style("font-weight", func() string {
			if isActive {
				return "600"
			}
			return "400"
		}()).
		Style("transition", "background 0.15s, color 0.15s").
		OnClick(func(ctx app.Context, e app.Event) {
			e.PreventDefault()
			ctx.Navigate(item.Path)
		}).
		Body(
			// Icon (SVG)
			app.Raw(`<svg style="width:18px;height:18px;flex-shrink:0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" d="`+item.Icon+`"/></svg>`),
			app.Span().Text(item.Label),
			// Badge untuk notifications
			app.If(item.Path == "/notifications" && s.notifCount > 0,
				func() app.UI {
					return app.Span().
						Style("margin-left", "auto").
						Style("background", "#EF4444").
						Style("color", "#fff").
						Style("border-radius", "10px").
						Style("padding", "1px 7px").
						Style("font-size", "11px").
						Style("font-weight", "700").
						Text(s.notifCount)
				},
			),
		)
}

// onLogout membersihkan auth state dan redirect ke /login.
func (s *Shell) onLogout(ctx app.Context, e app.Event) {
	s.authStore.Clear()
	ctx.Navigate("/login")
}
