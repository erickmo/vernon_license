//go:build !wasm

// Package main — server-side route registration untuk Vernon App.
// go-app memerlukan routes yang sama terdaftar di server agar app.Handler
// dapat mengembalikan HTML yang benar untuk setiap path.
package main

import (
	"github.com/flashlab/vernon-license/internal/ui/pages"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func init() {
	app.Route("/setup", func() app.Composer { return &pages.SetupPage{} })
	app.Route("/login", func() app.Composer { return &pages.LoginPage{} })
	app.Route("/", func() app.Composer { return &pages.DashboardPage{} })
	app.Route("/companies", func() app.Composer { return &pages.CompaniesListPage{} })
	app.Route("/projects/{id}", func() app.Composer { return &pages.ProjectDetailPage{} })
	app.Route("/licenses", func() app.Composer { return &pages.LicensesListPage{} })
	app.Route("/licenses/{id}", func() app.Composer { return &pages.LicenseDetailPage{} })
	app.Route("/licenses/create", func() app.Composer { return &pages.LicenseCreatePage{} })
	app.Route("/proposals/{id}", func() app.Composer { return &pages.ProposalDetailPage{} })
	app.Route("/proposals/{id}/edit", func() app.Composer { return &pages.ProposalFormPage{} })
	app.Route("/products", func() app.Composer { return &pages.ProductsListPage{} })
	app.Route("/users", func() app.Composer { return &pages.UsersListPage{} })
	app.Route("/notifications", func() app.Composer { return &pages.NotificationsPage{} })
	app.RouteWithRegexp(".*", func() app.Composer { return &pages.NotFoundPage{} })
}
