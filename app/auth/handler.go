package auth

import (
	user_module "GoFiber-API/app/user"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/infra/response"
	internal_log "GoFiber-API/internal/log"
	"GoFiber-API/internal/queue"
	"GoFiber-API/internal/utils"
	"errors"
	"time"

	pasetoware "github.com/gofiber/contrib/paseto"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

func Register(ctx *fiber.Ctx) error {
	registerReq := &RegisterRequest{}

	if err := ctx.BodyParser(registerReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(registerReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	hashedPassword, err := utils.HashPassword(registerReq.Password)

	if err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to register user"),
			response.WithError(response.ErrorInternal),
		).Send(ctx)
	}

	user := &user_module.User{
		Name:     registerReq.Name,
		Email:    registerReq.Email,
		Password: hashedPassword,
	}

	if err := database.Connection.First(user, "email = ?", registerReq.Email).Error; err == nil {
		return response.NewResponse(
			response.WithMessage("email already registered"),
			response.WithError(response.ErrorUnprocessableEntity),
		).Send(ctx)
	}

	if err := database.Connection.Create(user).Error; err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed register user"),
			response.WithError(response.ErrorInternal),
		).Send(ctx)
	}

	return response.NewResponse(
		response.WithMessage("success register"),
		response.WithData(user),
	).Send(ctx)

}

func Login(ctx *fiber.Ctx) error {

	loginReq := &LoginRequest{}

	if err := ctx.BodyParser(loginReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(loginReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	var user = &user_module.User{}

	findUser := database.Connection.First(user, "email = ?", loginReq.Email)

	if findUser.Error != nil {
		if errors.Is(findUser.Error, gorm.ErrRecordNotFound) {
			return response.NewResponse(
				response.WithMessage("invalid email or password"),
				response.WithError(response.ErrorUnauthorized),
			).Send(ctx)
		}

		return response.NewResponse(
			response.WithMessage(findUser.Error.Error()),
			response.WithError(response.ErrorUnauthorized),
		).Send(ctx)
	}

	if err := utils.ComparePassword(user.Password, loginReq.Password); err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or password"),
			response.WithError(response.ErrorUnauthorized),
		).Send(ctx)
	}

	token, err := utils.GenerateLocalPaseto(user.UID)

	if err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed generate token"),
			response.WithError(response.ErrorInternal),
		).Send(ctx)
	}

	return response.NewResponse(
		response.WithMessage("success login"),
		response.WithData(fiber.Map{
			"token": token,
		}),
	).Send(ctx)
}

func Me(ctx *fiber.Ctx) error {
	payload := ctx.Locals(pasetoware.DefaultContextKey).(*user_module.User)

	return response.NewResponse(
		response.WithMessage("success get user data"),
		response.WithData(payload),
	).Send(ctx)
}

func ForgotPassword(ctx *fiber.Ctx) error {
	resetPasswordReq := &ForgotPasswordRequest{}

	if err := ctx.BodyParser(resetPasswordReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(resetPasswordReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	var user user_module.User

	err := database.Connection.Where(&user_module.User{
		Email: resetPasswordReq.Email,
	}).First(&user).Error

	if err == nil {
		var checkResetPasswordToken user_module.UserToken

		database.Connection.Order("created_at desc").Where(&user_module.UserToken{
			UserID: user.UID,
			Type:   user_module.UserTokenTypeResetPassword,
		}).First(&checkResetPasswordToken)

		waitMin := time.Minute * 10

		diff := time.Since(checkResetPasswordToken.CreatedAt)

		if diff < waitMin {
			return response.NewResponse(
				response.WithMessage("Please wait 10 minutes before requesting a new password reset email."),
				response.WithError(response.ErrorBadRequest),
			).Send(ctx)
		}
	}

	task, err := NewAuthResetPasswordJob(resetPasswordReq.Email, user.UID)
	if err != nil {
		internal_log.Logger.Sugar().Errorf("NewAuthResetPasswordJob Error: %v", err)
	}

	_, err = queue.QueueClient.Enqueue(task, asynq.Retention(1*time.Hour))

	if err != nil {
		internal_log.Logger.Sugar().Errorf("NewAuthResetPasswordJob Enqueue Error: %v", err)
	}

	database.Connection.Delete(&user_module.UserToken{}, "user_id = ? AND type = ?", user.UID, user_module.UserTokenTypeResetPassword)

	return response.NewResponse(
		response.WithMessage("A password reset email will be sent if the email is registered in our system."),
	).Send(ctx)
}

func ValidateResetPasswordToken(ctx *fiber.Ctx) error {
	validateResetPasswordTokenReq := &ValidateResetPasswordTokenRequest{}

	if err := ctx.BodyParser(validateResetPasswordTokenReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(validateResetPasswordTokenReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	var user user_module.User

	err := database.Connection.Where(&user_module.User{
		Email: validateResetPasswordTokenReq.Email,
	}).First(&user).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or token"),
			response.WithError(response.ErrorNotFound),
		).Send(ctx)
	}

	var userToken user_module.UserToken

	err = database.Connection.Where(&user_module.UserToken{
		UserID: user.UID,
		Token:  validateResetPasswordTokenReq.Token,
		Type:   user_module.UserTokenTypeResetPassword,
	}).First(&userToken).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or token"),
			response.WithError(response.ErrorNotFound),
		).Send(ctx)
	}

	return response.NewResponse(
		response.WithMessage("success validate token"),
	).Send(ctx)
}

func ResetPassword(ctx *fiber.Ctx) error {
	resetPasswordReq := &ResetPasswordRequest{}

	if err := ctx.BodyParser(resetPasswordReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(resetPasswordReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	var user user_module.User

	err := database.Connection.Where(&user_module.User{
		Email: resetPasswordReq.Email,
	}).First(&user).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or token"),
			response.WithError(response.ErrorNotFound),
		).Send(ctx)
	}

	var userToken user_module.UserToken

	err = database.Connection.Where(&user_module.UserToken{
		UserID: user.UID,
		Token:  resetPasswordReq.Token,
		Type:   user_module.UserTokenTypeResetPassword,
	}).First(&userToken).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or token"),
			response.WithError(response.ErrorNotFound),
		).Send(ctx)
	}

	hashedPassword, err := utils.HashPassword(resetPasswordReq.Password)

	if err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to reset password"),
			response.WithError(response.ErrorInternal),
		).Send(ctx)
	}

	user.Password = hashedPassword

	if err := database.Connection.Save(&user).Error; err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to reset password"),
			response.WithError(response.ErrorInternal),
		).Send(ctx)
	}

	return response.NewResponse(
		response.WithMessage("success reset password"),
	).Send(ctx)

}
