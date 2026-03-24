//go:build wasm

package pages

import (
	"context"
	"fmt"
	"strings"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ProposalDetail adalah DTO detail proposal yang dikembalikan dari API.
type ProposalDetail struct {
	ID               string         `json:"id"`
	ProjectID        string         `json:"project_id"`
	ProjectName      string         `json:"project_name"`
	CompanyID        string         `json:"company_id"`
	CompanyName      string         `json:"company_name"`
	ProductID        string         `json:"product_id"`
	ProductName      string         `json:"product_name"`
	Version          int            `json:"version"`
	Status           string         `json:"status"`
	Plan             string         `json:"plan"`
	Modules          []string       `json:"modules"`
	Apps             []string       `json:"apps"`
	MaxUsers         *int           `json:"max_users"`
	MaxTransPerMonth *int           `json:"max_trans_per_month"`
	MaxTransPerDay   *int           `json:"max_trans_per_day"`
	MaxItems         *int           `json:"max_items"`
	MaxCustomers     *int           `json:"max_customers"`
	MaxBranches      *int           `json:"max_branches"`
	MaxStorage       *int           `json:"max_storage"`
	ContractAmount   *float64       `json:"contract_amount"`
	ExpiresAt        *string        `json:"expires_at"`
	Notes            string         `json:"notes"`
	OwnerNotes       string         `json:"owner_notes"`
	RejectionReason  string         `json:"rejection_reason"`
	Changelog        *ChangelogView `json:"changelog"`
	PDFPath          string         `json:"pdf_path"`
	SubmittedByName  string         `json:"submitted_by_name"`
	ReviewedByName   string         `json:"reviewed_by_name"`
	ReviewedAt       *string        `json:"reviewed_at"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
}

// ChangelogView adalah tampilan changelog untuk UI.
type ChangelogView struct {
	ComparedToVersion int              `json:"compared_to_version"`
	Summary           string           `json:"summary"`
	Changes           []ChangelogEntry `json:"changes"`
	Unchanged         []string         `json:"unchanged"`
}

// ChangelogEntry adalah satu baris diff pada changelog.
type ChangelogEntry struct {
	Field    string `json:"field"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
}

// ProposalDetailPage menampilkan detail proposal dengan 2 tabs: Overview dan Changelog.
type ProposalDetailPage struct {
	app.Compo
	proposalID      string
	proposal        *ProposalDetail
	activeTab       string // "overview" | "changelog"
	loading         bool
	errMsg          string
	authStore       store.AuthStore
	showRejectModal bool
	rejectReason    string
	actionMsg       string
	actionErr       string
}

// OnNav dipanggil saat navigasi ke halaman ini. Mengambil proposalID dari URL.
func (p *ProposalDetailPage) OnNav(ctx app.Context) {
	p.proposalID = ctx.Page().URL().Path
	// Ekstrak ID dari path /proposals/{id}
	parts := strings.Split(strings.Trim(p.proposalID, "/"), "/")
	if len(parts) >= 2 {
		p.proposalID = parts[1]
	}

	if p.activeTab == "" {
		p.activeTab = "overview"
	}

	p.loading = true

	go p.loadProposal(ctx)
}

// loadProposal mengambil data proposal dari API.
func (p *ProposalDetailPage) loadProposal(ctx app.Context) {
	user := p.authStore.GetUser()
	if user == nil {
		ctx.Dispatch(func(ctx app.Context) {
			ctx.Navigate("/login")
		})
		return
	}

	client := api.NewClient("", user.Token)
	var proposal ProposalDetail
	if err := client.Get(context.Background(), "/api/internal/proposals/"+p.proposalID, &proposal); err != nil {
		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			p.errMsg = "Gagal memuat data proposal: " + err.Error()
		})
		return
	}

	ctx.Dispatch(func(ctx app.Context) {
		p.loading = false
		p.proposal = &proposal
	})
}

// onTabClick mengganti tab aktif.
func (p *ProposalDetailPage) onTabClick(tab string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		p.activeTab = tab
	}
}

