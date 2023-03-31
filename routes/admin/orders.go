package admin

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin/order"
)

func OrderRoutes(r fiber.Router) {
	r.Get("/", order.Index)
	r.Get("/:id<int>", order.Show)
	r.Post("/update/status/:id<int>", order.UpdateStatus)
	r.Post("/delete/:id<int>", order.DeleteOrder)
}
