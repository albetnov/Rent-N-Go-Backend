package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"rent-n-go-backend/utils"
)

func ApiRoutes(r fiber.Router) {
	// use cors
	router := r.Use(cors.New())

	// set default global router
	utils.SetGlobalRouter(router)

	utils.RegisterWithPrefix(AuthRoutes, "/auth")
}
