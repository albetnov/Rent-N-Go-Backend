package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

/*
*
The Rent-N-Go-Backend Entrypoint
This function loads kernel middlewares
as well as close unused instance.

This function will also register all defined routes.
*/
func main() {
	app := fiber.New()

	if viper.GetString("APP_ENV") == "production" {
		defer registerGlobalMiddlewares(app).Close()
	} else {
		registerGlobalMiddlewares(app)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println("Logging Middleware!")
		return c.Next()
	}, func(c *fiber.Ctx) error {
		c.SendStatus(200)

		return c.JSON(fiber.Map{
			"message": "Welcome to Rent-N-Go Backend Entrypoint",
			"status":  200,
		})
	})

	app.Listen(fmt.Sprintf(":%d", viper.GetInt("PORT")))
}
