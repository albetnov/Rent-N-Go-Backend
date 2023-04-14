package admin

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin/car"
)

func CarRoutes(r fiber.Router) {
	r.Get("/", car.Index)
	r.Get("/:id<int>", car.Show)
}
