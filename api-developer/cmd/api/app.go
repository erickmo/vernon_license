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
	app.Route("/licenses/create", func() app.Composer { return &pages.LicenseCreatePage{} })
	app.Route("/proposals/create", func() app.Composer { return &pages.ProposalFormPage{} })
	app.Route("/products", func() app.Composer { return &pages.ProductsListPage{} })
	app.Route("/users", func() app.Composer { return &pages.UsersListPage{} })
	app.Route("/notifications", func() app.Composer { return &pages.NotificationsPage{} })
	app.Route("/logs", func() app.Composer { return &pages.ActivityLogPage{} })
	app.RouteWithRegexp(`^/licenses/[^/]+$`, func() app.Composer { return &pages.LicenseDetailPage{} })
	app.RouteWithRegexp(`^/projects/[^/]+$`, func() app.Composer { return &pages.ProjectDetailPage{} })
	app.RouteWithRegexp(`^/proposals/[^/]+/edit$`, func() app.Composer { return &pages.ProposalFormPage{} })
	app.RouteWithRegexp(`^/proposals/[^/]+$`, func() app.Composer { return &pages.ProposalDetailPage{} })
	app.RouteWithRegexp(".*", func() app.Composer { return &pages.NotFoundPage{} })
}
