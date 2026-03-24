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

// ProjectDetail adalah representasi project detail dari API.
type ProjectDetail struct {
	ID          string  `json:"id"`
	CompanyID   string  `json:"company_id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
}

// companyDetailDTO adalah representasi company dari API (digunakan untuk breadcrumb).
type companyDetailDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// LicenseItem adalah representasi license dalam tab Licenses.
type LicenseItem struct {
	ID           string  `json:"id"`
	LicenseKey   string  `json:"license_key"`
	Plan         string  `json:"plan"`
	Status       string  `json:"status"`
	ExpiresAt    *string `json:"expires_at"`
	IsRegistered bool    `json:"is_registered"`
}

// licensesResponse adalah format response dari GET /api/internal/projects/{id}/licenses.
type licensesResponse struct {
	Data []LicenseItem `json:"data"`
}

// ProposalItem adalah representasi proposal dalam tab Proposals.
type ProposalItem struct {
	ID      string  `json:"id"`
	Version int     `json:"version"`
	Status  string  `json:"status"`
	Plan    string  `json:"plan"`
	Notes   *string `json:"notes"`
}

// AuditItem adalah representasi satu entri audit log.
type AuditItem struct {
	ID        string `json:"id"`
	Action    string `json:"action"`
	ActorName string `json:"actor_name"`
	CreatedAt string `json:"created_at"`
}

// auditResponse adalah format response dari audit log API.
type auditResponse struct {
	Data []AuditItem `json:"data"`
}

// ProjectDetailPage menampilkan detail project dengan 3 tabs: Licenses, Proposals, Activity.
type ProjectDetailPage struct {
	app.Compo
	authStore store.AuthStore

	projectID   string
	project     *ProjectDetail
	companyName string

	activeTab string // "licenses" | "proposals" | "activity"

	licenses  []LicenseItem
	proposals []ProposalItem
	auditLogs []AuditItem

	loading    bool
	errMsg     string

	// Project form (edit)
	showEditForm  bool
	editFormName  string
	editFormDesc  string
	editFormStatus string
	editFormErr   string
	editFormSaving bool

	// Project form (add project — tidak digunakan di detail, untuk kelengkapan)
	showProjectForm  bool
	projFormName     string
	projFormDesc     string
	projFormErr      string
	projFormSaving   bool
}

// OnNav dipanggil saat halaman ini di-navigasi.
// Ambil project ID dari URL, redirect ke /login jika belum login.
func (p *ProjectDetailPage) OnNav(ctx app.Context) {
	if !p.authStore.IsLoggedIn() {
		ctx.Navigate("/login")
		return
	}

	// Ambil ID dari URL: /projects/{id}
	path := ctx.Page().URL().Path
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == "projects" {
		p.projectID = parts[1]
	}

	if p.activeTab == "" {
		p.activeTab = "licenses"
	}

	p.fetchProject(ctx)
}

// fetchProject mengambil data project dari API.
func (p *ProjectDetailPage) fetchProject(ctx app.Context) {
	if p.projectID == "" {
		p.errMsg = "Project ID tidak ditemukan"
		return
	}

	p.loading = true
	p.errMsg = ""

	token := p.authStore.GetToken()
	projectID := p.projectID
	ctx.Async(func() {
		client := api.NewClient("", token)
		var project ProjectDetail
		err := client.Get(context.Background(), "/api/internal/projects/"+projectID, &project)
		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			if err != nil {
				p.errMsg = fmt.Sprintf("Gagal memuat project: %s", err.Error())
				return
			}
			p.project = &project
			p.fetchCompanyName(ctx, project.CompanyID)
			p.fetchTabData(ctx)
		})
	})
}

// fetchCompanyName mengambil nama company untuk breadcrumb.
func (p *ProjectDetailPage) fetchCompanyName(ctx app.Context, companyID string) {
	token := p.authStore.GetToken()
	ctx.Async(func() {
		client := api.NewClient("", token)
		var c companyDetailDTO
		err := client.Get(context.Background(), "/api/internal/companies/"+companyID, &c)
		ctx.Dispatch(func(ctx app.Context) {
			if err == nil {
				p.companyName = c.Name
			}
		})
	})
}

// fetchTabData mengambil data sesuai tab aktif.
func (p *ProjectDetailPage) fetchTabData(ctx app.Context) {
	switch p.activeTab {
	case "licenses":
		p.fetchLicenses(ctx)
	case "proposals":
		p.fetchProposals(ctx)
	case "activity":
		p.fetchAuditLogs(ctx)
	}
}

// fetchLicenses mengambil daftar licenses untuk project ini.
func (p *ProjectDetailPage) fetchLicenses(ctx app.Context) {
	token := p.authStore.GetToken()
	projectID := p.projectID
	ctx.Async(func() {
		client := api.NewClient("", token)
		var resp licensesResponse
		err := client.Get(context.Background(), "/api/internal/projects/"+projectID+"/licenses", &resp)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				app.Log("fetchLicenses error:", err.Error())
				return
			}
			p.licenses = resp.Data
		})
	})
}

// fetchProposals mengambil daftar proposals untuk project ini.
func (p *ProjectDetailPage) fetchProposals(ctx app.Context) {
	// Proposals belum ada endpoint spesifik per project — tampilkan kosong dengan pesan
	// Akan diisi di Phase 5.4
	p.proposals = []ProposalItem{}
}

// fetchAuditLogs mengambil audit log untuk project ini.
func (p *ProjectDetailPage) fetchAuditLogs(ctx app.Context) {
	token := p.authStore.GetToken()
	projectID := p.projectID
	ctx.Async(func() {
		client := api.NewClient("", token)
		var resp auditResponse
		err := client.Get(context.Background(), "/api/internal/projects/"+projectID+"/audit", &resp)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				app.Log("fetchAuditLogs error:", err.Error())
				p.auditLogs = []AuditItem{}
				return
			}
			p.auditLogs = resp.Data
		})
	})
}

// onTabClick menangani klik tab.
func (p *ProjectDetailPage) onTabClick(tab string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		e.PreventDefault()
		p.activeTab = tab
		p.fetchTabData(ctx)
	}
}

// onEditProject membuka form edit project.
func (p *ProjectDetailPage) onEditProject(ctx app.Context, e app.Event) {
	if p.project == nil {
		return
	}
	p.editFormName = p.project.Name
	p.editFormDesc = safeDeref(p.project.Description)
	p.editFormStatus = p.project.Status
	p.editFormErr = ""
	p.showEditForm = true
}

// onHideEditForm menutup form edit.
func (p *ProjectDetailPage) onHideEditForm(ctx app.Context, e app.Event) {
	p.showEditForm = false
	p.editFormErr = ""
}

// onEditFormNameChange menangani input nama.
func (p *ProjectDetailPage) onEditFormNameChange(ctx app.Context, e app.Event) {
	p.editFormName = ctx.JSSrc().Get("value").String()
}

// onEditFormDescChange menangani input deskripsi.
func (p *ProjectDetailPage) onEditFormDescChange(ctx app.Context, e app.Event) {
	p.editFormDesc = ctx.JSSrc().Get("value").String()
}

// onEditFormStatusChange menangani perubahan status.
func (p *ProjectDetailPage) onEditFormStatusChange(ctx app.Context, e app.Event) {
	p.editFormStatus = ctx.JSSrc().Get("value").String()
}

// onEditFormSubmit mengirim form edit project.
func (p *ProjectDetailPage) onEditFormSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()
	if strings.TrimSpace(p.editFormName) == "" {
		p.editFormErr = "Nama project wajib diisi"
		return
	}

	p.editFormErr = ""
	p.editFormSaving = true

	token := p.authStore.GetToken()
	projectID := p.projectID
	body := map[string]any{
		"name":        p.editFormName,
		"description": strPtr(p.editFormDesc),
		"status":      p.editFormStatus,
	}

	ctx.Async(func() {
		client := api.NewClient("", token)
		var resp ProjectDetail
		err := client.Put(context.Background(), "/api/internal/projects/"+projectID, body, &resp)
		ctx.Dispatch(func(ctx app.Context) {
			p.editFormSaving = false
			if err != nil {
				p.editFormErr = fmt.Sprintf("Gagal menyimpan: %s", err.Error())
				return
			}
			p.project = &resp
			p.showEditForm = false
		})
	})
}

// onCreateLicense navigasi ke halaman create license.
func (p *ProjectDetailPage) onCreateLicense(ctx app.Context, e app.Event) {
	ctx.Navigate("/licenses/create?project=" + p.projectID)
}

// onCreateProposal navigasi ke halaman create proposal.
func (p *ProjectDetailPage) onCreateProposal(ctx app.Context, e app.Event) {
	ctx.Navigate("/proposals/new?project=" + p.projectID)
}

// Render menampilkan halaman detail project.
func (p *ProjectDetailPage) Render() app.UI {
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
func (p *ProjectDetailPage) renderContent() app.UI {
	return app.Div().
		Style("padding", "32px").
		Body(
			// Breadcrumb
			p.renderBreadcrumb(),

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

			// Error state
			app.If(p.errMsg != "",
				func() app.UI {
					return app.Div().
						Style("background", "rgba(239,68,68,0.15)").
						Style("border", "1px solid rgba(239,68,68,0.4)").
						Style("border-radius", "8px").
						Style("padding", "12px 16px").
						Style("color", "#EF4444").
						Style("font-size", "14px").
						Style("margin-top", "16px").
						Text(p.errMsg)
				},
			),

			// Main content (project header + tabs)
			app.If(!p.loading && p.project != nil,
				func() app.UI {
					return app.Div().
						Style("margin-top", "16px").
						Body(
							p.renderProjectHeader(),
							p.renderTabs(),
							p.renderTabContent(),
						)
				},
			),

			// Edit form modal
			app.If(p.showEditForm,
				func() app.UI {
					return p.renderEditForm()
				},
			),
		)
}

// renderBreadcrumb merender navigasi breadcrumb.
func (p *ProjectDetailPage) renderBreadcrumb() app.UI {
	projectName := ""
	if p.project != nil {
		projectName = p.project.Name
	}

	return app.Div().
		Style("display", "flex").
		Style("align-items", "center").
		Style("gap", "8px").
		Style("margin-bottom", "8px").
		Body(
			app.A().
				Href("/companies").
				Style("color", "#9B8DB5").
				Style("font-size", "13px").
				Style("text-decoration", "none").
				OnClick(func(ctx app.Context, e app.Event) {
					e.PreventDefault()
					ctx.Navigate("/companies")
				}).
				Text("Companies"),
			app.Span().
				Style("color", "#6B5E8A").
				Style("font-size", "13px").
				Text("›"),
			app.If(p.companyName != "",
				func() app.UI {
					return app.Span().
						Style("color", "#9B8DB5").
						Style("font-size", "13px").
						Text(p.companyName)
				},
			),
			app.If(p.companyName != "",
				func() app.UI {
					return app.Span().
						Style("color", "#6B5E8A").
						Style("font-size", "13px").
						Text("›")
				},
			),
			app.Span().
				Style("color", "#E2D9F3").
				Style("font-size", "13px").
				Text(projectName),
		)
}

// renderProjectHeader merender header project dengan nama, status badge, dan actions.
func (p *ProjectDetailPage) renderProjectHeader() app.UI {
	if p.project == nil {
		return app.Div()
	}

	desc := ""
	if p.project.Description != nil {
		desc = *p.project.Description
	}

	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("padding", "24px").
		Style("margin-bottom", "24px").
		Body(
			app.Div().
				Style("display", "flex").
				Style("align-items", "flex-start").
				Style("justify-content", "space-between").
				Body(
					app.Div().
						Body(
							app.Div().
								Style("display", "flex").
								Style("align-items", "center").
								Style("gap", "12px").
								Style("margin-bottom", "8px").
								Body(
									app.H1().
										Style("color", "#E2D9F3").
										Style("font-size", "22px").
										Style("font-weight", "700").
										Style("margin", "0").
										Text(p.project.Name),
									renderStatusBadge(p.project.Status),
								),
							app.If(p.companyName != "",
								func() app.UI {
									return app.Div().
										Style("color", "#9B8DB5").
										Style("font-size", "13px").
										Style("margin-bottom", "6px").
										Text("Company: "+p.companyName)
								},
							),
							app.If(desc != "",
								func() app.UI {
									return app.P().
										Style("color", "#9B8DB5").
										Style("font-size", "13px").
										Style("margin", "0").
										Text(desc)
								},
							),
						),
					app.Button().
						Style("background", "rgba(77,41,117,0.3)").
						Style("color", "#E2D9F3").
						Style("border", "1px solid rgba(77,41,117,0.4)").
						Style("border-radius", "8px").
						Style("padding", "8px 16px").
						Style("font-size", "13px").
						Style("cursor", "pointer").
						Style("flex-shrink", "0").
						OnClick(p.onEditProject).
						Text("Edit Project"),
				),
		)
}

// renderTabs merender tab navigation.
func (p *ProjectDetailPage) renderTabs() app.UI {
	return app.Div().
		Style("display", "flex").
		Style("gap", "0").
		Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
		Style("margin-bottom", "24px").
		Body(
			p.renderTab("licenses", "Licenses"),
			p.renderTab("proposals", "Proposals"),
			p.renderTab("activity", "Activity"),
		)
}

// renderTab merender satu tab.
func (p *ProjectDetailPage) renderTab(key, label string) app.UI {
	isActive := p.activeTab == key
	color := "#9B8DB5"
	borderBottom := "3px solid transparent"
	fontWeight := "400"
	if isActive {
		color = "#26B8B0"
		borderBottom = "3px solid #26B8B0"
		fontWeight = "600"
	}

	return app.A().
		Href("#").
		Style("padding", "12px 20px").
		Style("color", color).
		Style("font-size", "14px").
		Style("font-weight", fontWeight).
		Style("text-decoration", "none").
		Style("border-bottom", borderBottom).
		Style("transition", "color 0.15s, border-color 0.15s").
		Style("cursor", "pointer").
		OnClick(p.onTabClick(key)).
		Text(label)
}

// renderTabContent merender konten sesuai tab aktif.
func (p *ProjectDetailPage) renderTabContent() app.UI {
	switch p.activeTab {
	case "licenses":
		return p.renderLicensesTab()
	case "proposals":
		return p.renderProposalsTab()
	case "activity":
		return p.renderActivityTab()
	default:
		return p.renderLicensesTab()
	}
}

// renderLicensesTab merender tab Licenses.
func (p *ProjectDetailPage) renderLicensesTab() app.UI {
	isPO := p.authStore.HasRole("project_owner")

	return app.Div().
		Body(
			// Action bar
			app.Div().
				Style("display", "flex").
				Style("justify-content", "space-between").
				Style("align-items", "center").
				Style("margin-bottom", "16px").
				Body(
					app.Div().
						Style("color", "#9B8DB5").
						Style("font-size", "13px").
						Text(fmt.Sprintf("%d license", len(p.licenses))),
					app.Div().
						Style("display", "flex").
						Style("gap", "8px").
						Body(
							app.Button().
								Style("background", "rgba(38,184,176,0.15)").
								Style("color", "#26B8B0").
								Style("border", "1px solid rgba(38,184,176,0.3)").
								Style("border-radius", "8px").
								Style("padding", "8px 16px").
								Style("font-size", "13px").
								Style("cursor", "pointer").
								OnClick(p.onCreateProposal).
								Text("Buat Proposal"),
							app.If(isPO,
								func() app.UI {
									return app.Button().
										Style("background", "#4D2975").
										Style("color", "#E2D9F3").
										Style("border", "none").
										Style("border-radius", "8px").
										Style("padding", "8px 16px").
										Style("font-size", "13px").
										Style("cursor", "pointer").
										OnClick(p.onCreateLicense).
										Text("Buat License")
								},
							),
						),
				),

			// License list / empty state
			app.If(len(p.licenses) == 0,
				func() app.UI {
					return renderEmptyState("Belum ada license", "Buat proposal atau license langsung untuk project ini")
				},
			),
			app.If(len(p.licenses) > 0,
				func() app.UI {
					return p.renderLicenseCards()
				},
			),
		)
}

// renderLicenseCards merender kartu-kartu license.
func (p *ProjectDetailPage) renderLicenseCards() app.UI {
	cards := make([]app.UI, 0, len(p.licenses))
	for _, l := range p.licenses {
		l := l
		cards = append(cards, p.renderLicenseCard(l))
	}
	return app.Div().
		Style("display", "flex").
		Style("flex-direction", "column").
		Style("gap", "12px").
		Body(cards...)
}

// renderLicenseCard merender satu kartu license.
func (p *ProjectDetailPage) renderLicenseCard(l LicenseItem) app.UI {
	expiresText := "Tidak ada batas"
	if l.ExpiresAt != nil && *l.ExpiresAt != "" {
		expiresText = "Expires: " + *l.ExpiresAt
	}

	registeredText := "Belum terdaftar"
	registeredColor := "#F59E0B"
	if l.IsRegistered {
		registeredText = "Terdaftar"
		registeredColor = "#22C55E"
	}

	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "10px").
		Style("padding", "16px 20px").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "space-between").
		Body(
			app.Div().
				Body(
					app.Div().
						Style("display", "flex").
						Style("align-items", "center").
						Style("gap", "10px").
						Style("margin-bottom", "6px").
						Body(
							app.Span().
								Style("color", "#E2D9F3").
								Style("font-size", "14px").
								Style("font-weight", "600").
								Style("font-family", "monospace").
								Text(l.LicenseKey),
							renderStatusBadge(l.Status),
						),
					app.Div().
						Style("display", "flex").
						Style("gap", "16px").
						Body(
							app.Span().
								Style("color", "#9B8DB5").
								Style("font-size", "12px").
								Text("Plan: "+l.Plan),
							app.Span().
								Style("color", "#9B8DB5").
								Style("font-size", "12px").
								Text(expiresText),
							app.Span().
								Style("color", registeredColor).
								Style("font-size", "12px").
								Text(registeredText),
						),
				),
			app.A().
				Href("/licenses/"+l.ID).
				Style("color", "#26B8B0").
				Style("font-size", "13px").
				Style("text-decoration", "none").
				OnClick(func(ctx app.Context, e app.Event) {
					e.PreventDefault()
					ctx.Navigate("/licenses/" + l.ID)
				}).
				Text("Detail →"),
		)
}

// renderProposalsTab merender tab Proposals.
func (p *ProjectDetailPage) renderProposalsTab() app.UI {
	return app.Div().
		Body(
			app.Div().
				Style("display", "flex").
				Style("justify-content", "flex-end").
				Style("margin-bottom", "16px").
				Body(
					app.Button().
						Style("background", "rgba(38,184,176,0.15)").
						Style("color", "#26B8B0").
						Style("border", "1px solid rgba(38,184,176,0.3)").
						Style("border-radius", "8px").
						Style("padding", "8px 16px").
						Style("font-size", "13px").
						Style("cursor", "pointer").
						OnClick(p.onCreateProposal).
						Text("Buat Proposal"),
				),

			app.If(len(p.proposals) == 0,
				func() app.UI {
					return renderEmptyState("Belum ada proposal", "Klik tombol Buat Proposal untuk membuat proposal baru")
				},
			),
			app.If(len(p.proposals) > 0,
				func() app.UI {
					return p.renderProposalCards()
				},
			),
		)
}

// renderProposalCards merender kartu proposal.
func (p *ProjectDetailPage) renderProposalCards() app.UI {
	cards := make([]app.UI, 0, len(p.proposals))
	for _, pr := range p.proposals {
		pr := pr
		cards = append(cards, p.renderProposalCard(pr))
	}
	return app.Div().
		Style("display", "flex").
		Style("flex-direction", "column").
		Style("gap", "12px").
		Body(cards...)
}

// renderProposalCard merender satu kartu proposal.
func (p *ProjectDetailPage) renderProposalCard(pr ProposalItem) app.UI {
	notes := ""
	if pr.Notes != nil {
		notes = *pr.Notes
	}

	return app.Div().
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "10px").
		Style("padding", "16px 20px").
		Style("display", "flex").
		Style("align-items", "center").
		Style("justify-content", "space-between").
		Body(
			app.Div().
				Body(
					app.Div().
						Style("display", "flex").
						Style("align-items", "center").
						Style("gap", "10px").
						Style("margin-bottom", "6px").
						Body(
							app.Span().
								Style("color", "#E2D9F3").
								Style("font-size", "14px").
								Style("font-weight", "600").
								Text(fmt.Sprintf("Proposal v%d", pr.Version)),
							renderStatusBadge(pr.Status),
						),
					app.Div().
						Style("display", "flex").
						Style("gap", "16px").
						Body(
							app.Span().
								Style("color", "#9B8DB5").
								Style("font-size", "12px").
								Text("Plan: "+pr.Plan),
							app.If(notes != "",
								func() app.UI {
									return app.Span().
										Style("color", "#9B8DB5").
										Style("font-size", "12px").
										Text(notes)
								},
							),
						),
				),
			app.A().
				Href("/proposals/"+pr.ID).
				Style("color", "#26B8B0").
				Style("font-size", "13px").
				Style("text-decoration", "none").
				OnClick(func(ctx app.Context, e app.Event) {
					e.PreventDefault()
					ctx.Navigate("/proposals/" + pr.ID)
				}).
				Text("Detail →"),
		)
}

// renderActivityTab merender tab Activity (audit log timeline).
func (p *ProjectDetailPage) renderActivityTab() app.UI {
	return app.Div().
		Body(
			app.If(len(p.auditLogs) == 0,
				func() app.UI {
					return renderEmptyState("Belum ada activity", "Activity akan muncul saat ada perubahan pada project ini")
				},
			),
			app.If(len(p.auditLogs) > 0,
				func() app.UI {
					return p.renderAuditTimeline()
				},
			),
		)
}

// renderAuditTimeline merender timeline audit log.
func (p *ProjectDetailPage) renderAuditTimeline() app.UI {
	entries := make([]app.UI, 0, len(p.auditLogs))
	for _, a := range p.auditLogs {
		a := a
		entries = append(entries, p.renderAuditEntry(a))
	}
	return app.Div().
		Style("display", "flex").
		Style("flex-direction", "column").
		Style("gap", "0").
		Body(entries...)
}

// renderAuditEntry merender satu entri audit log.
func (p *ProjectDetailPage) renderAuditEntry(a AuditItem) app.UI {
	return app.Div().
		Style("display", "flex").
		Style("gap", "16px").
		Style("padding", "16px 0").
		Style("border-bottom", "1px solid rgba(77,41,117,0.2)").
		Body(
			// Timeline dot
			app.Div().
				Style("flex-shrink", "0").
				Style("width", "10px").
				Style("height", "10px").
				Style("background", "#26B8B0").
				Style("border-radius", "50%").
				Style("margin-top", "4px"),
			// Content
			app.Div().
				Body(
					app.Div().
						Style("color", "#E2D9F3").
						Style("font-size", "14px").
						Style("font-weight", "500").
						Style("margin-bottom", "4px").
						Text(formatAction(a.Action)),
					app.Div().
						Style("display", "flex").
						Style("gap", "12px").
						Body(
							app.Span().
								Style("color", "#9B8DB5").
								Style("font-size", "12px").
								Text("oleh "+a.ActorName),
							app.Span().
								Style("color", "#6B5E8A").
								Style("font-size", "12px").
								Text(a.CreatedAt),
						),
				),
		)
}

// renderEditForm merender modal form edit project.
func (p *ProjectDetailPage) renderEditForm() app.UI {
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
				Style("max-width", "480px").
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
								Text("Edit Project"),
							app.Button().
								Style("background", "none").
								Style("border", "none").
								Style("color", "#9B8DB5").
								Style("cursor", "pointer").
								Style("font-size", "20px").
								Style("padding", "0").
								OnClick(p.onHideEditForm).
								Text("×"),
						),

					// Form error
					app.If(p.editFormErr != "",
						func() app.UI {
							return app.Div().
								Style("background", "rgba(239,68,68,0.15)").
								Style("border", "1px solid rgba(239,68,68,0.4)").
								Style("border-radius", "8px").
								Style("padding", "10px 14px").
								Style("color", "#EF4444").
								Style("font-size", "13px").
								Style("margin-bottom", "16px").
								Text(p.editFormErr)
						},
					),

					app.Form().
						OnSubmit(p.onEditFormSubmit).
						Body(
							renderFormField("Nama *", "text", p.editFormName, "Nama project", p.onEditFormNameChange),
							renderTextareaField("Deskripsi", p.editFormDesc, "Deskripsi project", p.onEditFormDescChange),

							// Status dropdown
							app.Div().
								Style("margin-bottom", "16px").
								Body(
									app.Label().
										Style("display", "block").
										Style("color", "#9B8DB5").
										Style("font-size", "12px").
										Style("font-weight", "600").
										Style("margin-bottom", "6px").
										Text("Status"),
									app.Select().
										Style("width", "100%").
										Style("background", "#0F0A1A").
										Style("border", "1px solid rgba(77,41,117,0.4)").
										Style("border-radius", "8px").
										Style("padding", "10px 14px").
										Style("color", "#E2D9F3").
										Style("font-size", "14px").
										Style("box-sizing", "border-box").
										OnChange(p.onEditFormStatusChange).
										Body(
											app.Option().Value("active").Selected(p.editFormStatus == "active").Text("Active"),
											app.Option().Value("completed").Selected(p.editFormStatus == "completed").Text("Completed"),
											app.Option().Value("cancelled").Selected(p.editFormStatus == "cancelled").Text("Cancelled"),
										),
								),

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
											if p.editFormSaving {
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
										OnClick(p.onHideEditForm).
										Text("Batal"),
								),
						),
				),
		)
}

// renderStatusBadge merender badge status dengan warna sesuai.
func renderStatusBadge(status string) app.UI {
	bg, textColor := statusColors(status)
	return app.Span().
		Style("background", bg).
		Style("color", textColor).
		Style("border-radius", "20px").
		Style("padding", "2px 10px").
		Style("font-size", "11px").
		Style("font-weight", "600").
		Style("text-transform", "capitalize").
		Text(status)
}

// statusColors mengembalikan pasangan warna background dan text untuk status.
func statusColors(status string) (bg, text string) {
	switch status {
	case "active", "approved":
		return "rgba(34,197,94,0.2)", "#22C55E"
	case "pending", "draft", "submitted", "trial":
		return "rgba(245,158,11,0.2)", "#F59E0B"
	case "suspended", "rejected", "cancelled":
		return "rgba(239,68,68,0.2)", "#EF4444"
	case "expired", "completed":
		return "rgba(155,141,181,0.2)", "#9B8DB5"
	default:
		return "rgba(77,41,117,0.2)", "#E2D9F3"
	}
}

// renderEmptyState merender tampilan kosong dengan ikon, judul, dan pesan.
func renderEmptyState(title, msg string) app.UI {
	return app.Div().
		Style("background", "#1A1035").
		Style("border-radius", "12px").
		Style("padding", "48px 32px").
		Style("text-align", "center").
		Body(
			app.Raw(`<svg style="width:40px;height:40px;color:#9B8DB5;margin:0 auto 12px;display:block" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0H4"/></svg>`),
			app.P().
				Style("color", "#9B8DB5").
				Style("font-size", "14px").
				Style("font-weight", "600").
				Style("margin", "0 0 6px").
				Text(title),
			app.P().
				Style("color", "#6B5E8A").
				Style("font-size", "13px").
				Style("margin", "0").
				Text(msg),
		)
}

// formatAction mengkonversi action string dari audit log menjadi teks yang readable.
func formatAction(action string) string {
	replacer := strings.NewReplacer(
		"_", " ",
		"company", "Company",
		"project", "Project",
		"license", "License",
		"proposal", "Proposal",
		"created", "dibuat",
		"updated", "diperbarui",
		"deleted", "dihapus",
		"approved", "disetujui",
		"rejected", "ditolak",
		"suspended", "disuspend",
		"activated", "diaktifkan",
	)
	return strings.Title(replacer.Replace(action))
}
