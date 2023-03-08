package admin

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin/user"
	"rent-n-go-backend/utils"
)

func UsersModuleRoutes(r fiber.Router) {
	r.Get("/", user.Index)
	r.Get("/create", user.Create)
	r.Post("/create", utils.InterceptWebRequest(new(user.CreateUserPayload)), user.Store)
	r.Get("/edit/:id<int>", user.Edit)
	r.Post("/edit/:id<int>", utils.InterceptWebRequest(new(user.UpdateUserPayload)), user.Update)
	r.Get("/:id<int>", user.Show)
}
