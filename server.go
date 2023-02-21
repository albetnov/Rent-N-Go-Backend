package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"rent-n-go-backend/routes"
	"rent-n-go-backend/utils"
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
		defer beforeHook(app).Close()
	} else {
		beforeHook(app)
	}

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(utils.GetApp())
	})

	api := app.Group("/api/v1")

	routes.ApiRoutes(api)

	afterHook(app)

	app.Listen(fmt.Sprintf(":%d", viper.GetInt("PORT")))
}
