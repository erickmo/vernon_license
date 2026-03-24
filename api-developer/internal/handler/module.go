//go:build !wasm

// Package handler menyediakan Uber FX module untuk internal API handlers.
package handler

import "go.uber.org/fx"

// Module adalah Uber FX module yang menyediakan semua internal API handlers.
var Module = fx.Options(
	fx.Provide(NewAuthHandler),
	fx.Provide(NewSetupHandler),
	fx.Provide(NewCompanyHandler),
	fx.Provide(NewProjectHandler),
	fx.Provide(NewLicenseHandler),
	fx.Provide(NewProposalHandler),
	fx.Provide(NewProductHandler),
	fx.Provide(NewUserHandler),
	fx.Provide(NewNotificationHandler),
	fx.Provide(NewDashboardHandler),
)
