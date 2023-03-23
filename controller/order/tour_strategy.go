package order

import "rent-n-go-backend/repositories/UserRepositories"

func tourStrategy(userId uint, payload PlaceOrderPayload) {
	err := UserRepositories.Order.CreateOrder(payload.StartPeriod, payload.EndPeriod, payload.PaymentMethod, userId)
}
