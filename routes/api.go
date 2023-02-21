package routes

import "github.com/gofiber/fiber/v2"

var globalRouter fiber.Router

func register(fn func(r fiber.Router)) {
	fn(globalRouter)
}

func ApiRoutes(r fiber.Router) {
	// set default global router
	globalRouter = r

	register(AuthRoutes)
}
