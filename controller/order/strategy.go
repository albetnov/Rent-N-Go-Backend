package order

import "github.com/gofiber/fiber/v2"

type strategy struct {
	Payload PlaceOrderPayload
	UserId  uint
}

func (s *strategy) Build(payload PlaceOrderPayload, userId uint) strategy {
	s.Payload = payload
	s.UserId = userId
	return *s
}

func (s strategy) UseStrategy(fn func(userId uint, payload PlaceOrderPayload) fiber.Map) fiber.Map {
	return fn(s.UserId, s.Payload)
}

var orderStrategy = strategy{}
