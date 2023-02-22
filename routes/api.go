package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	router.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(viper.GetString("APP_KEY")),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"app":     utils.GetApp(),
				"message": "Invalid Credentials",
			})
		},
	}))

	router.Get("/restricted", func(ctx *fiber.Ctx) error {
		user := ctx.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		username := claims["name"].(string)
		return ctx.JSON(fiber.Map{
			"message":  "nt",
			"username": username,
		})
	})
}