// onSubmit mengubah status proposal dari draft → submitted.
func (p *ProposalDetailPage) onSubmit(ctx app.Context, e app.Event) {
	user := p.authStore.GetUser()
	if user == nil {
		return
	}
	p.actionMsg = ""
	p.actionErr = ""

	go func() {
		client := api.NewClient("", user.Token)
		if err := client.Put(context.Background(), "/api/internal/proposals/"+p.proposalID+"/submit", nil, nil); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.actionErr = "Gagal submit: " + err.Error()
			})
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.actionMsg = "Proposal berhasil di-submit untuk review."
			p.loading = true
			go p.loadProposal(ctx)
		})
	}()
}

// onApprove menyetujui proposal (PO/superuser).
func (p *ProposalDetailPage) onApprove(ctx app.Context, e app.Event) {
	user := p.authStore.GetUser()
	if user == nil {
		return
	}
	p.actionMsg = ""
	p.actionErr = ""

	go func() {
		client := api.NewClient("", user.Token)
		type approveReq struct{}
		var result map[string]any
		if err := client.Put(context.Background(), "/api/internal/proposals/"+p.proposalID+"/approve", approveReq{}, &result); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.actionErr = "Gagal approve: " + err.Error()
			})
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.actionMsg = "Proposal disetujui. Lisensi telah dibuat."
			p.loading = true
			go p.loadProposal(ctx)
		})
	}()
}

// onShowRejectModal menampilkan modal konfirmasi penolakan.
func (p *ProposalDetailPage) onShowRejectModal(ctx app.Context, e app.Event) {
	p.showRejectModal = true
	p.rejectReason = ""
}

// onHideRejectModal menyembunyikan modal penolakan.
func (p *ProposalDetailPage) onHideRejectModal(ctx app.Context, e app.Event) {
	p.showRejectModal = false
}

// onRejectReasonInput menangkap input alasan penolakan.
func (p *ProposalDetailPage) onRejectReasonInput(ctx app.Context, e app.Event) {
	p.rejectReason = ctx.JSSrc().Get("value").String()
}

// onConfirmReject mengirim penolakan ke API.
func (p *ProposalDetailPage) onConfirmReject(ctx app.Context, e app.Event) {
	if p.rejectReason == "" {
		return
	}
	user := p.authStore.GetUser()
	if user == nil {
		return
	}

	p.showRejectModal = false
	p.actionMsg = ""
	p.actionErr = ""

	reason := p.rejectReason
	go func() {
		client := api.NewClient("", user.Token)
		type rejectReq struct {
			Reason string `json:"reason"`
		}
		if err := client.Put(context.Background(), "/api/internal/proposals/"+p.proposalID+"/reject", rejectReq{Reason: reason}, nil); err != nil {
			ctx.Dispatch(func(ctx app.Context) {
				p.actionErr = "Gagal menolak: " + err.Error()
			})
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			p.actionMsg = "Proposal ditolak."
			p.loading = true
			go p.loadProposal(ctx)
		})
	}()
}

// onDownloadPDF membuka tab baru untuk download PDF.
func (p *ProposalDetailPage) onDownloadPDF(ctx app.Context, e app.Event) {
	user := p.authStore.GetUser()
	if user == nil {
		return
	}
	// Buka PDF di tab baru dengan token sebagai query param
	url := "/api/internal/proposals/" + p.proposalID + "/pdf?token=" + user.Token
	app.Window().Call("open", url, "_blank")
}

// Render menampilkan halaman detail proposal.
func (p *ProposalDetailPage) Render() app.UI {
	if p.loading {
		return p.renderLoading()
	}
	if p.errMsg != "" {
		return p.renderError()
	}
	if p.proposal == nil {
		return p.renderLoading()
	}

	return app.Div().
		Style("padding", "32px").
		Style("max-width", "960px").
		Style("margin", "0 auto").
		Body(
			p.renderHeader(),
			p.renderActionFeedback(),
			p.renderTabs(),
			app.If(p.activeTab == "overview", func() app.UI { return p.renderOverviewTab() }),
			app.If(p.activeTab == "changelog", func() app.UI { return p.renderChangelogTab() }),
			app.If(p.showRejectModal, func() app.UI { return p.renderRejectModal() }),
		)
}

// renderLoading menampilkan state loading.
func (p *ProposalDetailPage) renderLoading() app.UI {
	return app.Div().
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "center").
		Style("min-height", "200px").
		Style("color", "#9B8DB5").
		Style("font-size", "15px").
		Text("Memuat proposal...")
}

