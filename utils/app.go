package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"log"
)

func GetApp() fiber.Map {
	return fiber.Map{
		"name":   "Rent-N-Go Backend",
		"slogan": "Your journey, our priority",
	}
}

func ShouldPanic(err error) {
	if viper.GetString("APP_ENV") == "production" {
		log.Fatalf("An error occurred: %v\n", err.Error())
	} else {
		panic(err)
	}
}

func RecordLog(err error) {
	if err != nil {
		log.Fatalf("Something went wrong: %v\n", err.Error())
	}
}
