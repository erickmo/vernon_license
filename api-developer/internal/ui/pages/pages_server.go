//go:build !wasm

// Package pages berisi semua halaman UI untuk Vernon App.
// File ini menyediakan stub types untuk server-side compilation.
// Implementasi lengkap ada di file dengan build tag wasm.
package pages

import "github.com/maxence-charriere/go-app/v10/pkg/app"

// SetupPage adalah halaman first-run setup superuser.
type SetupPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *SetupPage) Render() app.UI { return app.Div() }

// LoginPage adalah halaman login Vernon App.
type LoginPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *LoginPage) Render() app.UI { return app.Div() }

// DashboardPage adalah placeholder untuk dashboard.
type DashboardPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *DashboardPage) Render() app.UI { return app.Div() }

// CompaniesListPage menampilkan daftar semua perusahaan.
type CompaniesListPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *CompaniesListPage) Render() app.UI { return app.Div() }

// ProjectDetailPage menampilkan detail project.
type ProjectDetailPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *ProjectDetailPage) Render() app.UI { return app.Div() }

// LicensesListPage menampilkan daftar semua lisensi.
type LicensesListPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *LicensesListPage) Render() app.UI { return app.Div() }

// LicenseDetailPage menampilkan detail lisensi.
type LicenseDetailPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *LicenseDetailPage) Render() app.UI { return app.Div() }

// LicenseCreatePage adalah form untuk membuat lisensi baru.
type LicenseCreatePage struct{ app.Compo }

// Render stub untuk server-side.
func (p *LicenseCreatePage) Render() app.UI { return app.Div() }

// ProposalsListPage menampilkan daftar semua proposals.
type ProposalsListPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *ProposalsListPage) Render() app.UI { return app.Div() }

// ProposalDetailPage menampilkan detail proposal.
type ProposalDetailPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *ProposalDetailPage) Render() app.UI { return app.Div() }

// ProposalFormPage adalah form untuk membuat atau mengedit proposal.
type ProposalFormPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *ProposalFormPage) Render() app.UI { return app.Div() }

// ProductsListPage menampilkan daftar produk.
type ProductsListPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *ProductsListPage) Render() app.UI { return app.Div() }

// UsersListPage menampilkan daftar user.
type UsersListPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *UsersListPage) Render() app.UI { return app.Div() }

// NotificationsPage menampilkan daftar notifikasi.
type NotificationsPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *NotificationsPage) Render() app.UI { return app.Div() }

// ActivityLogPage menampilkan activity log sistem.
type ActivityLogPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *ActivityLogPage) Render() app.UI { return app.Div() }

// NotFoundPage adalah halaman 404.
type NotFoundPage struct{ app.Compo }

// Render stub untuk server-side.
func (p *NotFoundPage) Render() app.UI { return app.Div() }
