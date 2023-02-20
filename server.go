package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"os"
)

/**
* Load configuration
 */
func init() {
	// load from file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed, please re run:", e.Name)
		os.Exit(1)
	})
	viper.WatchConfig()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error regrading config file: %w", err))
	}

	// set default
	viper.SetDefault("PORT", 3000)
}

func main() {
	app := fiber.New()

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
