package routes

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/home"
)

func WebRoutes(r fiber.Router) {
	r.Get("/home", home.Index)
}
