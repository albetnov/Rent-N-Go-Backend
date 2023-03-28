package order

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/repositories/UserRepositories"
)

func tourStrategy(c *fiber.Ctx, userId uint, payload PlaceOrderPayload) fiber.Map {
	err := UserRepositories.Order.CreateOrder(c, payload.StartPeriod, payload.EndPeriod, payload.PaymentMethod, userId).
		CreateTourOrder(payload.TourId)

	if err != nil {
		if carHandled := handleCarErrorResponse(err); carHandled != nil {
			return carHandled
		}

		if driverHandled := handleDriverErrorResponse(err); driverHandled != nil {
			return driverHandled
		}

		if errors.Is(err, UserRepositories.TourIsNotAvailableErr) {
			return tourNotAvailable()
		}

		return orderErr(err)
	}

	return orderOk()
}
