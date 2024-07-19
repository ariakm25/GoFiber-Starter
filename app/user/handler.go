package user

import (
	user_entities "GoFiber-API/app/user/entities"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/infra/response"
	"GoFiber-API/internal/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	var users []user_entities.User

	query := database.Connection.Model(&user_entities.User{})

	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if email := c.Query("email"); email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	page := utils.Paginate(&utils.PaginateParam{
		DB:      query,
		Page:    c.QueryInt("page", 1),
		Limit:   c.QueryInt("limit", 10),
		ShowSQL: true,
	}, &users)

	return response.NewResponse(
		response.WithMessage("get list products success"),
		response.WithData(page),
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
