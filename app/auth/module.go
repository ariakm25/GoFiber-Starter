package auth

import (
	"GoFiber-API/infra/middleware"

	"github.com/gofiber/fiber/v2"
)

func Module(router fiber.Router) {
	route := router.Group("/auth")

	route.Post("/login", Login)
	route.Post("/register", Register)

	route.Post("/reset-password", ResetPassword)

	route.Get("/me", middleware.AuthMiddleware(), Me)

}
