package auth

import (
	"GoFiber-API/entities"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/external/database/redis"
	"GoFiber-API/infra/response"
	"GoFiber-API/internal/config"
	internal_log "GoFiber-API/internal/log"
	"GoFiber-API/internal/queue"
	"GoFiber-API/internal/utils"
	"context"
	"errors"
	"strings"
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

	hashedPassword, err := utils.HashPassword(registerReq.Password)

	if err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to register user"),
			response.WithError(response.ErrorInternal),
		).Send(ctx)
	}

	user := &entities.User{
		Name:     registerReq.Name,
		Email:    registerReq.Email,
		Password: hashedPassword,
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

	var user = &entities.User{}

	findUser := database.Connection.First(user, "email = ?", loginReq.Email)

	if findUser.Error != nil {
		if errors.Is(findUser.Error, gorm.ErrRecordNotFound) {
			return response.NewResponse(
				response.WithMessage("invalid email or password"),
				response.WithData(fiber.Map{
					"email": "invalid email or password",
				}),
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
			response.WithData(fiber.Map{
				"email": "invalid email or password",
			}),
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

	refresh_token := utils.GenerateRefreshToken(user.UID)

	user_agent := ctx.Get("User-Agent")

	if err := database.Connection.Create(&entities.UserSession{
		UserID:       user.UID,
		RefreshToken: refresh_token,
		DeviceInfo:   user_agent,
		ExpiredAt:    time.Now().Add(time.Hour * time.Duration(24*config.GetConfig.REFRESH_TOKEN_EXPIRATION_DAYS)),
	}).Error; err != nil {
		internal_log.Logger.Error(err.Error())

		return response.NewResponse(
			response.WithMessage("failed to login"),
			response.WithError(response.ErrorInternal),
		).Send(ctx)
	}

	return response.NewResponse(
		response.WithMessage("success login"),
		response.WithData(fiber.Map{
			"token":         token,
			"refresh_token": refresh_token,
		}),
	).Send(ctx)
}

func Me(ctx *fiber.Ctx) error {
	payload := ctx.Locals(pasetoware.DefaultContextKey).(*entities.User)

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

	var user entities.User

	err := database.Connection.Where(&entities.User{
		Email: resetPasswordReq.Email,
	}).First(&user).Error

	if err == nil {
		var checkResetPasswordToken entities.UserToken

		database.Connection.Order("created_at desc").Where(&entities.UserToken{
			UserID: user.UID,
			Type:   entities.UserTokenTypeResetPassword,
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

	database.Connection.Delete(&entities.UserToken{}, "user_id = ? AND type = ?", user.UID, entities.UserTokenTypeResetPassword)

	return response.NewResponse(
		response.WithMessage("A password reset email will be sent if the email is registered in our system."),
	).Send(ctx)
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

	var user entities.User

	err := database.Connection.Where(&entities.User{
		Email: validateResetPasswordTokenReq.Email,
	}).First(&user).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or token"),
			response.WithError(response.ErrorNotFound),
		).Send(ctx)
	}

	var userToken entities.UserToken

	err = database.Connection.Where(&entities.UserToken{
		UserID: user.UID,
		Token:  validateResetPasswordTokenReq.Token,
		Type:   entities.UserTokenTypeResetPassword,
	}).First(&userToken).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or token"),
			response.WithError(response.ErrorNotFound),
		).Send(ctx)
	}

	if userToken.ExpiredAt.Before(time.Now()) {
		return response.NewResponse(
			response.WithMessage("token has been expired"),
			response.WithError(response.ErrorBadRequest),
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

	var user entities.User

	err := database.Connection.Where(&entities.User{
		Email: resetPasswordReq.Email,
	}).First(&user).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or token"),
			response.WithError(response.ErrorNotFound),
		).Send(ctx)
	}

	var userToken entities.UserToken

	err = database.Connection.Where(&entities.UserToken{
		UserID: user.UID,
		Token:  resetPasswordReq.Token,
		Type:   entities.UserTokenTypeResetPassword,
	}).First(&userToken).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or token"),
			response.WithError(response.ErrorNotFound),
		).Send(ctx)
	}

	if userToken.ExpiredAt.Before(time.Now()) {
		return response.NewResponse(
			response.WithMessage("token has been expired"),
			response.WithError(response.ErrorBadRequest),
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

func Logout(ctx *fiber.Ctx) error {

	token := strings.Split(ctx.Get(fiber.HeaderAuthorization), "Bearer ")[1]

	if token != "" {
		data, err := utils.DecryptPaseto(token)

		if err != nil {
			internal_log.Logger.Error(err.Error())
		}

		redis.RedisStore.Conn().Set(context.Background(), "blacklist_token:"+token, "true", time.Until(data.Expiration))
	}

	logoutReq := &LogoutRequest{}

	if err := ctx.BodyParser(logoutReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorBadRequest),
			response.WithMessage("invalid request"),
		).Send(ctx)
	}

	if logoutReq.RefreshToken != "" {
		database.Connection.Delete(&entities.UserSession{}, "refresh_token = ?", logoutReq.RefreshToken)
	}

	return response.NewResponse(
		response.WithMessage("success logout"),
	).Send(ctx)
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

	var userSession entities.UserSession

	err := database.Connection.Where(&entities.UserSession{
		RefreshToken: refreshTokenReq.RefreshToken,
	}).First(&userSession).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid refresh token"),
			response.WithError(response.ErrorUnauthorized),
		).Send(ctx)
	}

	if userSession.ExpiredAt.Before(time.Now()) {
		return response.NewResponse(
			response.WithMessage("refresh token has been expired"),
			response.WithError(response.ErrorUnauthorized),
		).Send(ctx)
	}

	var user entities.User

	err = database.Connection.Where(&entities.User{
		UID: userSession.UserID,
	}).First(&user).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("user not found"),
			response.WithError(response.ErrorNotFound),
		).Send(ctx)
	}

	newToken, err := utils.GenerateLocalPaseto(user.UID)

	if err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed generate token"),
			response.WithError(response.ErrorInternal),
		).Send(ctx)
	}

	newRefreshToken := utils.GenerateRefreshToken(user.UID)

	user_agent := ctx.Get("User-Agent")

	userSession.RefreshToken = newRefreshToken
	userSession.DeviceInfo = user_agent

	if err := database.Connection.Save(&userSession).Error; err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to refresh token"),
			response.WithError(response.ErrorInternal),
		).Send(ctx)
	}

	token := strings.Split(ctx.Get(fiber.HeaderAuthorization), "Bearer ")[1]

	if token != "" {
		data, err := utils.DecryptPaseto(token)

		if err != nil {
			internal_log.Logger.Error(err.Error())
		}

		redis.RedisStore.Conn().Set(context.Background(), "blacklist_token:"+token, "true", time.Until(data.Expiration))
	}

	return response.NewResponse(
		response.WithMessage("success refresh token"),
		response.WithData(fiber.Map{
			"token":         newToken,
			"refresh_token": newRefreshToken,
		}),
	).Send(ctx)
}
