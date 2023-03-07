package routes

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin/auth"
	"rent-n-go-backend/utils"
)

func WebRoutes(r fiber.Router) {
	utils.RegisterWithPrefix(r, func(authRouter fiber.Router) {
		authRouter.Get("/login", auth.LoginView)
		authRouter.Post("/login", utils.InterceptWebRequest(new(auth.LoginRequest)), auth.LoginHandler)
	}, "auth", auth.Guest)

	r.Get("/logout", auth.Authed, auth.Logout)
}
