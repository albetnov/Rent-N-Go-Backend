package auth

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/repositories/UserRepositories"
	"rent-n-go-backend/utils"
)

func LoginView(c *fiber.Ctx) error {
	sess := utils.Session.Provide(c)

	return c.Render("auth/login", utils.Wrap(fiber.Map{}, c, sess).Message().Csrf().Validation().Get())
}

func LoginHandler(c *fiber.Ctx) error {
	payload := utils.GetPayload[LoginRequest](c)

	currentUser, err := UserRepositories.User.GetByEmail(payload.Email)

	sess := utils.Session.Provide(c)

	if err != nil || !utils.ComparePassword(payload.Password, currentUser.Password) || currentUser.Role != "admin" {
		sess.SetSession("message", "Wrong credential being passed")
		return c.RedirectBack("/auth/login")
	}

	sess.SetSession("authed", currentUser.ID)
	return c.Redirect("/admin/dashboard")
}

func Logout(c *fiber.Ctx) error {
	utils.Session.Provide(c).DeleteSession("authed")
	return c.Redirect("/auth/login")
}
