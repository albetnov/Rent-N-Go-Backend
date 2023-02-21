package utils

import "github.com/gofiber/fiber/v2"

func GetApp() fiber.Map {
	return fiber.Map{
		"name":   "Rent-N-Go Backend",
		"slogan": "Your journey, our priority",
	}
}
