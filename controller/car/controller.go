package car

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/ServiceRepositories"
	"rent-n-go-backend/utils"
	"strconv"
	"strings"
)

func Recommendation(c *fiber.Ctx) error {
	cars, err := ServiceRepositories.Car.Ctx(c).GetRandom()

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "There's no recommended car for now",
			"error":   true,
		})
	}

	return c.JSON(fiber.Map{
		"cars":    cars,
		"message": "Recommendation fetched successfully",
	})
}

func Index(c *fiber.Ctx) error {
	search := c.Query("search", "")
	filterSeat, _ := strconv.Atoi(c.Query("seats", "0"))
	filterPrice, _ := strconv.Atoi(c.Query("price", "0"))

	cars, err := ServiceRepositories.Car.Ctx(c).GetAll(
		c,
		"%"+strings.ToLower(search)+"%",
		filterSeat,
		filterPrice,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Ups, no car available at the moment",
			"error":   true,
		})
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

	car, err := ServiceRepositories.Car.Ctx(c).GetById(uint(carId))

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
