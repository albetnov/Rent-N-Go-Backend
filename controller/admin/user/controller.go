package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/UserRepositories"
	"rent-n-go-backend/utils"
)

func Dashboard(c *fiber.Ctx) error {
	userId := utils.Session.Provide(c).GetSession("authed").(uint)
	user, _ := UserRepositories.User.GetById(userId)

	if user == nil {
		return c.Render("dashboard", fiber.Map{
			"Message": "Failed to fetch current user, logout advised.",
		}, "layout")
	}

	userPhoto, err := query.User.Photo.Model(user).Find()

	propic := "https://source.unsplash.com/500x500?potrait"

	if err == nil {
		propic = fmt.Sprintf("/public/files/user/%s", userPhoto.PhotoPath)
	}

	return c.Render("dashboard", fiber.Map{
		"Name":       user.Name,
		"ModuleName": "Dashboard",
		"Propic":     propic,
	}, "layout")
}

func Index(c *fiber.Ctx) error {
	users, err := query.User.Find()

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	return c.Render("users/index", fiber.Map{
		"Users": users,
	})
}
