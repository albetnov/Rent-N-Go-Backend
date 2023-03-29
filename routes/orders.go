package routes

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/order"
	"rent-n-go-backend/utils"
)

func OrderRoutes(r fiber.Router) {
	r.Get("/", order.History)
	r.Get("/active", order.HasActive)
	r.Post("/place", utils.InterceptRequest(new(order.PlaceOrderPayload)), order.Place)
}
