package user

import (
	user_entities "GoFiber-API/app/user/entities"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/infra/response"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	var users []user_entities.User
	database.Connection.Limit(100).Find(&users)
	return response.NewResponse(
		response.WithMessage("get list products success"),
		response.WithData(users),
	).Send(c)
}

func CreateUser(c *fiber.Ctx) error {
	user := new(user_entities.User)
	if err := c.BodyParser(user); err != nil {
		return response.NewResponse(
			response.WithHttpCode(http.StatusBadRequest),
			response.WithMessage("invalid request"),
			response.WithError(err),
		).Send(c)
	}

	database.Connection.Create(&user)

	return response.NewResponse(
		response.WithHttpCode(http.StatusCreated),
		response.WithMessage("create user success"),
		response.WithData(user),
	).Send(c)
}
