package auth

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/utils"
)

func Guest(c *fiber.Ctx) error {
	session := utils.Session.Provide(c)

	if session.GetSession("authed") != nil {
		return c.Redirect("/dashboard")
	}

	return c.Next()
}

func Authed(c *fiber.Ctx) error {
	if utils.Session.Provide(c).GetSession("authed") != nil {
		return c.Next()
	}

	return c.Redirect("/login")
}
