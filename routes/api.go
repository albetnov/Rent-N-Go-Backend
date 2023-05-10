package routes

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/spf13/viper"
	"rent-n-go-backend/utils"
)

func ApiRoutes(r fiber.Router) {
	utils.RegisterWithPrefix(r, AuthRoutes, "/auth")
	utils.RegisterWithPrefix(r, CarsRoutes, "/cars")
	utils.RegisterWithPrefix(r, TourRoutes, "/tours")

	r.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(viper.GetString("APP_KEY")),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"app":     utils.GetApp(),
				"message": "Invalid Credentials",
				"error":   err.Error(),
				"details": "It is also possible that the route you're looking for is not found.",
			})
		}},
	))

	utils.RegisterWithPrefix(r, ProfileRoutes, "/profiles")
	utils.RegisterWithPrefix(r, OrderRoutes, "/orders")
}
