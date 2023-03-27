package order

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/repositories/UserRepositories"
	"rent-n-go-backend/utils"
	"sync"
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
			"type":           v.Type,
		}

		if v.CarId != nil {
			data["car"] = v.Car
		}

		if v.DriverId != nil {
			data["driver"] = v.Driver
		}

		if v.TourId != nil {
			data["tour"] = v.Tour
		}

		clearedOrder = append(clearedOrder, data)
	}

	return c.JSON(fiber.Map{
		"data":    clearedOrder,
		"message": "Order fetched successfully",
	})
}

func Place(c *fiber.Ctx) error {
	payload := utils.GetPayload[PlaceOrderPayload](c)

	userId := utils.GetUserId(c)
	if alreadyHasOrder := UserRepositories.Order.HasOrder(userId); alreadyHasOrder {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "You already have an order!",
			"status":  fiber.StatusBadRequest,
		})
	}

	var res fiber.Map
	mtx := new(sync.Mutex)

	mtx.Lock()
	orderStrategy.Build(payload, userId)
	if payload.Type == IsCar {
		res = orderStrategy.UseStrategy(carStrategy)
	} else if payload.Type == IsDriver {
		res = orderStrategy.UseStrategy(driverStrategy)
	} else {
		res = orderStrategy.UseStrategy(tourStrategy)
	}
	mtx.Unlock()

	return c.Status(res["status"].(int)).JSON(res)
}
