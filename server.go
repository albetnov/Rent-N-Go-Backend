package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	// init html engine
	engine := html.New("./views", ".gohtml").AddFuncMap(RegisterViewFunc())
	// init fiber
	app := fiber.New(fiber.Config{
		Views:   engine,
		AppName: "Rent-N-Go",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			if err == fiber.ErrRequestEntityTooLarge {
				if utils.WantsJson(ctx) {
					return ctx.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
						"message": "Payload is too large, rejecting.",
					})
				}
				return ctx.Status(fiber.StatusRequestEntityTooLarge).Render("error", fiber.Map{
					"Code":    fiber.StatusRequestEntityTooLarge,
					"Message": "Your file is too large. Rejecting.",
				})
			}

			return fiber.DefaultErrorHandler(ctx, err)
		},
	})

	// run beforeHook of kernel.go
	stream := beforeHook(app)

	// check if app is in production
	if utils.IsProduction() {
		defer func() {
			// close the stream (returned by beforeHook only when in production)
			err := stream.Close()
			if err != nil {
				utils.ShouldPanic(err)
			}
		}()
	}

	// init monitor views
	app.Get("/", monitor.New(monitor.Config{
		Title: "Rent N Go Backend Status",
	}))

	// init static files serving
	app.Static("/public", utils.PublicPath())

	// init api namespace with cors middleware.
	api := app.Group("/api/v1", cors.New())

	// register routes
	routes.ApiRoutes(api)
	routes.WebRoutes(app)

	// init after hook that will be run after routes.
	afterHook(app)

	// serve the app
	err := app.Listen(fmt.Sprintf(":%d", viper.GetInt("PORT")))

	if err != nil {
		utils.ShouldPanic(err)
	}
}
