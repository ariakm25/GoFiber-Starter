package user

import "github.com/gofiber/fiber/v2"

func Module(router fiber.Router) {
	route := router.Group("/users")

	route.Get("/", GetUsers)
}
