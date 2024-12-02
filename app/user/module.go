package user

import (
	"GoFiber-API/infra/middleware"

	"github.com/gofiber/fiber/v2"
)

func Module(router fiber.Router) {
	route := router.Group("/users")

	route.Post("/", middleware.AuthMiddleware(), middleware.Rbac.RequiresPermissions([]string{"user:create"}), CreateUser)
	route.Get("/", middleware.AuthMiddleware(), middleware.Rbac.RequiresPermissions([]string{"user:read"}), GetUsers)
	route.Put("/", middleware.AuthMiddleware(), middleware.Rbac.RequiresPermissions([]string{"user:update"}), UpdateUser)
	route.Get("/:id", middleware.AuthMiddleware(), middleware.Rbac.RequiresPermissions([]string{"user:read"}), GetUserByID)
	route.Delete("/:id", middleware.AuthMiddleware(), middleware.Rbac.RequiresPermissions([]string{"user:delete"}), DeleteUser)

	route.Get("/trash/list", middleware.AuthMiddleware(), middleware.Rbac.RequiresPermissions([]string{"user:restore"}), TrashUsers)
	route.Put("/restore/:id", middleware.AuthMiddleware(), middleware.Rbac.RequiresPermissions([]string{"user:restore"}), RestoreUser)
}
