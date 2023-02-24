package routes

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/auth"
	"rent-n-go-backend/utils"
)

func AuthRoutes(r fiber.Router) {
	r.Post("/login", utils.InterceptRequest(new(auth.LoginPayload)), auth.Login)
	r.Post("/refresh", utils.InterceptRequest(new(auth.RefreshPayload)), auth.Refresh)
	r.Post("/register", utils.InterceptRequest(new(auth.RegisterPayload)), auth.Register)
}
