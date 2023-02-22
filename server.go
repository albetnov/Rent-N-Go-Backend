package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html"
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
	engine := html.New("./views", ".gohtml")

	app := fiber.New(fiber.Config{
		Views:   engine,
		AppName: "Rent-N-Go",
	})

	stream := beforeHook(app)

	if utils.IsProduction() {
		defer func() {
			err := stream.Close()
			if err != nil {
				utils.ShouldPanic(err)
			}
		}()
	}

	app.Get("/", monitor.New(monitor.Config{
		Title: "Rent N Go Backend Status",
	}))

	app.Static("/public", "./public")

	api := app.Group("/api/v1")

	routes.ApiRoutes(api)

	routes.WebRoutes(app)

	afterHook(app)

	err := app.Listen(fmt.Sprintf(":%d", viper.GetInt("PORT")))

	if err != nil {
		utils.ShouldPanic(err)
	}
}
