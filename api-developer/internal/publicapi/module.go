package publicapi

import "go.uber.org/fx"

// Module adalah Uber FX module yang menyediakan semua public API handlers.
// Handlers ini menangani 3 public endpoint: register, validate, dan validate_otp.
var Module = fx.Options(
	fx.Provide(NewRegisterHandler),
	fx.Provide(NewValidateHandler),
	fx.Provide(NewValidateOTPHandler),
)
