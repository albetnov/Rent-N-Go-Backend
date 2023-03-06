package auth

import "github.com/gofiber/fiber/v2"

func LoginView(c *fiber.Ctx) error {
	return c.Render("auth/login", nil, "auth/layout")
}
