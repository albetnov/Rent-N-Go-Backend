package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"rent-n-go-backend/controller/admin/auth"
	"rent-n-go-backend/controller/admin/user"
	"rent-n-go-backend/utils"
	"time"
)

func WebRoutes(r fiber.Router) {
	csrfMiddleware := csrf.New(csrf.Config{
		KeyLookup:      "form:_csrf",
		CookieName:     "csrf_",
		CookieSameSite: "Strict",
		Expiration:     1 * time.Hour,
		ContextKey:     "token",
	})

	utils.RegisterWithPrefix(r, func(authRouter fiber.Router) {
		authRouter.Get("/login", auth.LoginView)
		authRouter.Post("/login", utils.InterceptWebRequest(new(auth.LoginRequest)), auth.LoginHandler)
	}, "auth", csrfMiddleware, auth.Guest)

	utils.RegisterWithPrefix(r, func(adminRouter fiber.Router) {
		adminRouter.Get("/dashboard", user.Dashboard)
		adminRouter.Get("/users/:id<int>", user.Show)
		adminRouter.Get("/users", user.Index)
		adminRouter.Get("/logout", auth.Logout)
	}, "admin", csrfMiddleware, auth.Authed)
}
