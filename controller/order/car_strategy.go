package order

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/repositories/UserRepositories"
)

func carStrategy(userId uint, payload PlaceOrderPayload) fiber.Map {
	err := UserRepositories.Order.
		CreateOrder(payload.StartPeriod, payload.EndPeriod, payload.PaymentMethod, userId).
		CreateCarOrder(payload.CarId)

	if err != nil {
		if errors.Is(err, UserRepositories.CarIsOutOfStockErr) {
			return carNotAvailable()
		}

		return orderErr(err)
	}

	return orderOk()
}
