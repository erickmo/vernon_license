// Package proposal menyediakan generator PDF untuk proposal yang disetujui.
package proposal

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/flashlab/vernon-license/internal/domain"
)

// PDFData adalah data yang dibutuhkan untuk generate PDF proposal.
type PDFData struct {
	Proposal    *domain.Proposal
	CompanyName string
	ProjectName string
	ProductName string
	ReviewerName string
	// Vendor info dari config
	VendorName    string
	VendorAddress string
	VendorPhone   string
	VendorEmail   string
}

// proposalTmplData adalah data template internal yang sudah di-format untuk rendering.
type proposalTmplData struct {
	ProposalID       string
	GeneratedAt      string
	CompanyName      string
	ProjectName      string
	ProductName      string
	Plan             string
	ReviewerName     string
	VendorName       string
	VendorAddress    string
	VendorPhone      string
	VendorEmail      string
	MaxUsers         string
	MaxTransPerMonth string
	MaxTransPerDay   string
	MaxItems         string
	MaxCustomers     string
	MaxBranches      string
	MaxStorage       string
	Modules          string
	Apps             string
	ContractAmount   string
	ExpiresAt        string
	Notes            string
	OwnerNotes       string
	RejectionReason  string
	HasChangelog     bool
	ChangelogSummary string
	ChangelogEntries []changelogEntryDisplay
	ReviewedAt       string
}

type changelogEntryDisplay struct {
	Field    string
	OldValue string
	NewValue string
}

