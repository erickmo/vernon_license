//go:build wasm

package pages

import (
	"context"
	"fmt"

	"github.com/flashlab/vernon-license/internal/ui/api"
	"github.com/flashlab/vernon-license/internal/ui/components"
	"github.com/flashlab/vernon-license/internal/ui/store"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// proposalListItem adalah DTO ringkas proposal untuk tampilan daftar.
type proposalListItem struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	CompanyID string `json:"company_id"`
	ProductID string `json:"product_id"`
	Version   int    `json:"version"`
	Status    string `json:"status"`
	Plan      string `json:"plan"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ProposalsListPage menampilkan daftar semua proposals dengan filter status.
type ProposalsListPage struct {
	app.Compo
	proposals    []proposalListItem
	loading      bool
	errMsg       string
	filterStatus string // "" = semua
	authStore    store.AuthStore
}

// OnNav dipanggil saat navigasi ke halaman proposals list.
func (p *ProposalsListPage) OnNav(ctx app.Context) {
	p.loading = true
	go p.loadProposals(ctx)
}

// loadProposals mengambil daftar proposals dari API.
func (p *ProposalsListPage) loadProposals(ctx app.Context) {
	user := p.authStore.GetUser()
	if user == nil {
		ctx.Dispatch(func(ctx app.Context) { ctx.Navigate("/login") })
		return
	}

	client := api.NewClient("", user.Token)
	var resp struct {
		Data []proposalListItem `json:"data"`
	}

	if err := client.Get(context.Background(), "/api/internal/proposals", &resp); err != nil {
		ctx.Dispatch(func(ctx app.Context) {
			p.loading = false
			p.errMsg = "Gagal memuat proposals: " + err.Error()
		})
		return
	}

	ctx.Dispatch(func(ctx app.Context) {
		p.loading = false
		p.proposals = resp.Data
	})
}

// filteredProposals mengembalikan proposals yang difilter berdasarkan status.
func (p *ProposalsListPage) filteredProposals() []proposalListItem {
	if p.filterStatus == "" {
		return p.proposals
	}
	var result []proposalListItem
	for _, pr := range p.proposals {
		if pr.Status == p.filterStatus {
			result = append(result, pr)
		}
	}
	return result
}

// onFilterChange menangani perubahan filter status.
func (p *ProposalsListPage) onFilterChange(status string) func(app.Context, app.Event) {
	return func(ctx app.Context, e app.Event) {
		p.filterStatus = status
	}
}

// Render menampilkan halaman daftar proposals dengan shell.
func (p *ProposalsListPage) Render() app.UI {
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

// renderContent merender area konten proposals list.
func (p *ProposalsListPage) renderContent() app.UI {
	return app.Div().
		Style("padding", "32px").
		Style("max-width", "960px").
		Style("margin", "0 auto").
		Body(
			// Header
			app.Div().
				Style("display", "flex").
				Style("align-items", "center").
				Style("justify-content", "space-between").
				Style("margin-bottom", "24px").
				Body(
					app.H1().
						Style("color", "#E2D9F3").
						Style("font-size", "24px").
						Style("font-weight", "700").
						Style("margin", "0").
						Text("Proposals"),
				),

			// Filter tabs
			p.renderFilterTabs(),

			// Content
			app.If(p.loading, func() app.UI {
				return app.Div().
					Style("text-align", "center").
					Style("padding", "60px").
					Style("color", "#9B8DB5").
					Text("Memuat proposals...")
			}),
			app.If(!p.loading && p.errMsg != "", func() app.UI {
				return app.Div().
					Style("color", "#EF4444").
					Style("padding", "24px").
					Text(p.errMsg)
			}),
			app.If(!p.loading && p.errMsg == "", func() app.UI {
				return p.renderProposalList()
			}),
		)
}

// renderFilterTabs menampilkan tab filter berdasarkan status.
func (p *ProposalsListPage) renderFilterTabs() app.UI {
	filters := []struct {
		key   string
		label string
	}{
		{"", "Semua"},
		{"draft", "Draft"},
		{"submitted", "Submitted"},
		{"approved", "Approved"},
		{"rejected", "Rejected"},
	}

	uis := make([]app.UI, len(filters))
	for i, f := range uis {
		_ = f
		isActive := p.filterStatus == filters[i].key
		color := "#9B8DB5"
		borderBottom := "3px solid transparent"
		if isActive {
			color = "#E2D9F3"
			borderBottom = "3px solid #4D2975"
		}
		key := filters[i].key
		uis[i] = app.Button().
			Style("background", "none").
			Style("border", "none").
			Style("border-bottom", borderBottom).
			Style("padding", "8px 16px").
			Style("color", color).
			Style("font-size", "14px").
			Style("cursor", "pointer").
			OnClick(p.onFilterChange(key)).
			Text(filters[i].label)
	}

	return app.Div().
		Style("border-bottom", "1px solid rgba(77,41,117,0.3)").
		Style("margin-bottom", "20px").
		Style("display", "flex").
		Body(uis...)
}

// renderProposalList menampilkan daftar proposal dalam bentuk tabel.
func (p *ProposalsListPage) renderProposalList() app.UI {
	filtered := p.filteredProposals()

	if len(filtered) == 0 {
		return app.Div().
			Style("background", "#1A1035").
			Style("border", "1px solid rgba(77,41,117,0.3)").
			Style("border-radius", "12px").
			Style("padding", "48px").
			Style("text-align", "center").
			Style("color", "#9B8DB5").
			Body(
				app.Div().Style("font-size", "14px").Text("Tidak ada proposal ditemukan."),
			)
	}

	rows := make([]app.UI, len(filtered))
	for i, pr := range filtered {
		proposal := pr
		rows[i] = app.Tr().
			Style("border-bottom", "1px solid rgba(77,41,117,0.15)").
			Style("cursor", "pointer").
			OnClick(func(ctx app.Context, e app.Event) {
				ctx.Navigate("/proposals/" + proposal.ID)
			}).
			Body(
				app.Td().
					Style("padding", "14px 16px").
					Style("color", "#E2D9F3").
					Style("font-weight", "500").
					Text(fmt.Sprintf("v%d", proposal.Version)),
				app.Td().
					Style("padding", "14px 16px").
					Style("color", "#9B8DB5").
					Style("font-size", "13px").
					Text(proposal.Plan),
				app.Td().
					Style("padding", "14px 16px").
					Body(renderProposalStatusBadge(proposal.Status)),
				app.Td().
					Style("padding", "14px 16px").
					Style("color", "#9B8DB5").
					Style("font-size", "13px").
					Text(formatDateStr(&proposal.CreatedAt)),
			)
	}

	return app.Table().
		Style("width", "100%").
		Style("border-collapse", "collapse").
		Style("background", "#1A1035").
		Style("border", "1px solid rgba(77,41,117,0.3)").
		Style("border-radius", "12px").
		Style("overflow", "hidden").
		Body(
			app.THead().Body(
				app.Tr().
					Style("background", "rgba(77,41,117,0.2)").
					Body(
						app.Th().Style("text-align", "left").Style("padding", "12px 16px").Style("color", "#9B8DB5").Style("font-weight", "600").Style("font-size", "12px").Text("VERSION"),
						app.Th().Style("text-align", "left").Style("padding", "12px 16px").Style("color", "#9B8DB5").Style("font-weight", "600").Style("font-size", "12px").Text("PLAN"),
						app.Th().Style("text-align", "left").Style("padding", "12px 16px").Style("color", "#9B8DB5").Style("font-weight", "600").Style("font-size", "12px").Text("STATUS"),
						app.Th().Style("text-align", "left").Style("padding", "12px 16px").Style("color", "#9B8DB5").Style("font-weight", "600").Style("font-size", "12px").Text("DIBUAT"),
					),
			),
			app.TBody().Body(rows...),
		)
}
