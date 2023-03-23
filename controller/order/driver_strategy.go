package order

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/repositories/UserRepositories"
	"sync"
)

func driverStrategy(res chan<- fiber.Map, mtx *sync.Mutex, userId uint, payload PlaceOrderPayload) {
	mtx.Lock()
	defer mtx.Unlock()

	err := UserRepositories.Order.CreateOrder(payload.StartPeriod, payload.EndPeriod, payload.PaymentMethod, userId).
		CreateDriverOrder(payload.CarId, payload.DriverId)

	if err != nil {
		if errors.Is(err, UserRepositories.CarIsOutOfStockErr) {
			res <- carNotAvailable()
			return
		}

		if errors.Is(err, UserRepositories.DriverIsNotAvailableErr) {
			res <- driverNotAvailable()

			return
		}

		res <- orderErr(err)
		return
	}

	res <- orderOk()
}
