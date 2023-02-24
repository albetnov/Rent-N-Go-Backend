package utils

import "github.com/gofiber/fiber/v2"

// RegisterWithPrefix
// A simple wrapper that will register the route under given prefix with "Before" middleware support
// That only applies to grouped route registered under this function.
func RegisterWithPrefix(r fiber.Router, routes func(r fiber.Router), prefix string, handlers ...fiber.Handler) {
	prefixedRoute := r.Group(prefix, handlers...)
	routes(prefixedRoute)
}

// Register
// A simple wrapper that will register the route without any prefix and "Before" middleware support (file scope).
// Please read how middleware works in https://docs.gofiber.io/api/app#route-handlers
func Register(r fiber.Router, routes func(r fiber.Router)) {
	routes(r)
}
