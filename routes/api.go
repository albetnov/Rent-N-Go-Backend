package routes

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"rent-n-go-backend/utils"
)

type Testing struct {
	Name       string `validate:"required,min=3,max=32"`
	Salary     int    `validate:"required,number"`
	IsEmployee *bool  `validate:"required"`
}

func ApiRoutes(r fiber.Router) {
	utils.RegisterWithPrefix(r, AuthRoutes, "/auth")

	r.Post("/test", utils.InterceptRequest(new(Testing)), func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"message": "mantap",
		})
	})

	r.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(viper.GetString("APP_KEY")),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"app":     utils.GetApp(),
				"message": "Invalid Credentials",
				"error":   err.Error(),
				"details": "It is also possible that the route you're looking for is not found.",
			})
		}},
	))

	r.Get("/restricted", func(ctx *fiber.Ctx) error {
		user := ctx.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		username := claims["name"].(string)
		return ctx.JSON(fiber.Map{
			"message":  "nt",
			"username": username,
		})
	})
}
