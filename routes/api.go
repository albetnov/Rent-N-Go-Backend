package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"rent-n-go-backend/utils"
)

type Testing struct {
	Name       string `validate:"required,min=3,max=32"`
	Salary     int    `validate:"required,number"`
	IsEmployee *bool  `validate:"required"`
}

func ApiRoutes(r fiber.Router) {
	// use cors
	router := r.Use(cors.New())

	// set default global router
	utils.SetGlobalRouter(router)

	utils.RegisterWithPrefix(AuthRoutes, "/auth")

	router.Post("/test", utils.InterceptRequest(new(Testing)), func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "mantap",
		})
	})
}
