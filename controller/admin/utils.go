package admin

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/UserRepositories"
	"rent-n-go-backend/utils"
)

// RenderTemplate
// Render admin template with specified layout and rules and default value.
func RenderTemplate(c *fiber.Ctx, name string, moduleName string, data fiber.Map) error {
	userId := utils.Session.Provide(c).GetSession("authed").(uint)
	user, _ := UserRepositories.User.GetById(userId)

	if user == nil {
		return c.Render("dashboard", fiber.Map{
			"Message": "Failed to fetch current user, logout advised.",
		}, "layout")
	}

	userPhoto, err := query.User.Photo.Model(user).Find()

	propic := "https://source.unsplash.com/500x500?potrait"

	if err == nil && userPhoto.ID != 0 {
		propic = fmt.Sprintf("/public/files/user/%s", userPhoto.PhotoPath)
	}

	if data == nil {
		data = fiber.Map{}
	}

	data["_Name"] = user.Name
	data["_Propic"] = propic
	data["_ModuleName"] = moduleName
	data["_csrf"] = c.Locals("token")

	return c.Render(name, data, "layout")
}
