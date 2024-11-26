package user

import (
	"GoFiber-API/entities"
	database "GoFiber-API/external/database/postgres"
	"GoFiber-API/infra/response"
	internal_log "GoFiber-API/internal/log"
	"GoFiber-API/internal/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	var users []entities.User

	query := database.Connection.Model(&entities.User{})

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
		response.WithMessage("get list users success"),
		response.WithData(page),
	).Send(c)
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

	hashedPassword, err := utils.HashPassword(createUserReq.Password)

	if err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to crete user"),
			response.WithError(response.ErrorInternal),
		).Send(c)
	}

	newUser := &entities.User{
		Name:     createUserReq.Name,
		Email:    createUserReq.Email,
		Status:   createUserReq.Status,
		Password: hashedPassword,
	}

	if err := database.Connection.Create(newUser).Error; err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed register user"),
			response.WithError(response.ErrorInternal),
		).Send(c)
	}

	return response.NewResponse(
		response.WithHttpCode(http.StatusCreated),
		response.WithMessage("user created successfully"),
		response.WithData(newUser),
	).Send(c)
}

func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	var user entities.User

	if err := database.Connection.First(&user, "uid = ?", id).Error; err != nil {
		return response.NewResponse(
			response.WithMessage("user not found"),
			response.WithError(response.ErrorNotFound),
		).Send(c)
	}

	return response.NewResponse(
		response.WithMessage("get user success"),
		response.WithData(user),
	).Send(c)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user entities.User

	if err := database.Connection.First(&user, "uid = ?", id).Error; err != nil {
		return response.NewResponse(
			response.WithMessage("user not found"),
			response.WithError(response.ErrorNotFound),
		).Send(c)
	}

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

	user.Name = updateUserReq.Name
	user.Status = updateUserReq.Status

	if user.Email != updateUserReq.Email {

		if err := database.Connection.Where("email = ? AND uid != ?", updateUserReq.Email, id).First(&entities.User{}).Error; err == nil {

			return response.NewResponse(
				response.WithMessage("invalid request"),
				response.WithData(&struct {
					Email string `json:"email"`
				}{Email: "The Email is already taken."}),
				response.WithError(response.ErrorUnprocessableEntity),
			).Send(c)
		}

		user.Email = updateUserReq.Email

	}

	if updateUserReq.Password != "" {
		hashedPassword, err := utils.HashPassword(updateUserReq.Password)

		if err != nil {
			internal_log.Logger.Error(err.Error())
			return response.NewResponse(
				response.WithMessage("failed to update user"),
				response.WithError(response.ErrorInternal),
			).Send(c)
		}

		user.Password = hashedPassword
	}

	if err := database.Connection.Save(&user).Error; err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to update user"),
			response.WithError(response.ErrorInternal),
		).Send(c)
	}

	return response.NewResponse(
		response.WithMessage("user updated successfully"),
		response.WithData(user),
	).Send(c)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user entities.User

	if err := database.Connection.First(&user, "uid = ?", id).Error; err != nil {
		return response.NewResponse(
			response.WithMessage("user not found"),
			response.WithError(response.ErrorNotFound),
		).Send(c)
	}

	database.Connection.Delete(&user)

	return response.NewResponse(
		response.WithMessage("user deleted successfully"),
	).Send(c)
}

func TrashUsers(c *fiber.Ctx) error {
	var users []entities.User

	query := database.Connection.Model(&entities.User{})

	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if email := c.Query("email"); email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	query = query.Unscoped().Where("deleted_at IS NOT NULL")

	page := utils.Paginate(&utils.PaginateParam{
		DB:      query,
		Page:    c.QueryInt("page", 1),
		Limit:   c.QueryInt("limit", 10),
		ShowSQL: true,
	}, &users)

	return response.NewResponse(
		response.WithMessage("get list trash users success"),
		response.WithData(page),
	).Send(c)
}

func RestoreUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user entities.User

	if err := database.Connection.Unscoped().First(&user, "uid = ?", id).Error; err != nil {
		return response.NewResponse(
			response.WithMessage("user not found"),
			response.WithError(response.ErrorNotFound),
		).Send(c)
	}

	if err := database.Connection.Unscoped().Model(&user).Update("deleted_at", nil).Error; err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to restore user"),
			response.WithError(response.ErrorInternal),
		).Send(c)
	}

	return response.NewResponse(
		response.WithMessage("user restored successfully"),
	).Send(c)
}
