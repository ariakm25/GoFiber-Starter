package auth

import (
	"GoFiber-API/infra/response"
	"GoFiber-API/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func Register(ctx *fiber.Ctx) error {
	registerReq := &RegisterRequest{}

	if err := ctx.BodyParser(registerReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(registerReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	return AuthService(ctx).Register(*registerReq)

}

func Login(ctx *fiber.Ctx) error {

	loginReq := &LoginRequest{}

	if err := ctx.BodyParser(loginReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(loginReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	return AuthService(ctx).Login(*loginReq)
}

func RefreshToken(ctx *fiber.Ctx) error {
	refreshTokenReq := &RefreshTokenRequest{}

	if err := ctx.BodyParser(refreshTokenReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(refreshTokenReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	return AuthService(ctx).RefreshToken(*refreshTokenReq)
}

func ForgotPassword(ctx *fiber.Ctx) error {
	forgotPasswordReq := &ForgotPasswordRequest{}

	if err := ctx.BodyParser(forgotPasswordReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(forgotPasswordReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	return AuthService(ctx).ForgotPassword(*forgotPasswordReq)
}

func ValidateResetPasswordToken(ctx *fiber.Ctx) error {
	validateResetPasswordTokenReq := &ValidateResetPasswordTokenRequest{}

	if err := ctx.BodyParser(validateResetPasswordTokenReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(validateResetPasswordTokenReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	return AuthService(ctx).ValidateResetPasswordToken(*validateResetPasswordTokenReq)
}

func ResetPassword(ctx *fiber.Ctx) error {
	resetPasswordReq := &ResetPasswordRequest{}

	if err := ctx.BodyParser(resetPasswordReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(resetPasswordReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	return AuthService(ctx).ResetPassword(*resetPasswordReq)

}

func Logout(ctx *fiber.Ctx) error {
	logoutReq := &LogoutRequest{}

	if err := ctx.BodyParser(logoutReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	return AuthService(ctx).Logout(*logoutReq)
}

func Me(ctx *fiber.Ctx) error {
	return AuthService(ctx).Me()
}
