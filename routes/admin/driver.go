package admin

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin/driver"
	"rent-n-go-backend/utils"
)

func DriverRoutes(r fiber.Router) {
	r.Get("/", driver.Index)
	r.Get("/create", driver.Create)
	r.Post("/create", utils.InterceptWebRequest(new(driver.DriverPayload)), driver.Store)
	r.Get("/edit/:id<int>", driver.Edit)
	r.Post("/edit/:id<int>", utils.InterceptWebRequest(new(driver.DriverPayload)), driver.Update)
	r.Get("/delete/:id<int>", driver.Delete)
	r.Get("/:id<int>", driver.Show)
}
