package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log"
	"os"
	"path"
	"rent-n-go-backend/query"
	"rent-n-go-backend/utils"
	"runtime"
)

/*
Create a log directory and its file if not exist, then return
the file instance.
*/
func getLogFile() *os.File {
	currentDir, _ := utils.GetCurrentDir()

	fileDir := path.Join(currentDir, "logs")

	err := os.MkdirAll(fileDir, 0700)

	if err != nil {
		utils.ShouldPanic(err)
	}

	file, err := os.OpenFile(path.Join(fileDir, "log.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		utils.ShouldPanic(err)
	}

	return file
}

/*
*
Return the corresponding Logger Output based on Application Environment
Return file if in production, return stdOut otherwise.
*/
func getLogOutput() *os.File {
	if utils.IsProduction() {
		return getLogFile()
	}

	return os.Stdout
}

/*
Register all application globals middleware.
return file instance to be deferred by server entry point.

Initiated at beginning of routing.
*/
func registerGlobalMiddlewares(app *fiber.App) *os.File {
	// Load encryptCookie middleware
	//app.Use(encryptcookie.New(encryptcookie.Config{
	//	Key: viper.GetString("APP_KEY"),
	//}))

	// Load requestId middleware
	app.Use(requestid.New())

	file := getLogOutput()

	// Load logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status}:${method} -> ${path} ::${locals:requestid} \n",
		Output: file,
	}))

	// Only if in production, then recover the app.
	if utils.IsProduction() {
		app.Use(fiberRecover.New())
	}

	return file
}

// RegisterViewFunc Register custom view utilities
func RegisterViewFunc() map[string]interface{} {
	return map[string]interface{}{
		"when": func(firstCond any, value any, fallback any) any {
			if firstCond != nil {
				return value
			}

			return fallback
		},
		"inc": func(a int) int {
			return a + 1
		},
		"dec": func(a int) int {
			return a - 1
		},
	}
}

// beforeHook bootstrap any process before begin serving.
func beforeHook(app *fiber.App) *os.File {
	// Log some welcome message
	log.Println("Welcome to Rent-N-Go Backend!")
	log.Println("Running in:", runtime.Version(), "Using:", runtime.GOOS)

	if utils.IsProduction() {
		log.Println("App is running in production mode.")
	}

	// Satisfy database connection
	utils.SatisfiesDbConnection()

	// set default db for query
	query.SetDefault(utils.GetDb())

	// register the global middleware
	file := registerGlobalMiddlewares(app)

	// get argument of app
	args := os.Args[1:]

	processMigration(args)

	utils.Session.InitStore()

	return file
}

/*
*
An global middleware that initiated at the end of routing.
*/
func afterHook(app *fiber.App) {
	// set up 404 handler
	app.Use(func(c *fiber.Ctx) error {
		statusCode := fiber.StatusNotFound
		message := "Ups, can't find that!"

		c.Status(statusCode)

		if utils.WantsJson(c) {
			return c.JSON(fiber.Map{
				"app":     utils.GetApp(),
				"message": message,
				"status":  statusCode,
			})
		}

		return c.Render("error", fiber.Map{
			"Code":    statusCode,
			"Message": message,
		})
	})
}
