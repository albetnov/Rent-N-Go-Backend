package admin

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin/car"
	"rent-n-go-backend/utils"
)

func CarRoutes(r fiber.Router) {
	r.Get("/", car.Index)
	r.Get("/create", car.Create)
	r.Post("/create", utils.InterceptWebRequest(new(car.CarPayload)), car.Store)
	r.Get("/edit/:id<int>", car.Edit)
	r.Post("/edit/:id<int>", utils.InterceptWebRequest(new(car.EditCarPayload)), car.Update)
	r.Get("/delete/:id<int>", car.Delete)
	r.Get("/:id<int>", car.Show)
}