const proposalHTMLTemplate = `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Proposal Lisensi - {{.ProposalID}}</title>
<style>
  :root {
    --primary: #4D2975;
    --primary-light: #6B3FA0;
    --accent: #26B8B0;
    --amber: #E9A800;
    --success: #22C55E;
    --error: #EF4444;
    --gray-900: #111827;
    --gray-700: #374151;
    --gray-500: #6B7280;
    --gray-400: #9CA3AF;
    --gray-300: #D1D5DB;
    --gray-200: #E5E7EB;
    --gray-100: #F3F4F6;
    --white: #FFFFFF;
  }
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body {
    font-family: -apple-system, 'Segoe UI', Arial, sans-serif;
    background: #E5E7EB;
    color: var(--gray-900);
    padding: 40px 20px;
  }
  .page {
    width: 794px;
    min-height: 1123px;
    background: var(--white);
    margin: 0 auto;
    box-shadow: 0 4px 24px rgba(0,0,0,0.12);
    display: flex;
    flex-direction: column;
  }
  .accent-bar {
    height: 4px;
    background: var(--gray-300);
    flex-shrink: 0;
  }
  .header {
    padding: 40px 56px 32px;
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    border-bottom: 1px solid var(--gray-200);
  }
  .vendor-block { }
  .vendor-name {
    font-size: 15px;
    font-weight: 700;
    color: var(--gray-900);
    letter-spacing: -0.3px;
  }
  .vendor-detail {
    font-size: 12px;
    color: var(--gray-500);
    margin-top: 4px;
    line-height: 1.6;
  }
  .doc-info { text-align: right; }
  .doc-label {
    font-size: 11px;
    font-weight: 600;
    color: var(--gray-500);
    text-transform: uppercase;
    letter-spacing: 1px;
  }
  .doc-title {
    font-size: 13px;
    font-weight: 700;
    color: var(--gray-700);
    margin-top: 4px;
  }
  .doc-id {
    font-size: 12px;
    color: var(--gray-500);
    margin-top: 4px;
    font-family: monospace;
  }
  .doc-date {
    font-size: 12px;
    color: var(--gray-500);
    margin-top: 2px;
  }
  .content {
    padding: 40px 56px;
    flex: 1;
  }
  .section {
    margin-bottom: 32px;
  }
  .section-title {
    font-size: 11px;
    font-weight: 700;
    color: var(--gray-700);
    text-transform: uppercase;
    letter-spacing: 1px;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 2px solid var(--gray-300);
    display: inline-block;
  }
  table.info-table {
    width: 100%;
    border-collapse: collapse;
  }
  table.info-table tr td {
    padding: 8px 12px;
    font-size: 12px;
    border-bottom: 1px solid var(--gray-100);
  }
  table.info-table tr td:first-child {
    color: var(--gray-500);
    font-weight: 500;
    width: 200px;
  }
  table.info-table tr td:last-child {
    color: var(--gray-900);
    font-weight: 600;
  }
  table.info-table tr:last-child td {
    border-bottom: none;
  }
  .two-col {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 24px;
  }
  .info-card {
    background: var(--gray-100);
    border-radius: 8px;
    padding: 20px;
  }
  .info-card .card-label {
    font-size: 11px;
    font-weight: 600;
    color: var(--gray-500);
    text-transform: uppercase;
    letter-spacing: 0.8px;
    margin-bottom: 6px;
  }
  .info-card .card-value {
    font-size: 14px;
    font-weight: 700;
    color: var(--gray-900);
  }
  .info-card .card-sub {
    font-size: 12px;
    color: var(--gray-500);
    margin-top: 2px;
  }
  .notes-box {
    background: var(--gray-100);
    border-radius: 8px;
    padding: 14px 18px;
    font-size: 12px;
    color: var(--gray-700);
    line-height: 1.6;
    border-left: 3px solid var(--gray-400);
  }
  .owner-notes-box {
    background: var(--gray-100);
    border-radius: 8px;
    padding: 14px 18px;
    font-size: 12px;
    color: var(--gray-700);
    line-height: 1.6;
    border-left: 3px solid var(--gray-500);
  }
  .changelog-box {
    background: var(--gray-100);
    border-radius: 8px;
    padding: 14px 18px;
    border-left: 3px solid var(--gray-400);
  }
  .changelog-summary {
    font-size: 13px;
    font-weight: 600;
    color: var(--gray-700);
    margin-bottom: 10px;
  }
  .changelog-entry {
    font-size: 12px;
    color: var(--gray-700);
    padding: 6px 0;
    border-bottom: 1px solid rgba(0,0,0,0.06);
    display: grid;
    grid-template-columns: 160px 1fr 1fr;
    gap: 8px;
  }
  .changelog-entry:last-child { border-bottom: none; }
  .changelog-field { font-weight: 600; color: var(--gray-900); }
  .changelog-old { color: var(--error); text-decoration: line-through; }
  .changelog-new { color: var(--success); }
  .tag {
    display: inline-block;
    background: var(--gray-200);
    color: var(--gray-700);
    padding: 2px 8px;
    border-radius: 100px;
    font-size: 11px;
    font-weight: 600;
    margin: 2px;
  }
  .approved-badge {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: #DCFCE7;
    color: #15803D;
    padding: 4px 14px;
    border-radius: 100px;
    font-size: 12px;
    font-weight: 700;
  }
  .footer {
    padding: 32px 56px;
    border-top: 1px solid var(--gray-200);
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
  }
  .signature-block { }
  .signature-label {
    font-size: 11px;
    color: var(--gray-500);
    text-transform: uppercase;
    letter-spacing: 0.8px;
    margin-bottom: 48px;
  }
  .signature-line {
    width: 200px;
    height: 1px;
    background: var(--gray-300);
    margin-bottom: 8px;
  }
  .signature-name {
    font-size: 13px;
    font-weight: 600;
    color: var(--gray-900);
  }
  .signature-date {
    font-size: 12px;
    color: var(--gray-500);
    margin-top: 2px;
  }
  .footer-right {
    text-align: right;
    font-size: 11px;
    color: var(--gray-400);
  }
  .bottom-bar {
    height: 2px;
    background: var(--gray-200);
    flex-shrink: 0;
  }
</style>
</head>
<body>
<div class="page">
  <div class="accent-bar"></div>

  <div class="header">
    <div class="vendor-block">
      <div class="vendor-name">{{.VendorName}}</div>
      <div class="vendor-detail">
        {{.VendorAddress}}<br>
        {{.VendorPhone}} &middot; {{.VendorEmail}}
      </div>
    </div>
    <div class="doc-info">
      <div class="doc-label">Proposal Lisensi</div>
      <div class="doc-title">DISETUJUI</div>
      <div class="doc-id">{{.ProposalID}}</div>
      <div class="doc-date">Dibuat: {{.GeneratedAt}}</div>
    </div>
  </div>

  <div class="content">

    <!-- Klien & Proyek -->
    <div class="section">
      <div class="section-title">Informasi Klien & Proyek</div>
      <table class="info-table">
        <tr><td>Perusahaan</td><td>{{.CompanyName}}</td></tr>
        <tr><td>Proyek</td><td>{{.ProjectName}}</td></tr>
        <tr><td>Produk</td><td>{{.ProductName}}</td></tr>
        <tr><td>Plan</td><td>{{.Plan}}</td></tr>
      </table>
    </div>

    <!-- Constraints -->
    <div class="section">
      <div class="section-title">Constraint Lisensi</div>
      <table class="info-table">
        {{if .MaxUsers}}<tr><td>Maks. Pengguna</td><td>{{.MaxUsers}}</td></tr>{{end}}
        {{if .MaxTransPerMonth}}<tr><td>Maks. Transaksi / Bulan</td><td>{{.MaxTransPerMonth}}</td></tr>{{end}}
        {{if .MaxTransPerDay}}<tr><td>Maks. Transaksi / Hari</td><td>{{.MaxTransPerDay}}</td></tr>{{end}}
        {{if .MaxItems}}<tr><td>Maks. Item</td><td>{{.MaxItems}}</td></tr>{{end}}
        {{if .MaxCustomers}}<tr><td>Maks. Pelanggan</td><td>{{.MaxCustomers}}</td></tr>{{end}}
        {{if .MaxBranches}}<tr><td>Maks. Cabang</td><td>{{.MaxBranches}}</td></tr>{{end}}
        {{if .MaxStorage}}<tr><td>Maks. Storage</td><td>{{.MaxStorage}} MB</td></tr>{{end}}
      </table>
      {{if .Modules}}
      <div style="margin-top:12px">
        <div style="font-size:12px;color:var(--gray-500);margin-bottom:6px;">Modul Aktif:</div>
        {{range splitComma .Modules}}<span class="tag">{{.}}</span>{{end}}
      </div>
      {{end}}
      {{if .Apps}}
      <div style="margin-top:10px">
        <div style="font-size:12px;color:var(--gray-500);margin-bottom:6px;">Apps Aktif:</div>
        {{range splitComma .Apps}}<span class="tag">{{.}}</span>{{end}}
      </div>
      {{end}}
    </div>

    <!-- Kontrak -->
    <div class="section">
      <div class="section-title">Nilai Kontrak</div>
      <div class="two-col">
        <div class="info-card">
          <div class="card-label">Contract Amount</div>
          <div class="card-value">{{.ContractAmount}}</div>
          <div class="card-sub">Nilai kontrak tahunan</div>
        </div>
        <div class="info-card">
          <div class="card-label">Berlaku Hingga</div>
          <div class="card-value">{{.ExpiresAt}}</div>
          <div class="card-sub">Tanggal kadaluarsa lisensi</div>
        </div>
      </div>
    </div>

    <!-- Catatan Sales -->
    {{if .Notes}}
    <div class="section">
      <div class="section-title">Catatan Sales</div>
      <div class="notes-box">{{.Notes}}</div>
    </div>
    {{end}}

    <!-- Catatan PO -->
    {{if .OwnerNotes}}
    <div class="section">
      <div class="section-title">Catatan Project Owner</div>
      <div class="owner-notes-box">{{.OwnerNotes}}</div>
    </div>
    {{end}}

    <!-- Changelog -->
    {{if .HasChangelog}}
    <div class="section">
      <div class="section-title">Changelog</div>
      <div class="changelog-box">
        <div class="changelog-summary">{{.ChangelogSummary}}</div>
        {{range .ChangelogEntries}}
        <div class="changelog-entry">
          <span class="changelog-field">{{.Field}}</span>
          <span class="changelog-old">{{.OldValue}}</span>
          <span class="changelog-new">{{.NewValue}}</span>
        </div>
        {{end}}
      </div>
    </div>
    {{end}}

  </div>

  <div class="footer">
    <div class="signature-block">
      <div class="signature-label">Disetujui oleh</div>
      <div class="signature-line"></div>
      <div class="signature-name">{{.ReviewerName}}</div>
      <div class="signature-date">{{.ReviewedAt}}</div>
    </div>
    <div class="footer-right">
      <div><span class="approved-badge">&#10003; DISETUJUI</span></div>
      <div style="margin-top:8px;">Dokumen ini digenerate otomatis oleh Vernon License System</div>
      <div>{{.GeneratedAt}}</div>
    </div>
  </div>

  <div class="bottom-bar"></div>
</div>
</body>
</html>`

