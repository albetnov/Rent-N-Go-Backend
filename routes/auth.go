package routes

import "github.com/gofiber/fiber/v2"

func AuthRoutes(r fiber.Router) {
	r.Get("/login", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "Welcome to Rent-N-Go Backend!",
			"status":  200,
		})
	})
}
