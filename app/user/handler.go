package user

import (
	"GoFiber-API/infra/response"
	"GoFiber-API/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	return UserService(c).GetUsers()
}

func CreateUser(c *fiber.Ctx) error {
	createUserReq := &CreateUserRequest{}

	if err := c.BodyParser(createUserReq); err != nil {
		return response.NewResponse(
			response.WithMessage("invalid request"),
			response.WithError(response.ErrorBadRequest),
		).Send(c)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(createUserReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(c)
	}

	return UserService(c).CreateUser(*createUserReq)
}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	return UserService(c).GetUserByID(id)
}

func UpdateUser(c *fiber.Ctx) error {
	updateUserReq := &UpdateUserRequest{}

	if err := c.BodyParser(updateUserReq); err != nil {
		return response.NewResponse(
			response.WithMessage("invalid request"),
			response.WithError(response.ErrorBadRequest),
		).Send(c)
	}

	validate := utils.NewValidator()

	if err := validate.Struct(updateUserReq); err != nil {
		return response.NewResponse(
			response.WithMessage(err.Error()),
			response.WithError(response.ErrorUnprocessableEntity),
			response.WithData(utils.ValidatorErrors(err)),
			response.WithMessage("invalid request"),
		).Send(c)
	}

	return UserService(c).UpdateUser(*updateUserReq)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	return UserService(c).DeleteUser(id)
}

func TrashUsers(c *fiber.Ctx) error {
	return UserService(c).TrashUsers()
}

func RestoreUser(c *fiber.Ctx) error {
	id := c.Params("id")

	return UserService(c).RestoreUser(id)
}