// GeneratePDF menghasilkan HTML proposal dalam format bytes.
// Output berisi: header perusahaan, detail proposal, constraints, tanda tangan reviewer.
// File disimpan sebagai HTML di STORAGE_PATH/proposals/{proposalID}.html.
// Return: bytes content, error
func GeneratePDF(data PDFData) ([]byte, error) {
	tmpl, err := template.New("proposal").Funcs(template.FuncMap{
		"splitComma": func(s string) []string {
			if s == "" {
				return nil
			}
			parts := strings.Split(s, ", ")
			var result []string
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					result = append(result, p)
				}
			}
			return result
		},
	}).Parse(proposalHTMLTemplate)
	if err != nil {
		return nil, fmt.Errorf("GeneratePDF parse template: %w", err)
	}

	p := data.Proposal
	td := proposalTmplData{
		ProposalID:   p.ID.String(),
		GeneratedAt:  time.Now().UTC().Format("02 January 2006, 15:04 UTC"),
		CompanyName:  data.CompanyName,
		ProjectName:  data.ProjectName,
		ProductName:  data.ProductName,
		Plan:         p.Plan,
		ReviewerName: data.ReviewerName,
		VendorName:   data.VendorName,
		VendorAddress: data.VendorAddress,
		VendorPhone:  data.VendorPhone,
		VendorEmail:  data.VendorEmail,
		Modules:      strings.Join(p.Modules, ", "),
		Apps:         strings.Join(p.Apps, ", "),
	}

	// Optional numeric fields
	if p.MaxUsers != nil {
		td.MaxUsers = fmt.Sprintf("%d", *p.MaxUsers)
	}
	if p.MaxTransPerMonth != nil {
		td.MaxTransPerMonth = fmt.Sprintf("%d", *p.MaxTransPerMonth)
	}
	if p.MaxTransPerDay != nil {
		td.MaxTransPerDay = fmt.Sprintf("%d", *p.MaxTransPerDay)
	}
	if p.MaxItems != nil {
		td.MaxItems = fmt.Sprintf("%d", *p.MaxItems)
	}
	if p.MaxCustomers != nil {
		td.MaxCustomers = fmt.Sprintf("%d", *p.MaxCustomers)
	}
	if p.MaxBranches != nil {
		td.MaxBranches = fmt.Sprintf("%d", *p.MaxBranches)
	}
	if p.MaxStorage != nil {
		td.MaxStorage = fmt.Sprintf("%d", *p.MaxStorage)
	}
	if p.ContractAmount != nil {
		td.ContractAmount = fmt.Sprintf("%.2f", *p.ContractAmount)
	} else {
		td.ContractAmount = "-"
	}
	if p.ExpiresAt != nil {
		td.ExpiresAt = p.ExpiresAt.UTC().Format("02 January 2006")
	} else {
		td.ExpiresAt = "Tidak ada batas waktu (perpetual)"
	}
	if p.Notes != nil {
		td.Notes = *p.Notes
	}
	if p.OwnerNotes != nil {
		td.OwnerNotes = *p.OwnerNotes
	}
	if p.ReviewedAt != nil {
		td.ReviewedAt = p.ReviewedAt.UTC().Format("02 January 2006")
	}

	// Changelog
	if len(p.Changelog) > 0 && string(p.Changelog) != "null" {
		cl, err := domain.ParseChangelog(p.Changelog)
		if err == nil && len(cl.Changes) > 0 {
			td.HasChangelog = true
			td.ChangelogSummary = cl.Summary
			for _, ch := range cl.Changes {
				td.ChangelogEntries = append(td.ChangelogEntries, changelogEntryDisplay{
					Field:    ch.Field,
					OldValue: fmt.Sprintf("%v", ch.OldValue),
					NewValue: fmt.Sprintf("%v", ch.NewValue),
				})
			}
		}
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, td); err != nil {
		return nil, fmt.Errorf("GeneratePDF execute template: %w", err)
	}
	return buf.Bytes(), nil
}

// SavePDF menyimpan konten HTML proposal ke STORAGE_PATH/proposals/{proposalID}.html.
// Direktori dibuat otomatis jika belum ada.
// Return: file path relatif (misal: "proposals/uuid.html"), error.
func SavePDF(storagePath string, proposalID string, content []byte) (string, error) {
	dir := filepath.Join(storagePath, "proposals")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("SavePDF mkdir: %w", err)
	}
	filename := proposalID + ".html"
	fullPath := filepath.Join(dir, filename)
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return "", fmt.Errorf("SavePDF write: %w", err)
	}
	return filepath.Join("proposals", filename), nil
}
