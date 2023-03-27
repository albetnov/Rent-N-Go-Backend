package order

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/repositories/UserRepositories"
)

func handleCarErrorResponse(err error) fiber.Map {
	if errors.Is(err, UserRepositories.CarIsOutOfStockErr) {
		return carNotAvailable()
	}

	if errors.Is(err, UserRepositories.CarNotFound) {
		return carNotFound()
	}

	return nil
}

func carStrategy(userId uint, payload PlaceOrderPayload) fiber.Map {
	err := UserRepositories.Order.
		CreateOrder(payload.StartPeriod, payload.EndPeriod, payload.PaymentMethod, userId).
		CreateCarOrder(payload.CarId)

	if err != nil {
		if handled := handleCarErrorResponse(err); handled != nil {
			return handled
		}

		return orderErr(err)
	}

	return orderOk()
}
