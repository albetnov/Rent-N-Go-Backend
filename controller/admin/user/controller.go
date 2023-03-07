package user

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/query"
	"rent-n-go-backend/utils"
)

func Dashboard(c *fiber.Ctx) error {
	return utils.RenderTemplate(c, "dashboard", "Dashboard", nil)
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

	return utils.RenderTemplate(c, "users/index", "Manage Users", fiber.Map{
		"Users": lists,
	})
}
