package routes

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/driver"
)

func DriverRoutes(r fiber.Router) {
	// Define routes for tour module
	r.Get("/:id", driver.Show)
	r.Get("/", driver.Index)
}
