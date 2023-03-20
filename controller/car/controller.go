package car

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/ServiceRepositories"
	"rent-n-go-backend/utils"
	"strconv"
)

func Index(c *fiber.Ctx) error {
	cars, err := ServiceRepositories.Car.GetAll(c)

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	total, err := query.Cars.Count()
	if err != nil {
		return utils.SafeThrow(c, err)
	}

	res := utils.Wrap(fiber.Map{
		"data":    cars,
		"message": "Car fetched successfully",
	}, c).WithMeta(total)

	return c.JSON(res.Get())
}

func Show(c *fiber.Ctx) error {
	carId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Car not found or invalid id given",
			"error":   true,
		})
	}

	car, err := ServiceRepositories.Car.GetById(uint(carId))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Car not found",
			"error":   true,
		})
	}

	res := fiber.Map{
		"data":    car,
		"message": "Car found!"}

	return c.JSON(res)
}
