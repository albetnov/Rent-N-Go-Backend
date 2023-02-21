package home

import (
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Name": "Asep Surasep",
	})
}