// renderError menampilkan pesan error.
func (p *ProposalDetailPage) renderError() app.UI {
	return app.Div().
		Style("padding", "24px").
		Style("color", "#EF4444").
		Text(p.errMsg)
}

// renderActionFeedback menampilkan pesan sukses atau error dari action.
func (p *ProposalDetailPage) renderActionFeedback() app.UI {
	if p.actionMsg == "" && p.actionErr == "" {
		return app.Text("")
	}

	color := "#22C55E"
	bg := "rgba(34,197,94,0.1)"
	msg := p.actionMsg
	if p.actionErr != "" {
		color = "#EF4444"
		bg = "rgba(239,68,68,0.1)"
		msg = p.actionErr
	}

	return app.Div().
		Style("background", bg).
		Style("border", "1px solid "+color).
		Style("border-radius", "8px").
		Style("padding", "12px 16px").
		Style("margin-bottom", "16px").
		Style("color", color).
		Style("font-size", "14px").
		Text(msg)
}

// renderHeader menampilkan header proposal dengan title dan status badge.
func (p *ProposalDetailPage) renderHeader() app.UI {
	return app.Div().
		Style("display", "flex").
		Style("align-items", "flex-start").
		Style("justify-content", "space-between").
		Style("margin-bottom", "24px").
		Body(
			app.Div().Body(
				app.Div().
					Style("display", "flex").
					Style("align-items", "center").
					Style("gap", "12px").
					Style("margin-bottom", "6px").
					Body(
						app.A().
							Href("/proposals").
							Style("color", "#9B8DB5").
							Style("text-decoration", "none").
							Style("font-size", "14px").
							OnClick(func(ctx app.Context, e app.Event) {
								e.PreventDefault()
								ctx.Navigate("/proposals")
							}).
							Text("← Proposals"),
					),
				app.Div().
					Style("display", "flex").
					Style("align-items", "center").
					Style("gap", "12px").
					Body(
						app.H1().
							Style("color", "#E2D9F3").
							Style("font-size", "24px").
							Style("font-weight", "700").
							Style("margin", "0").
							Text(fmt.Sprintf("Proposal v%d", p.proposal.Version)),
						renderProposalStatusBadge(p.proposal.Status),
					),
				app.Div().
					Style("color", "#9B8DB5").
					Style("font-size", "13px").
					Style("margin-top", "4px").
					Text(p.proposal.CompanyName+" · "+p.proposal.ProjectName),
			),
			p.renderActionButtons(),
		)
}

// renderProposalStatusBadge mengembalikan badge berwarna berdasarkan status.
func renderProposalStatusBadge(status string) app.UI {
	bg, color := "#9B8DB5", "#0F0A1A"
	switch status {
	case "draft":
		bg, color = "rgba(155,141,181,0.2)", "#9B8DB5"
	case "submitted":
		bg, color = "rgba(233,168,0,0.2)", "#E9A800"
	case "approved":
		bg, color = "rgba(34,197,94,0.2)", "#22C55E"
	case "rejected":
		bg, color = "rgba(239,68,68,0.2)", "#EF4444"
	}
	return app.Span().
		Style("background", bg).
		Style("color", color).
		Style("border-radius", "20px").
		Style("padding", "4px 12px").
		Style("font-size", "12px").
		Style("font-weight", "600").
		Style("text-transform", "capitalize").
		Text(status)
}

