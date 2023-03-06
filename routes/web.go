package routes

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin/auth"
	"rent-n-go-backend/controller/home"
)

func WebRoutes(r fiber.Router) {
	r.Get("/dashboard", auth.Authed, home.Index)
	r.Get("/login", auth.Guest, auth.LoginView)
}
