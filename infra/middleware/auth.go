package middleware

import (
	"GoFiber-API/app/user"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/infra/response"
	"GoFiber-API/internal/config"
	"encoding/json"
	"errors"
	"time"

	pasetoware "github.com/gofiber/contrib/paseto"
	"github.com/gofiber/fiber/v2"
	"github.com/o1egl/paseto"
)

func AuthMiddleware() func(*fiber.Ctx) error {
	return pasetoware.New(pasetoware.Config{
		SymmetricKey: []byte(config.GetConfig.PASETO_LOCAL_SECRET_SYMMETRIC_KEY),

		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			errorStatus := response.ErrorBadRequest
			errorMessage := err.Error()
			if errors.Is(err, pasetoware.ErrDataUnmarshal) || errors.Is(err, pasetoware.ErrExpiredToken) || errors.Is(err, pasetoware.ErrMissingToken) || errors.Is(err, pasetoware.ErrIncorrectTokenPrefix) {
				errorStatus = response.ErrorUnauthorized
				errorMessage = "invalid token authentication"
			}

			return response.NewResponse(
				response.WithMessage(errorMessage),
				response.WithError(errorStatus),
			).Send(ctx)

		},
		Validate: func(data []byte) (interface{}, error) {
			const (
				pasetoTokenAudience = "gofiber.gophers"
				pasetoTokenSubject  = "user-token"
				pasetoTokenField    = "data"
			)

			var payload paseto.JSONToken
			if err := json.Unmarshal(data, &payload); err != nil {
				return nil, paseto.ErrDataUnmarshal
			}

			if time.Now().After(payload.Expiration) {
				return nil, pasetoware.ErrExpiredToken
			}
			if err := payload.Validate(
				paseto.ValidAt(time.Now()), paseto.Subject(pasetoTokenSubject),
				paseto.ForAudience(pasetoTokenAudience),
			); err != nil {
				return "", err
			}

			var user = &user.User{}

			findUser := database.Connection.First(user, "uid = ?", payload.Get(pasetoTokenField))

			if findUser.Error != nil {
				return nil, pasetoware.ErrExpiredToken
			}

			return user, nil
		},
		TokenPrefix: "Bearer",
	})

}
