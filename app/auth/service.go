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

type authService struct {
	ctx *fiber.Ctx
}

func AuthService(ctx *fiber.Ctx) *authService {
	return &authService{ctx: ctx}
}

func (s *authService) Register(registerReq RegisterRequest) error {
	hashedPassword, err := utils.HashPassword(registerReq.Password)

	if err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to register user"),
			response.WithError(response.ErrorInternal),
		).Send(s.ctx)
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
		).Send(s.ctx)
	}

	return response.NewResponse(
		response.WithMessage("success register"),
		response.WithData(user),
	).Send(s.ctx)
}

func (s *authService) Login(loginReq LoginRequest) error {
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
			).Send(s.ctx)
		}

		return response.NewResponse(
			response.WithMessage(findUser.Error.Error()),
			response.WithError(response.ErrorUnauthorized),
		).Send(s.ctx)
	}

	if err := utils.ComparePassword(user.Password, loginReq.Password); err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or password"),
			response.WithData(fiber.Map{
				"email": "invalid email or password",
			}),
			response.WithError(response.ErrorUnauthorized),
		).Send(s.ctx)
	}

	token, err := utils.GenerateLocalPaseto(user.UID)

	if err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed generate token"),
			response.WithError(response.ErrorInternal),
		).Send(s.ctx)
	}

	refresh_token := utils.GenerateRefreshToken(user.UID)

	user_agent := s.ctx.Get("User-Agent")

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
		).Send(s.ctx)
	}

	return response.NewResponse(
		response.WithMessage("success login"),
		response.WithData(fiber.Map{
			"token":         token,
			"refresh_token": refresh_token,
		}),
	).Send(s.ctx)
}

func (s *authService) RefreshToken(refreshTokenReq RefreshTokenRequest) error {
	var userSession entities.UserSession

	err := database.Connection.Where(&entities.UserSession{
		RefreshToken: refreshTokenReq.RefreshToken,
	}).First(&userSession).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid refresh token"),
			response.WithError(response.ErrorUnauthorized),
		).Send(s.ctx)
	}

	if userSession.ExpiredAt.Before(time.Now()) {
		return response.NewResponse(
			response.WithMessage("refresh token has been expired"),
			response.WithError(response.ErrorUnauthorized),
		).Send(s.ctx)
	}

	var user entities.User

	err = database.Connection.Where(&entities.User{
		UID: userSession.UserID,
	}).First(&user).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("user not found"),
			response.WithError(response.ErrorNotFound),
		).Send(s.ctx)
	}

	newToken, err := utils.GenerateLocalPaseto(user.UID)

	if err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed generate token"),
			response.WithError(response.ErrorInternal),
		).Send(s.ctx)
	}

	newRefreshToken := utils.GenerateRefreshToken(user.UID)

	user_agent := s.ctx.Get("User-Agent")

	userSession.RefreshToken = newRefreshToken
	userSession.DeviceInfo = user_agent

	if err := database.Connection.Save(&userSession).Error; err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to refresh token"),
			response.WithError(response.ErrorInternal),
		).Send(s.ctx)
	}

	token := strings.Split(s.ctx.Get(fiber.HeaderAuthorization), "Bearer ")[1]

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
	).Send(s.ctx)
}

func (s *authService) ForgotPassword(forgotPasswordReq ForgotPasswordRequest) error {
	var user entities.User

	err := database.Connection.Where(&entities.User{
		Email: forgotPasswordReq.Email,
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
			).Send(s.ctx)
		}
	}

	task, err := NewAuthResetPasswordJob(forgotPasswordReq.Email, user.UID)
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
	).Send(s.ctx)
}

func (s *authService) ValidateResetPasswordToken(validateResetPasswordTokenReq ValidateResetPasswordTokenRequest) error {
	var user entities.User

	err := database.Connection.Where(&entities.User{
		Email: validateResetPasswordTokenReq.Email,
	}).First(&user).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or token"),
			response.WithError(response.ErrorNotFound),
		).Send(s.ctx)
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
		).Send(s.ctx)
	}

	if userToken.ExpiredAt.Before(time.Now()) {
		return response.NewResponse(
			response.WithMessage("token has been expired"),
			response.WithError(response.ErrorBadRequest),
		).Send(s.ctx)

	}

	return response.NewResponse(
		response.WithMessage("success validate token"),
	).Send(s.ctx)
}

func (s *authService) ResetPassword(resetPasswordReq ResetPasswordRequest) error {
	var user entities.User

	err := database.Connection.Where(&entities.User{
		Email: resetPasswordReq.Email,
	}).First(&user).Error

	if err != nil {
		return response.NewResponse(
			response.WithMessage("invalid email or token"),
			response.WithError(response.ErrorNotFound),
		).Send(s.ctx)
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
		).Send(s.ctx)
	}

	if userToken.ExpiredAt.Before(time.Now()) {
		return response.NewResponse(
			response.WithMessage("token has been expired"),
			response.WithError(response.ErrorBadRequest),
		).Send(s.ctx)
	}

	hashedPassword, err := utils.HashPassword(resetPasswordReq.Password)

	if err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to reset password"),
			response.WithError(response.ErrorInternal),
		).Send(s.ctx)
	}

	user.Password = hashedPassword

	if err := database.Connection.Save(&user).Error; err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to reset password"),
			response.WithError(response.ErrorInternal),
		).Send(s.ctx)
	}

	return response.NewResponse(
		response.WithMessage("success reset password"),
	).Send(s.ctx)
}

func (s *authService) Logout(logoutReq LogoutRequest) error {
	token := strings.Split(s.ctx.Get(fiber.HeaderAuthorization), "Bearer ")[1]

	if token != "" {
		data, err := utils.DecryptPaseto(token)

		if err != nil {
			internal_log.Logger.Error(err.Error())
		}

		redis.RedisStore.Conn().Set(context.Background(), "blacklist_token:"+token, "true", time.Until(data.Expiration))
	}

	if logoutReq.RefreshToken != "" {
		database.Connection.Delete(&entities.UserSession{}, "refresh_token = ?", logoutReq.RefreshToken)
	}

	return response.NewResponse(
		response.WithMessage("success logout"),
	).Send(s.ctx)
}

func (s *authService) Me() error {
	payload := s.ctx.Locals(pasetoware.DefaultContextKey).(*entities.User)

	return response.NewResponse(
		response.WithMessage("success get user data"),
		response.WithData(payload),
	).Send(s.ctx)
}