// renderActionButtons menampilkan tombol aksi berdasarkan status dan role.
func (p *ProposalDetailPage) renderActionButtons() app.UI {
	role := p.authStore.GetRole()
	status := p.proposal.Status

	buttons := []app.UI{}

	// Edit button: draft (semua role) atau submitted (PO/superuser)
	canEdit := status == "draft" ||
		(status == "submitted" && (role == "project_owner" || role == "superuser"))
	if canEdit {
		buttons = append(buttons,
			app.A().
				Href("/proposals/"+p.proposalID+"/edit").
				Style("display", "inline-flex").
				Style("align-items", "center").
				Style("padding", "8px 16px").
				Style("background", "transparent").
				Style("border", "1px solid rgba(77,41,117,0.5)").
				Style("border-radius", "8px").
				Style("color", "#E2D9F3").
				Style("font-size", "13px").
				Style("text-decoration", "none").
				Style("cursor", "pointer").
				OnClick(func(ctx app.Context, e app.Event) {
					e.PreventDefault()
					ctx.Navigate("/proposals/" + p.proposalID + "/edit")
				}).
				Text("Edit"),
		)
	}

	// Submit button: hanya draft
	if status == "draft" {
		buttons = append(buttons,
			app.Button().
				Style("padding", "8px 16px").
				Style("background", "#4D2975").
				Style("border", "none").
				Style("border-radius", "8px").
				Style("color", "#E2D9F3").
				Style("font-size", "13px").
				Style("cursor", "pointer").
				OnClick(p.onSubmit).
				Text("Submit"),
		)
	}

	// Approve + Reject: submitted + PO/superuser
	if status == "submitted" && (role == "project_owner" || role == "superuser") {
		buttons = append(buttons,
			app.Button().
				Style("padding", "8px 16px").
				Style("background", "rgba(34,197,94,0.15)").
				Style("border", "1px solid #22C55E").
				Style("border-radius", "8px").
				Style("color", "#22C55E").
				Style("font-size", "13px").
				Style("cursor", "pointer").
				OnClick(p.onApprove).
				Text("Approve"),
			app.Button().
				Style("padding", "8px 16px").
				Style("background", "rgba(239,68,68,0.15)").
				Style("border", "1px solid #EF4444").
				Style("border-radius", "8px").
				Style("color", "#EF4444").
				Style("font-size", "13px").
				Style("cursor", "pointer").
				OnClick(p.onShowRejectModal).
				Text("Reject"),
		)
	}

	// Download PDF: approved
	if status == "approved" {
		buttons = append(buttons,
			app.Button().
				Style("padding", "8px 16px").
				Style("background", "rgba(38,184,176,0.15)").
				Style("border", "1px solid #26B8B0").
				Style("border-radius", "8px").
				Style("color", "#26B8B0").
				Style("font-size", "13px").
				Style("cursor", "pointer").
				OnClick(p.onDownloadPDF).
				Text("Download PDF"),
		)
	}

	if len(buttons) == 0 {
		return app.Text("")
	}

	uis := make([]app.UI, len(buttons))
	copy(uis, buttons)

	return app.Div().
		Style("display", "flex").
		Style("gap", "10px").
		Style("align-items", "center").
		Body(uis...)
}

// renderTabs menampilkan tab switcher Overview / Changelog.
func (p *ProposalDetailPage) renderTabs() app.UI {
	tabs := []struct {
		key   string
		label string
	}{
		{"overview", "Overview"},
		{"changelog", "Changelog"},
	}

	uis := make([]app.UI, len(tabs))
	for i, t := range tabs {
		isActive := p.activeTab == t.key
		color := "#9B8DB5"
		borderBottom := "3px solid transparent"
		if isActive {
			color = "#E2D9F3"
			borderBottom = "3px solid #4D2975"
		}
		key := t.key
		uis[i] = app.Button().
			Style("background", "none").
			Style("border", "none").
			Style("border-bottom", borderBottom).
			Style("padding", "10px 20px").
			Style("color", color).
			Style("font-size", "14px").
			Style("font-weight", func() string {
				if isActive {
					return "600"
				}
				return "400"
			}()).
			Style("cursor", "pointer").
			Style("transition", "color 0.15s, border-bottom 0.15s").
			OnClick(p.onTabClick(key)).
			Text(t.label)
	}

	return app.Div().
		Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
		Style("margin-bottom", "24px").
		Style("display", "flex").
		Body(uis...)
}

