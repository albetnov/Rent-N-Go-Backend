package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin"
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/UserRepositories"
	"rent-n-go-backend/utils"
	"strconv"
)

func Dashboard(c *fiber.Ctx) error {
	return admin.RenderTemplate(c, "dashboard", "Dashboard", nil)
}

func Index(c *fiber.Ctx) error {
	users, err := query.User.Find()

	lists := make([]UserModels.User, len(users))

	for i, v := range users {
		lists[i] = *v
	}

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	return admin.RenderTemplate(c, "users/index", "Manage Users", fiber.Map{
		"Users": lists,
		"Error": utils.Session.Provide(c).GetFlash("error"),
	})
}

func Show(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/users")
	}

	id := uint(userId)

	user, err := UserRepositories.User.GetAllById(id)

	if err != nil {
		sess.SetSession("error", "Ups user not found")
		return c.RedirectBack("/admin/users")
	}

	return admin.RenderTemplate(c, "users/show", fmt.Sprintf("%s Detail", user.Name), fiber.Map{
		"Name":        user.Name,
		"Email":       user.Email,
		"PhoneNumber": user.PhoneNumber,
		"Role":        user.Role,
		"SIM":         user.Sim,
		"NIK":         user.Nik,
		"Photo":       user.Photo,
	})
}
