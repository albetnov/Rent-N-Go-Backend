package order

import "github.com/gofiber/fiber/v2"

func orderOk() fiber.Map {
	return fiber.Map{
		"message": "Order created successfully!",
		"status":  fiber.StatusOK,
	}
}

func orderErr(err error) fiber.Map {
	return fiber.Map{
		"message": "Something went wrong",
		"details": err.Error(),
		"status":  fiber.StatusInternalServerError,
	}
}

func driverNotAvailable() fiber.Map {
	return fiber.Map{
		"message": "Ups, seems like the driver is not available at the moment",
		"status":  fiber.StatusBadRequest,
	}
}

func carNotAvailable() fiber.Map {
	return fiber.Map{
		"message": "The car you've selected is out of stock, please choose another car.",
		"status":  fiber.StatusBadRequest,
	}
}

func carNotFound() fiber.Map {
	return fiber.Map{
		"message": "Car not found.",
		"status":  fiber.StatusNotFound,
	}
}

func tourNotAvailable() fiber.Map {
	return fiber.Map{
		"message": "The tour you've looking for is not longer available",
		"status":  fiber.StatusBadRequest,
	}
}