// renderOverviewTab menampilkan tab overview proposal.
func (p *ProposalDetailPage) renderOverviewTab() app.UI {
	pr := p.proposal

	return app.Div().Body(
		// Info utama
		renderProposalSection("Informasi Dasar",
			renderProposalInfoGrid(
				"Product", pr.ProductName,
				"Plan", strings.ToUpper(pr.Plan),
				"Company", pr.CompanyName,
				"Project", pr.ProjectName,
			),
		),

		// Modules & Apps
		app.If(len(pr.Modules) > 0 || len(pr.Apps) > 0, func() app.UI {
			return renderProposalSection("Modules & Apps",
				app.Div().
					Style("display", "flex").
					Style("gap", "32px").
					Body(
						app.If(len(pr.Modules) > 0, func() app.UI {
							return renderChipList("Modules", pr.Modules, "#4D2975")
						}),
						app.If(len(pr.Apps) > 0, func() app.UI {
							return renderChipList("Apps", pr.Apps, "#26B8B0")
						}),
					),
			)
		}),

		// Constraints
		renderProposalSection("Constraints",
			renderConstraintsGrid(pr),
		),

		// Kontrak
		renderProposalSection("Kontrak",
			renderProposalInfoGrid(
				"Contract Amount", formatAmount(pr.ContractAmount),
				"Expires At", formatDateStr(pr.ExpiresAt),
			),
		),

		// Notes
		app.If(pr.Notes != "" || pr.OwnerNotes != "", func() app.UI {
			return renderProposalSection("Catatan",
				app.Div().Body(
					app.If(pr.Notes != "", func() app.UI {
						return renderNoteBlock("Catatan Sales", pr.Notes, "#9B8DB5")
					}),
					app.If(pr.OwnerNotes != "", func() app.UI {
						return renderNoteBlock("Catatan PO", pr.OwnerNotes, "#E9A800")
					}),
				),
			)
		}),

		// Rejection reason
		app.If(pr.Status == "rejected" && pr.RejectionReason != "", func() app.UI {
			return renderProposalSection("Alasan Penolakan",
				app.Div().
					Style("background", "rgba(239,68,68,0.1)").
					Style("border", "1px solid rgba(239,68,68,0.3)").
					Style("border-radius", "8px").
					Style("padding", "14px 16px").
					Style("color", "#EF4444").
					Style("font-size", "14px").
					Text(pr.RejectionReason),
			)
		}),

		// Review info
		app.If(pr.ReviewedByName != "", func() app.UI {
			return renderProposalSection("Review",
				renderProposalInfoGrid(
					"Direview oleh", pr.ReviewedByName,
					"Direview pada", formatDateStr(pr.ReviewedAt),
				),
			)
		}),

		// Submitted by
		renderProposalSection("Diajukan oleh",
			renderProposalInfoGrid(
				"Nama", pr.SubmittedByName,
				"Dibuat", pr.CreatedAt,
			),
		),
	)
}

// renderChangelogTab menampilkan tab changelog proposal.
func (p *ProposalDetailPage) renderChangelogTab() app.UI {
	if p.proposal.Changelog == nil || len(p.proposal.Changelog.Changes) == 0 {
		return app.Div().
			Style("background", "#1A1035").
			Style("border", "1px solid rgba(77,41,117,0.3)").
			Style("border-radius", "12px").
			Style("padding", "40px").
			Style("text-align", "center").
			Style("color", "#9B8DB5").
			Body(
				app.Div().Style("font-size", "32px").Style("margin-bottom", "12px").Text("📄"),
				app.Div().Style("font-size", "15px").Text("Ini adalah versi pertama proposal."),
				app.Div().Style("font-size", "13px").Style("margin-top", "8px").Text("Tidak ada perubahan yang dapat ditampilkan."),
			)
	}

	cl := p.proposal.Changelog

	return app.Div().Body(
		// Summary
		app.Div().
			Style("background", "rgba(77,41,117,0.15)").
			Style("border", "1px solid rgba(77,41,117,0.3)").
			Style("border-radius", "8px").
			Style("padding", "12px 16px").
			Style("margin-bottom", "20px").
			Style("color", "#E2D9F3").
			Style("font-size", "14px").
			Body(
				app.Span().Style("font-weight", "600").Text("Ringkasan: "),
				app.Span().Text(cl.Summary),
				app.Span().Style("color", "#9B8DB5").
					Style("margin-left", "16px").
					Style("font-size", "13px").
					Text(fmt.Sprintf("dibanding v%d", cl.ComparedToVersion)),
			),

		// Tabel perubahan
		app.Table().
			Style("width", "100%").
			Style("border-collapse", "collapse").
			Style("font-size", "14px").
			Body(
				app.THead().Body(
					app.Tr().
						Style("background", "rgba(77,41,117,0.2)").
						Body(
							app.Th().
								Style("text-align", "left").
								Style("padding", "10px 14px").
								Style("color", "#9B8DB5").
								Style("font-weight", "600").
								Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
								Text("Field"),
							app.Th().
								Style("text-align", "left").
								Style("padding", "10px 14px").
								Style("color", "#9B8DB5").
								Style("font-weight", "600").
								Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
								Text("Sebelumnya"),
							app.Th().
								Style("text-align", "left").
								Style("padding", "10px 14px").
								Style("color", "#9B8DB5").
								Style("font-weight", "600").
								Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
								Text("Sesudah"),
						),
				),
				app.TBody().Body(p.renderChangelogRows()...),
			),

		// Unchanged fields
		app.If(len(cl.Unchanged) > 0, func() app.UI {
			return app.Div().
				Style("margin-top", "20px").
				Body(
					app.Div().
						Style("color", "#9B8DB5").
						Style("font-size", "13px").
						Style("margin-bottom", "8px").
						Text(fmt.Sprintf("Tidak berubah: %s", strings.Join(cl.Unchanged, ", "))),
				)
		}),
	)
}

