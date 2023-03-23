package routes

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/order"
)

func OrderRoutes(r fiber.Router) {
	r.Get("/", order.History)
}
