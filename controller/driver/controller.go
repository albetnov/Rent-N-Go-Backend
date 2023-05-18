package driver

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/ServiceRepositories"
	"rent-n-go-backend/utils"
	"strconv"
)

func Index(c *fiber.Ctx) error {
	search := c.Query("search", "")

	drivers, err := ServiceRepositories.Driver.Ctx(c).GetAll(c, search)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Ups, no driver available at the moment",
			"error":   true,
		})
	}

	total, err := query.Driver.Count()
	if err != nil {
		return utils.SafeThrow(c, err)
	}

	res := utils.Wrap(fiber.Map{
		"data":    drivers,
		"message": "drivers fetched successfully",
	}, c).WithMeta(total)

	return c.JSON(res.Get())
}

func Show(c *fiber.Ctx) error {
	driverId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "drivers not found or invalid id given",
			"error":   true,
		})
	}

	drivers, err := ServiceRepositories.Driver.Ctx(c).GetById(uint(driverId))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "drivers not found",
			"error":   true,
		})
	}

	res := fiber.Map{
		"data":    drivers,
		"message": "drivers found!"}

	return c.JSON(res)
}