// renderChangelogRows menghasilkan baris-baris tabel changelog.
func (p *ProposalDetailPage) renderChangelogRows() []app.UI {
	rows := make([]app.UI, len(p.proposal.Changelog.Changes))
	for i, entry := range p.proposal.Changelog.Changes {
		rows[i] = app.Tr().
			Style("background", "rgba(233,168,0,0.06)").
			Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
			Body(
				app.Td().
					Style("padding", "10px 14px").
					Style("color", "#E9A800").
					Style("font-weight", "500").
					Text(formatFieldName(entry.Field)),
				app.Td().
					Style("padding", "10px 14px").
					Style("color", "#9B8DB5").
					Style("text-decoration", "line-through").
					Text(fmt.Sprintf("%v", entry.OldValue)),
				app.Td().
					Style("padding", "10px 14px").
					Style("color", "#E2D9F3").
					Style("font-weight", "500").
					Text(fmt.Sprintf("%v", entry.NewValue)),
			)
	}
	return rows
}

// renderRejectModal menampilkan modal untuk konfirmasi penolakan.
func (p *ProposalDetailPage) renderRejectModal() app.UI {
	return app.Div().
		Style("position", "fixed").
		Style("inset", "0").
		Style("background", "rgba(15,10,26,0.85)").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "center").
		Style("z-index", "100").
		Body(
			app.Div().
				Style("background", "#1A1035").
				Style("border", "1px solid rgba(239,68,68,0.4)").
				Style("border-radius", "12px").
				Style("padding", "28px").
				Style("width", "420px").
				Style("max-width", "90vw").
				Body(
					app.H3().
						Style("color", "#E2D9F3").
						Style("font-size", "18px").
						Style("margin", "0 0 8px").
						Text("Tolak Proposal"),
					app.P().
						Style("color", "#9B8DB5").
						Style("font-size", "14px").
						Style("margin", "0 0 16px").
						Text("Masukkan alasan penolakan. Alasan ini akan dikirimkan ke sales yang mengajukan proposal."),
					app.Textarea().
						Style("width", "100%").
						Style("background", "rgba(77,41,117,0.15)").
						Style("border", "1px solid rgba(77,41,117,0.4)").
						Style("border-radius", "8px").
						Style("padding", "10px 12px").
						Style("color", "#E2D9F3").
						Style("font-size", "14px").
						Style("min-height", "100px").
						Style("resize", "vertical").
						Style("box-sizing", "border-box").
						Placeholder("Tulis alasan penolakan...").
						OnInput(p.onRejectReasonInput),
					app.Div().
						Style("display", "flex").
						Style("gap", "10px").
						Style("justify-content", "flex-end").
						Style("margin-top", "16px").
						Body(
							app.Button().
								Style("padding", "8px 16px").
								Style("background", "transparent").
								Style("border", "1px solid rgba(155,141,181,0.3)").
								Style("border-radius", "8px").
								Style("color", "#9B8DB5").
								Style("font-size", "13px").
								Style("cursor", "pointer").
								OnClick(p.onHideRejectModal).
								Text("Batal"),
							app.Button().
								Style("padding", "8px 16px").
								Style("background", "rgba(239,68,68,0.2)").
								Style("border", "1px solid #EF4444").
								Style("border-radius", "8px").
								Style("color", "#EF4444").
								Style("font-size", "13px").
								Style("cursor", "pointer").
								OnClick(p.onConfirmReject).
								Text("Konfirmasi Tolak"),
						),
				),
		)
}

// ---- Helper render functions ----

