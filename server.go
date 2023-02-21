package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"rent-n-go-backend/routes"
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

	api := app.Group("/api/v1")

	routes.ApiRoutes(api)

	app.Listen(fmt.Sprintf(":%d", viper.GetInt("PORT")))
}
