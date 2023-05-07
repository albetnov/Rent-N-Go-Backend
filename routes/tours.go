package routes

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/tour"
)

func TourRoutes(r fiber.Router) {
	// Define routes for tour module
	r.Get("/:id", tour.Show)
	r.Get("/:id/stock", tour.Stocks)
	r.Get("/", tour.Index)
}
