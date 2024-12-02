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

type userService struct {
	ctx *fiber.Ctx
}

func UserService(ctx *fiber.Ctx) *userService {
	return &userService{ctx: ctx}
}

func (s *userService) GetUsers() error {
	var users []entities.User

	query := database.Connection.Model(&entities.User{})

	if name := s.ctx.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if email := s.ctx.Query("email"); email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	page := utils.Paginate(&utils.PaginateParam{
		DB:    query,
		Page:  s.ctx.QueryInt("page", 1),
		Limit: s.ctx.QueryInt("limit", 10),
	}, &users)

	return response.NewResponse(
		response.WithMessage("get list users success"),
		response.WithData(page),
	).Send(s.ctx)
}

func (s *userService) CreateUser(createUserReq CreateUserRequest) error {
	hashedPassword, err := utils.HashPassword(createUserReq.Password)

	if err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to crete user"),
			response.WithError(response.ErrorInternal),
		).Send(s.ctx)
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
		).Send(s.ctx)
	}

	return response.NewResponse(
		response.WithHttpCode(http.StatusCreated),
		response.WithMessage("user created successfully"),
		response.WithData(newUser),
	).Send(s.ctx)
}

func (s *userService) GetUserByID(id string) error {
	var user entities.User

	if err := database.Connection.First(&user, "uid = ?", id).Error; err != nil {
		return response.NewResponse(
			response.WithMessage("user not found"),
			response.WithError(response.ErrorNotFound),
		).Send(s.ctx)
	}

	return response.NewResponse(
		response.WithMessage("get user success"),
		response.WithData(user),
	).Send(s.ctx)
}

func (s *userService) UpdateUser(updateUserReq UpdateUserRequest) error {
	var user entities.User

	if err := database.Connection.First(&user, "uid = ?", updateUserReq.ID).Error; err != nil {
		return response.NewResponse(
			response.WithMessage("user not found"),
			response.WithError(response.ErrorNotFound),
		).Send(s.ctx)
	}

	user.Name = updateUserReq.Name
	user.Status = updateUserReq.Status

	if user.Email != updateUserReq.Email {

		if err := database.Connection.Where("email = ? AND uid != ?", updateUserReq.Email, updateUserReq.ID).First(&entities.User{}).Error; err == nil {

			return response.NewResponse(
				response.WithMessage("invalid request"),
				response.WithData(&struct {
					Email string `json:"email"`
				}{Email: "The Email is already taken."}),
				response.WithError(response.ErrorUnprocessableEntity),
			).Send(s.ctx)
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
			).Send(s.ctx)
		}

		user.Password = hashedPassword
	}

	if err := database.Connection.Save(&user).Error; err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to update user"),
			response.WithError(response.ErrorInternal),
		).Send(s.ctx)
	}

	return response.NewResponse(
		response.WithMessage("user updated successfully"),
		response.WithData(user),
	).Send(s.ctx)
}

func (s *userService) DeleteUser(id string) error {
	var user entities.User

	if err := database.Connection.First(&user, "uid = ?", id).Error; err != nil {
		return response.NewResponse(
			response.WithMessage("user not found"),
			response.WithError(response.ErrorNotFound),
		).Send(s.ctx)
	}

	database.Connection.Delete(&user)

	return response.NewResponse(
		response.WithMessage("user deleted successfully"),
	).Send(s.ctx)
}

func (s *userService) TrashUsers() error {
	var users []entities.User

	query := database.Connection.Model(&entities.User{})

	if name := s.ctx.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if email := s.ctx.Query("email"); email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	query = query.Unscoped().Where("deleted_at IS NOT NULL")

	page := utils.Paginate(&utils.PaginateParam{
		DB:    query,
		Page:  s.ctx.QueryInt("page", 1),
		Limit: s.ctx.QueryInt("limit", 10),
	}, &users)

	return response.NewResponse(
		response.WithMessage("get list trash users success"),
		response.WithData(page),
	).Send(s.ctx)
}

func (s *userService) RestoreUser(id string) error {
	var user entities.User

	if err := database.Connection.Unscoped().First(&user, "uid = ?", id).Error; err != nil {
		return response.NewResponse(
			response.WithMessage("user not found"),
			response.WithError(response.ErrorNotFound),
		).Send(s.ctx)
	}

	if err := database.Connection.Unscoped().Model(&user).Update("deleted_at", nil).Error; err != nil {
		internal_log.Logger.Error(err.Error())
		return response.NewResponse(
			response.WithMessage("failed to restore user"),
			response.WithError(response.ErrorInternal),
		).Send(s.ctx)
	}

	return response.NewResponse(
		response.WithMessage("user restored successfully"),
	).Send(s.ctx)
}
