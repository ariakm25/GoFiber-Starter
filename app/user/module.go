package user

import (
	"GoFiber-API/infra/middleware"

	"github.com/gofiber/fiber/v2"
)

func Module(router fiber.Router) {
	route := router.Group("/users")

	route.Get("/", middleware.AuthMiddleware(), middleware.Rbac.RequiresPermissions([]string{"user:read"}), GetUsers)
}
