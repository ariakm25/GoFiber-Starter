package app

import (
	"GoFiber-API/app/auth"
	"GoFiber-API/app/monitoring"
	"GoFiber-API/app/user"

	"github.com/gofiber/fiber/v2"
)

func MainModule(router fiber.Router) {
	monitoring.Module(router)
	auth.Module(router)
	user.Module(router)
}
