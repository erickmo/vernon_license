package service

import "go.uber.org/fx"

// ServiceModule adalah kumpulan provider Uber FX untuk semua service di layer ini.
// Gunakan module ini saat inisialisasi aplikasi dengan fx.New(...).
var ServiceModule = fx.Options(
	fx.Provide(NewCompanyService),
	fx.Provide(NewProjectService),
	fx.Provide(NewProductService),
	fx.Provide(NewLicenseService),
	fx.Provide(NewUserService),
	fx.Provide(NewAuthService),
	fx.Provide(NewAuditService),
	fx.Provide(NewNotificationService),
	fx.Provide(NewSetupService),
	fx.Provide(NewProposalService),
)
