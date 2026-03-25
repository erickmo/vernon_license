//go:build wasm

package pages

import (
	"fmt"
	"time"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/components"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// NotificationItem merepresentasikan satu notifikasi.
type NotificationItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

// notificationsListResponse adalah response dari GET /api/internal/notifications.
type notificationsListResponse struct {
	Data []NotificationItem `json:"data"`
}

// NotificationsPage menampilkan daftar notifikasi user.
type NotificationsPage struct {
	app.Compo
	notifications []NotificationItem
	loading       bool
	markingAll    bool
	errMsg        string
	authStore     store.AuthStore
}

// OnNav dipanggil saat halaman ini di-navigasi.
func (p *NotificationsPage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}
	p.loadNotifications(ctx)
}

// loadNotifications mengambil daftar notifikasi dari API.
func (p *NotificationsPage) loadNotifications(ctx app.Context) {
	p.loading = true
	p.errMsg = ""

	token := p.authStore.GetToken()

	ctx.Async(func() {
		client := api.NewClient("", token)
		var resp notificationsListResponse
		err := client.Get(ctx, "/api/internal/notifications", &resp)

		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			if err != nil {
				p.errMsg = "Gagal memuat notifikasi."
				return
			}
			p.notifications = resp.Data
		})
	})
}

// onMarkAllRead menandai semua notifikasi sebagai sudah dibaca.
func (p *NotificationsPage) onMarkAllRead(ctx app.Context, e app.Event) {
	if p.markingAll {
		return
	}
	p.markingAll = true

	token := p.authStore.GetToken()

	ctx.Async(func() {
		client := api.NewClient("", token)
		err := client.Put(ctx, "/api/internal/notifications/read-all", nil, nil)

		ctx.Dispatch(func(ctx app.Context) {
			p.markingAll = false
			if err != nil {
				p.errMsg = "Gagal menandai semua dibaca."
				return
			}
			p.loadNotifications(ctx)
		})
	})
}

// onMarkRead menandai satu notifikasi sebagai sudah dibaca.
func (p *NotificationsPage) onMarkRead(id string, isRead bool) func(app.Context, app.Event) {
	return func(ctx app.Context, e app.Event) {
		if isRead {
			return
		}
		token := p.authStore.GetToken()

		ctx.Async(func() {
			client := api.NewClient("", token)
			err := client.Put(ctx, "/api/internal/notifications/"+id+"/read", nil, nil)

			ctx.Dispatch(func(ctx app.Context) {
				if err != nil {
					return
				}
				// Update state lokal
				for i := range p.notifications {
					if p.notifications[i].ID == id {
						p.notifications[i].IsRead = true
						break
					}
				}
			})
		})
	}
}

// Render menampilkan halaman notifications.
func (p *NotificationsPage) Render() app.UI {
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

// renderContent merender area konten notifications.
func (p *NotificationsPage) renderContent() app.UI {
	hasUnread := false
	for _, n := range p.notifications {
		if !n.IsRead {
			hasUnread = true
			break
		}
	}

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
								Text("Notifications"),
							app.P().
								Style("color", "#9B8DB5").
								Style("font-size", "14px").
								Style("margin", "0").
								Text("Notifikasi sistem Vernon"),
						),
					app.If(hasUnread,
						func() app.UI {
							return app.Button().
								Disabled(p.markingAll).
								Style("background", "rgba(77,41,117,0.3)").
								Style("color", "#E2D9F3").
								Style("border", "1px solid rgba(77,41,117,0.5)").
								Style("border-radius", "8px").
								Style("padding", "8px 16px").
								Style("font-size", "13px").
								Style("cursor", func() string {
									if p.markingAll {
										return "not-allowed"
									}
									return "pointer"
								}()).
								OnClick(p.onMarkAllRead).
								Text(func() string {
									if p.markingAll {
										return "Memproses..."
									}
									return "Tandai Semua Dibaca"
								}())
						},
					),
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
						Text("Memuat notifikasi...")
				},
			),

			// List atau empty state
			app.If(!p.loading,
				func() app.UI {
					return p.renderList()
				},
			),
		)
}

// renderList merender daftar notifikasi.
func (p *NotificationsPage) renderList() app.UI {
	if len(p.notifications) == 0 {
		return app.Div().
			Style("text-align", "center").
			Style("color", "#9B8DB5").
			Style("padding", "64px 0").
			Body(
				app.Div().
					Style("font-size", "40px").
					Style("margin-bottom", "12px").
					Text("🔔"),
				app.Div().
					Style("font-size", "16px").
					Style("font-weight", "500").
					Style("color", "#E2D9F3").
					Style("margin-bottom", "8px").
					Text("Tidak ada notifikasi"),
				app.Div().
					Style("font-size", "14px").
					Text("Notifikasi akan muncul di sini saat ada aktivitas terkait akun Anda."),
			)
	}

	items := make([]app.UI, 0, len(p.notifications))
	for _, notif := range p.notifications {
		notif := notif
		items = append(items, p.renderNotifItem(notif))
	}

	return app.Div().
		Style("width", "100%").
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("overflow", "hidden").
		Body(items...)
}

