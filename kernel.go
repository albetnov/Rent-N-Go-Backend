package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/viper"
	"os"
	"path"
	"rent-n-go-backend/utils"
)

/*
Create a log directory and its file if not exist, then return
the file instance.
*/
func getLogFile() *os.File {
	currentDir, _ := utils.GetCurrentDir()

	fileDir := path.Join(currentDir, "logs")

	os.MkdirAll(fileDir, 0700)

	file, err := os.OpenFile(path.Join(fileDir, "log.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(fmt.Sprintf("Error opening file: %v", err.Error()))
	}

	return file
}

/*
*
Return the corresponding Logger Output based on Application Environment
Return file if in production, return stdOut otherwise.
*/
func getLogOutput() *os.File {
	if viper.GetString("APP_ENV") == "production" {
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
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: viper.GetString("APP_KEY"),
	}))

	// Load requestId middleware
	app.Use(requestid.New())

	file := getLogOutput()

	// Load logger middleware
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status}:${method} -> ${path} ::${locals:requestid} \n",
		Output: file,
	}))

	// Only if in production, then recover the app.
	if viper.GetString("APP_ENV") == "production" {
		app.Use(fiberRecover.New())
	}

	return file
}

/*
*
An global middleware that initiated at the end of routing.
*/
func afterHook(app *fiber.App) {
	// set up 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"context": "Rent-N-Go Backend",
			"message": "Ups, can't find that!",
			"status":  fiber.StatusNotFound,
		})
	})
}
