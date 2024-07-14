package middleware

import (
	user_entities "GoFiber-API/app/user/entities"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/infra/response"
	"GoFiber-API/internal/config"
	internal_log "GoFiber-API/internal/log"
	"encoding/json"
	"errors"
	"time"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/contrib/casbin"
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
				internal_log.Logger.Sugar().Errorf("Error unmarshal data: %s", err.Error())
				return nil, paseto.ErrDataUnmarshal
			}

			if err := payload.Validate(
				paseto.ValidAt(time.Now()), paseto.Subject(pasetoTokenSubject),
				paseto.ForAudience(pasetoTokenAudience),
			); err != nil {
				internal_log.Logger.Sugar().Errorf("Error validate token: %s", err.Error())
				return "", err
			}

			var user = &user_entities.User{}

			findUser := database.Connection.First(user, "uid = ?", payload.Get(pasetoTokenField))

			if findUser.Error != nil {
				return nil, pasetoware.ErrExpiredToken
			}

			return user, nil
		},
		TokenPrefix: "Bearer",
	})

}

var Rbac *casbin.Middleware

func InitRbac(pathModel string, adapter *gormadapter.Adapter) {
	Rbac = casbin.New(casbin.Config{
		ModelFilePath: pathModel,
		PolicyAdapter: adapter,
		Forbidden: func(c *fiber.Ctx) error {
			return response.NewResponse(
				response.WithMessage("You don't have permission to access this resource"),
				response.WithError(response.ErrorForbiddenAccess),
			).Send(c)
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return response.NewResponse(
				response.WithMessage("You don't have permission to access this resource"),
				response.WithError(response.ErrorUnauthorized),
			).Send(c)
		},
		Lookup: func(c *fiber.Ctx) string {
			checkLocal := c.Locals(pasetoware.DefaultContextKey)

			if checkLocal == nil {
				return ""
			}

			user := checkLocal.(*user_entities.User)

			return user.UID
		},
	})
}
