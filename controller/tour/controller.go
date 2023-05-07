package tour

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/ServiceRepositories"
	"rent-n-go-backend/utils"
	"strconv"
)

func Show(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
			"error":   true,
		})
	}

	result, err := ServiceRepositories.Tour.Ctx(c).GetById(uint(id))
	if err != nil {
		if e, ok := err.(*fiber.Error); ok && e.Code == fiber.StatusNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Tour not found",
				"error":   true,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching tour",
			"error":   true,
		})
	}

	return c.JSON(fiber.Map{
		"data":    result,
		"message": "Tour fetched successfully",
	})
}

func Index(c *fiber.Ctx) error {
	t := ServiceRepositories.Tour

	results, err := t.GetTours(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching tours",
			"error":   true,
		})
	}

	total, err := query.Tour.Count()
	if err != nil {
		return utils.SafeThrow(c, err)
	}

	res := utils.Wrap(fiber.Map{
		"data":    results,
		"message": "Tours fetched successfully",
	}, c).WithMeta(total)

	return c.JSON(res.Get())
}

func Stocks(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid ID",
			"error":   true,
		})
	}

	availableStock, _, err := ServiceRepositories.Tour.Ctx(c).CheckStock(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error checking tour stock",
			"error":   true,
		})
	}

	return c.JSON(fiber.Map{
		"available_stock": availableStock,
		"message":         "Tour stock checked successfully",
	})
}
