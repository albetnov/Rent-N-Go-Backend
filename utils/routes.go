package utils

import "github.com/gofiber/fiber/v2"

var globalRouter fiber.Router

// SetGlobalRouter
// Should only be used for internal routing only.
func SetGlobalRouter(r fiber.Router) {
	globalRouter = r
}

// RegisterWithPrefix
// Register a routes with prefix.
// This function exist thanks to Go not supporting overloading. What a pain.
// Ref: https://github.com/golang/go/wiki/GoForCPPProgrammers
func RegisterWithPrefix(fn func(r fiber.Router), prefix string) {
	prefixedRoute := globalRouter.Group(prefix)
	fn(prefixedRoute)
}

// Register
// Register a routes without prefix.
// This function exist thanks to Go not supporting overloading. What a pain.
// Ref: https://github.com/golang/go/wiki/GoForCPPProgrammers
func Register(fn func(r fiber.Router)) {
	fn(globalRouter)
}
