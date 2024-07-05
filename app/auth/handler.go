package auth

import (
	"GoFiber-API/app/user"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/infra/response"
	internal_log "GoFiber-API/internal/log"
	"GoFiber-API/internal/queue"
	"GoFiber-API/internal/utils"
	"errors"
	"time"

	pasetoware "github.com/gofiber/contrib/paseto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	user := &user.User{
		Name:     registerReq.Name,
		Email:    registerReq.Email,
		Password: hashedPassword,
		UID:      uuid.New().String(),
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

	var user = &user.User{}

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
	payload := ctx.Locals(pasetoware.DefaultContextKey).(*user.User)

	return response.NewResponse(
		response.WithMessage("success get user data"),
		response.WithData(payload),
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

	var user = &user.User{}

	findUser := database.Connection.First(user, "email = ?", resetPasswordReq.Email)

	if findUser.Error == nil {
		task, err := NewAuthResetPasswordJob(resetPasswordReq.Email)
		if err != nil {
			internal_log.Logger.Sugar().Errorf("NewAuthResetPasswordJob Error: %v", err)
		}

		_, err = queue.QueueClient.Enqueue(task, asynq.Retention(1*time.Hour))

		if err != nil {
			internal_log.Logger.Sugar().Errorf("NewAuthResetPasswordJob Enqueue Error: %v", err)
		}
	}

	return response.NewResponse(
		response.WithMessage("A password reset email will be sent if the email is registered in our system."),
	).Send(ctx)
}
