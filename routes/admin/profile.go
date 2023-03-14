package admin

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin/profile"
	"rent-n-go-backend/utils"
)

func ProfileRoutes(r fiber.Router) {
	r.Get("/", profile.Index)
	r.Post("/update", utils.InterceptWebRequest(&profile.UpdateProfilePayload{}), profile.UpdateProfile)
	r.Post("/update/picture", profile.UpdatePicture)
}