// renderProposalSection membungkus konten dalam card section.
func renderProposalSection(title string, content app.UI) app.UI {
	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "20px").
		Style("margin-bottom", "16px").
		Body(
			app.Div().
				Style("color", "#9B8DB5").
				Style("font-size", "12px").
				Style("font-weight", "600").
				Style("text-transform", "uppercase").
				Style("letter-spacing", "0.08em").
				Style("margin-bottom", "14px").
				Text(title),
			content,
		)
}

// renderProposalInfoGrid menampilkan key-value pairs dalam 2 kolom. Argumen harus berpasangan.
func renderProposalInfoGrid(pairs ...string) app.UI {
	cells := make([]app.UI, 0, len(pairs))
	for i := 0; i+1 < len(pairs); i += 2 {
		label := pairs[i]
		value := pairs[i+1]
		cells = append(cells,
			app.Div().
				Style("display", "flex").
				Style("flex-direction", "column").
				Style("gap", "3px").
				Body(
					app.Span().
						Style("color", "#9B8DB5").
						Style("font-size", "12px").
						Text(label),
					app.Span().
						Style("color", "#E2D9F3").
						Style("font-size", "14px").
						Style("font-weight", "500").
						Text(value),
				),
		)
	}
	return app.Div().
		Style("display", "grid").
		Style("grid-template-columns", "repeat(auto-fill, minmax(200px, 1fr))").
		Style("gap", "16px").
		Body(cells...)
}

// renderConstraintsGrid menampilkan constraint fields dari proposal.
func renderConstraintsGrid(pr *ProposalDetail) app.UI {
	formatInt := func(v *int) string {
		if v == nil {
			return "—"
		}
		return fmt.Sprintf("%d", *v)
	}

	pairs := []string{
		"Max Users", formatInt(pr.MaxUsers),
		"Max Trans/Bulan", formatInt(pr.MaxTransPerMonth),
		"Max Trans/Hari", formatInt(pr.MaxTransPerDay),
		"Max Items", formatInt(pr.MaxItems),
		"Max Customers", formatInt(pr.MaxCustomers),
		"Max Branches", formatInt(pr.MaxBranches),
		"Max Storage (GB)", formatInt(pr.MaxStorage),
	}
	return renderProposalInfoGrid(pairs...)
}

// renderChipList menampilkan daftar item sebagai chip.
func renderChipList(label string, items []string, chipColor string) app.UI {
	chips := make([]app.UI, len(items))
	for i, item := range items {
		chips[i] = app.Span().
			Style("background", "rgba(77,41,117,0.2)").
			Style("border", "1px solid rgba(77,41,117,0.4)").
			Style("color", chipColor).
			Style("border-radius", "6px").
			Style("padding", "3px 10px").
			Style("font-size", "12px").
			Text(item)
	}
	return app.Div().Body(
		app.Div().
			Style("color", "#9B8DB5").
			Style("font-size", "12px").
			Style("margin-bottom", "8px").
			Text(label),
		app.Div().
			Style("display", "flex").
			Style("flex-wrap", "wrap").
			Style("gap", "6px").
			Body(chips...),
	)
}

// renderNoteBlock menampilkan blok catatan.
func renderNoteBlock(label, text, borderColor string) app.UI {
	return app.Div().
		Style("border-left", "3px solid "+borderColor).
		Style("padding", "10px 14px").
		Style("margin-bottom", "12px").
		Body(
			app.Div().
				Style("color", "#9B8DB5").
				Style("font-size", "12px").
				Style("margin-bottom", "4px").
				Text(label),
			app.Div().
				Style("color", "#E2D9F3").
				Style("font-size", "14px").
				Style("line-height", "1.6").
				Text(text),
		)
}

// formatAmount memformat angka float sebagai currency string.
func formatAmount(v *float64) string {
	if v == nil {
		return "—"
	}
	return fmt.Sprintf("%.2f", *v)
}

// formatDateStr memformat string tanggal RFC3339 menjadi lebih readable.
func formatDateStr(s *string) string {
	if s == nil || *s == "" {
		return "—"
	}
	// Ambil hanya bagian tanggal (YYYY-MM-DD)
	if len(*s) >= 10 {
		return (*s)[:10]
	}
	return *s
}

// formatFieldName mengubah snake_case field name menjadi Title Case.
func formatFieldName(field string) string {
	parts := strings.Split(field, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
