package auth

import (
	"GoFiber-API/infra/middleware"

	"github.com/gofiber/fiber/v2"
)

func Module(router fiber.Router) {
	route := router.Group("/auth")

	route.Post("/login", Login)
	route.Post("/register", Register)

	route.Post("/forgot-password", ForgotPassword)
	route.Post("/validate-reset-password-token", ValidateResetPasswordToken)
	route.Post("/reset-password", ResetPassword)

	route.Post("/refresh-token", RefreshToken)

	route.Post("/logout", Logout)

	route.Get("/me", middleware.AuthMiddleware(), Me)

}
