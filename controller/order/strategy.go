package order

import "github.com/gofiber/fiber/v2"

const IsCar = "car"
const IsDriver = "driver"
const IsTour = "tour"

type strategy struct {
	Payload PlaceOrderPayload
	UserId  uint
	c       *fiber.Ctx
}

func (s *strategy) Build(ctx *fiber.Ctx, payload PlaceOrderPayload, userId uint) strategy {
	s.Payload = payload
	s.UserId = userId
	s.c = ctx
	return *s
}

func (s strategy) UseStrategy(fn func(c *fiber.Ctx, userId uint, payload PlaceOrderPayload) fiber.Map) fiber.Map {
	return fn(s.c, s.UserId, s.Payload)
}

var orderStrategy = strategy{}
