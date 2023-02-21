package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"rent-n-go-backend/utils"
	"time"
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

	router.Post("/login", func(ctx *fiber.Ctx) error {
		type User struct {
			Username string
			Password string
		}

		user := new(User)

		if err := ctx.BodyParser(user); err != nil {
			return ctx.JSON(fiber.Map{
				"errors": err.Error(),
			})
		}

		claims := jwt.MapClaims{
			"name":  user.Username,
			"admin": false,
			"exp":   time.Now().Add(time.Minute * 30).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(viper.GetString("APP_KEY")))
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		return ctx.JSON(fiber.Map{
			"token": t,
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
