package order

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/repositories/UserRepositories"
	"sync"
)

func carStrategy(res chan<- fiber.Map, mtx *sync.Mutex, userId uint, payload PlaceOrderPayload) {
	mtx.Lock()
	defer mtx.Unlock()

	err := UserRepositories.Order.
		CreateOrder(payload.StartPeriod, payload.EndPeriod, payload.PaymentMethod, userId).
		CreateCarOrder(payload.CarId)

	if err != nil {
		if errors.Is(err, UserRepositories.CarIsOutOfStockErr) {
			res <- carNotAvailable()
		} else {
			res <- orderErr(err)
		}
		return
	}

	res <- orderOk()
}
