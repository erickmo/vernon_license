//go:build wasm

// Package main adalah WASM entry point untuk Vernon App.
// File ini hanya dikompilasi saat GOARCH=wasm GOOS=js.
package main

import (
	"github.com/flashlab/vernon-license/internal/ui/pages"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	app.Route("/setup", func() app.Composer { return &pages.SetupPage{} })
	app.Route("/login", func() app.Composer { return &pages.LoginPage{} })
	app.Route("/", func() app.Composer { return &pages.DashboardPage{} })
	app.Route("/companies", func() app.Composer { return &pages.CompaniesListPage{} })
	app.Route("/licenses", func() app.Composer { return &pages.LicensesListPage{} })
	app.Route("/licenses/create", func() app.Composer { return &pages.LicenseCreatePage{} })
	app.Route("/proposals/create", func() app.Composer { return &pages.ProposalFormPage{} })
	app.RouteWithRegexp(`^/licenses/[^/]+$`, func() app.Composer { return &pages.LicenseDetailPage{} })
	app.RouteWithRegexp(`^/projects/[^/]+$`, func() app.Composer { return &pages.ProjectDetailPage{} })
	app.RouteWithRegexp(`^/proposals/[^/]+/edit$`, func() app.Composer { return &pages.ProposalFormPage{} })
	app.RouteWithRegexp(`^/proposals/[^/]+$`, func() app.Composer { return &pages.ProposalDetailPage{} })
	app.RouteWithRegexp(`^/companies/[^/]+$`, func() app.Composer { return &pages.CompanyDetailPage{} })
	app.Route("/products", func() app.Composer { return &pages.ProductsListPage{} })
	app.RouteWithRegexp(`^/products/[^/]+$`, func() app.Composer { return &pages.ProductDetailPage{} })
	app.Route("/users", func() app.Composer { return &pages.UsersListPage{} })
	app.RouteWithRegexp(`^/users/[^/]+$`, func() app.Composer { return &pages.UserDetailPage{} })
	app.Route("/notifications", func() app.Composer { return &pages.NotificationsPage{} })
	app.Route("/logs", func() app.Composer { return &pages.ActivityLogPage{} })
	app.RouteWithRegexp(".*", func() app.Composer { return &pages.NotFoundPage{} })

	app.RunWhenOnBrowser()
}
