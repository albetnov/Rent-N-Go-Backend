package routes

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/car"
)

func CarsRoutes(r fiber.Router) {
	r.Get("/", car.Index)
	r.Get("/recommendation", car.Recommendation)
	r.Get("/:id", car.Show)
}
