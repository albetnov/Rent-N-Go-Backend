package order

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/repositories/UserRepositories"
	"rent-n-go-backend/utils"
)

func History(c *fiber.Ctx) error {
	userId := utils.GetUserId(c)
	order, err := UserRepositories.Order.GetUserOrder(userId)

	if err != nil || len(order) < 1 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Ups, you seems like not having any order.",
			"error":   true,
		})
	}

	var clearedOrder []fiber.Map

	for _, v := range order {
		data := fiber.Map{
			"id":             v.ID,
			"total_amount":   v.TotalAmount,
			"status":         v.Status,
			"start_period":   v.StartPeriod,
			"end_period":     v.EndPeriod,
			"payment_method": v.PaymentMethod,
		}

		if v.CarId != 0 {
			data["car"] = v.Car
		}

		if v.DriverId != 0 {
			data["driver"] = v.Driver
		}

		if v.TourId != 0 {
			data["tour"] = v.Tour
		}

		clearedOrder = append(clearedOrder, data)
	}

	return c.JSON(fiber.Map{
		"data":    clearedOrder,
		"message": "Order fetched successfully",
	})
}
