package profile

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin"
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/repositories/UserRepositories"
	"rent-n-go-backend/utils"
)

func Index(c *fiber.Ctx) error {
	sess := utils.Session.Provide(c)
	res := utils.Wrap(fiber.Map{}, nil, sess).Message().Error().Validation().Get()
	return admin.RenderTemplate(c, "profile", "Profile", res)
}

func UpdatePicture(c *fiber.Ctx) error {
	sess := utils.Session.Provide(c)
	userId := sess.GetSession("authed").(uint)

	err := UserRepositories.User.CreatePhoto(c, sess, "photo", userId, "/admin/profile")

	if err != nil {
		return err
	}

	sess.SetSession("message", "Profile Updated Successfully")
	return c.RedirectBack("/admin/profile")
}

func UpdateProfile(c *fiber.Ctx) error {
	payload := utils.GetPayload[UpdateProfilePayload](c)

	userData := UserModels.User{
		Name: payload.Name,
	}

	if payload.Password != "" {
		hashedPass, _ := utils.HashPassword(payload.Password)
		userData.Password = hashedPass
	}

	sess := utils.Session.Provide(c)
	userId := sess.GetSession("authed").(uint)

	err := UserRepositories.User.UpdateById(c, userId, &userData)

	if err != nil {
		return err
	}

	sess.SetSession("message", "Profile Info Updated Successfully")
	return c.RedirectBack("/admin/profile")
}
