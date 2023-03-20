package car

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/ServiceRepositories"
	"rent-n-go-backend/utils"
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
		"cars":    cars,
		"message": "Car fetched successfully",
	}, c).WithMeta(total)

	return c.JSON(res.Get())
}
