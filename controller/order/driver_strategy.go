package order

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/repositories/UserRepositories"
)

func driverStrategy(userId uint, payload PlaceOrderPayload) fiber.Map {
	err := UserRepositories.Order.CreateOrder(payload.StartPeriod, payload.EndPeriod, payload.PaymentMethod, userId).
		CreateDriverOrder(payload.CarId, payload.DriverId)

	if err != nil {
		if errors.Is(err, UserRepositories.CarIsOutOfStockErr) {
			return carNotAvailable()
		}

		if errors.Is(err, UserRepositories.DriverIsNotAvailableErr) {
			return driverNotAvailable()

		}

		return orderErr(err)
	}

	return orderOk()
}
