package routes

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/utils"
)

func ApiRoutes(r fiber.Router) {
	// set default global router
	utils.SetGlobalRouter(r)

	utils.RegisterWithPrefix(AuthRoutes, "/auth")
}
