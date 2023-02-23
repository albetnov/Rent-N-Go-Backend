package utils

import "github.com/gofiber/fiber/v2"

func RegisterWithPrefix(r fiber.Router, routes func(r fiber.Router), prefix string, handlers ...fiber.Handler) {
	prefixedRoute := r.Group(prefix, handlers...)
	routes(prefixedRoute)
}

func Register(r fiber.Router, routes func(r fiber.Router)) {
	routes(r)
}
