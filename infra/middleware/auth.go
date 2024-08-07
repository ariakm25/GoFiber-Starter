package middleware

import (
	"GoFiber-API/entities"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/external/database/redis"
	"GoFiber-API/infra/response"
	internal_casbin "GoFiber-API/internal/casbin"
	"GoFiber-API/internal/config"
	internal_log "GoFiber-API/internal/log"
	"context"
	"encoding/json"
	"errors"
	"strings"
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
			errorStatus := response.ErrorUnauthorized
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
				return "", err
			}

			var user = &entities.User{}

			findUser := database.Connection.First(user, "uid = ?", payload.Get(pasetoTokenField))

			if findUser.Error != nil {
				return nil, pasetoware.ErrExpiredToken
			}

			if user.Status != "active" {
				return nil, errors.New("user is " + user.Status)
			}

			return user, nil
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			originalToken := strings.Split(c.Get("Authorization"), "Bearer ")[1]

			checkBlacklist, _ := redis.RedisStore.Conn().Get(context.Background(), "blacklist_token:"+originalToken).Result()

			if checkBlacklist == "true" {
				return response.NewResponse(
					response.WithMessage(pasetoware.ErrExpiredToken.Error()),
					response.WithError(response.ErrorUnauthorized),
				).Send(c)
			}

			return c.Next()
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

			user := checkLocal.(*entities.User)

			role, err := internal_casbin.CasbinEnforcer.GetRolesForUser(user.UID)

			if err != nil {
				return ""
			}

			return role[0]

		},
	})
}
