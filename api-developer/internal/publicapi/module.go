package publicapi

import "go.uber.org/fx"

// Module adalah Uber FX module yang menyediakan semua public API handlers.
// Handlers ini menangani 2 public endpoint: register dan validate.
var Module = fx.Options(
	fx.Provide(NewRegisterHandler),
	fx.Provide(NewValidateHandler),
)
