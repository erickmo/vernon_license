//go:build wasm

// Package main adalah WASM entry point untuk Vernon App.
// File ini hanya dikompilasi saat GOARCH=wasm GOOS=js.
package main

import (
	"github.com/flashlab/vernon-license/internal/ui/pages"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	app.Route("/login", func() app.Composer { return &pages.LoginPage{} })
	app.Route("/", func() app.Composer { return &pages.DashboardPage{} })
	app.Route("/companies", func() app.Composer { return &pages.CompaniesListPage{} })
	app.Route("/projects/{id}", func() app.Composer { return &pages.ProjectDetailPage{} })
	app.Route("/licenses", func() app.Composer { return &pages.LicensesListPage{} })
	app.Route("/licenses/{id}", func() app.Composer { return &pages.LicenseDetailPage{} })
	app.Route("/licenses/create", func() app.Composer { return &pages.LicenseCreatePage{} })
	app.Route("/proposals", func() app.Composer { return &pages.ProposalsListPage{} })
	app.Route("/proposals/create", func() app.Composer { return &pages.ProposalFormPage{} })
	app.Route("/proposals/{id}", func() app.Composer { return &pages.ProposalDetailPage{} })
	app.Route("/proposals/{id}/edit", func() app.Composer { return &pages.ProposalFormPage{} })
	app.Route("/products", func() app.Composer { return &pages.ProductsListPage{} })
	app.Route("/users", func() app.Composer { return &pages.UsersListPage{} })
	app.Route("/notifications", func() app.Composer { return &pages.NotificationsPage{} })
	app.RouteWithRegexp(".*", func() app.Composer { return &pages.NotFoundPage{} })

	app.RunWhenOnBrowser()
}