// renderNotifItem merender satu item notifikasi.
func (p *NotificationsPage) renderNotifItem(notif NotificationItem) app.UI {
	bgStyle := "transparent"
	if !notif.IsRead {
		bgStyle = "rgba(77,41,117,0.1)"
	}

	timeStr := formatNotifTime(notif.CreatedAt)

	return app.Div().
		Style("display", "flex").
		Style("align-items", "flex-start").
		Style("gap", "14px").
		Style("padding", "16px 20px").
		Style("border-bottom", "1px solid rgba(77,41,117,0.2)").
		Style("background", bgStyle).
		Style("cursor", func() string {
			if !notif.IsRead {
				return "pointer"
			}
			return "default"
		}()).
		Style("transition", "background 0.15s").
		OnClick(p.onMarkRead(notif.ID, notif.IsRead)).
		Body(
			// Icon berdasarkan type
			app.Div().
				Style("width", "36px").
				Style("height", "36px").
				Style("border-radius", "50%").
				Style("display", "flex").
				Style("align-items", "center").
				Style("justify-content", "center").
				Style("flex-shrink", "0").
				Style("background", notifIconBg(notif.Type)).
				Style("color", notifIconColor(notif.Type)).
				Body(
					app.Raw(`<svg style="width:16px;height:16px" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="`+notifIconPath(notif.Type)+`"/></svg>`),
				),

			// Content
			app.Div().
				Style("flex", "1").
				Style("min-width", "0").
				Body(
					app.Div().
						Style("display", "flex").
						Style("align-items", "flex-start").
						Style("justify-content", "space-between").
						Style("gap", "8px").
						Body(
							app.Div().
								Style("color", "#E2D9F3").
								Style("font-size", "14px").
								Style("font-weight", func() string {
									if !notif.IsRead {
										return "600"
									}
									return "400"
								}()).
								Text(notif.Title),
							// Unread dot
							app.If(!notif.IsRead,
								func() app.UI {
									return app.Div().
										Style("width", "8px").
										Style("height", "8px").
										Style("border-radius", "50%").
										Style("background", "#EF4444").
										Style("flex-shrink", "0").
										Style("margin-top", "4px")
								},
							),
						),
					app.Div().
						Style("color", "#9B8DB5").
						Style("font-size", "13px").
						Style("margin-top", "3px").
						Text(notif.Body),
					app.Div().
						Style("color", "#9B8DB5").
						Style("font-size", "12px").
						Style("margin-top", "6px").
						Text(timeStr),
				),
		)
}

// notifIconPath mengembalikan SVG path icon berdasarkan notification type.
func notifIconPath(notifType string) string {
	switch notifType {
	case "proposal_submitted":
		return "M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
	case "proposal_approved":
		return "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
	case "proposal_rejected":
		return "M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
	case "license_expiring":
		return "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
	case "license_suspended":
		return "M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
	case "client_registered":
		return "M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"
	default:
		return "M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
	}
}

// notifIconBg mengembalikan background color berdasarkan notification type.
func notifIconBg(notifType string) string {
	switch notifType {
	case "proposal_approved", "client_registered":
		return "rgba(34,197,94,0.15)"
	case "proposal_rejected", "license_suspended":
		return "rgba(239,68,68,0.15)"
	case "license_expiring":
		return "rgba(245,158,11,0.15)"
	default:
		return "rgba(77,41,117,0.2)"
	}
}

// notifIconColor mengembalikan icon color berdasarkan notification type.
func notifIconColor(notifType string) string {
	switch notifType {
	case "proposal_approved", "client_registered":
		return "#22C55E"
	case "proposal_rejected", "license_suspended":
		return "#EF4444"
	case "license_expiring":
		return "#F59E0B"
	default:
		return "#26B8B0"
	}
}

// formatNotifTime mem-format ISO timestamp menjadi string yang mudah dibaca.
func formatNotifTime(isoTime string) string {
	t, err := time.Parse("2006-01-02T15:04:05Z", isoTime)
	if err != nil {
		return isoTime
	}
	now := time.Now().UTC()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "Baru saja"
	case diff < time.Hour:
		return fmt.Sprintf("%d menit yang lalu", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d jam yang lalu", int(diff.Hours()))
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("%d hari yang lalu", int(diff.Hours()/24))
	default:
		return t.Format("02 Jan 2006")
	}
}
