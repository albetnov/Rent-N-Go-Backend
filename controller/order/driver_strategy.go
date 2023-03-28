package order

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/repositories/UserRepositories"
)

func handleDriverErrorResponse(err error) fiber.Map {
	if errors.Is(err, UserRepositories.DriverIsNotAvailableErr) {
		return driverNotAvailable()
	}

	return nil
}

func driverStrategy(c *fiber.Ctx, userId uint, payload PlaceOrderPayload) fiber.Map {
	err := UserRepositories.Order.CreateOrder(c, payload.StartPeriod, payload.EndPeriod, payload.PaymentMethod, userId).
		CreateDriverOrder(payload.CarId, payload.DriverId)

	if err != nil {
		if carHandled := handleCarErrorResponse(err); carHandled != nil {
			return carHandled
		}

		if driverHandled := handleDriverErrorResponse(err); driverHandled != nil {
			return driverHandled
		}

		return orderErr(err)
	}

	return orderOk()
}
