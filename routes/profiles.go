package routes

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/profile"
	"rent-n-go-backend/utils"
)

func ProfileRoutes(r fiber.Router) {
	r.Get("/current", profile.CurrentUser)
	r.Get("/status", profile.CompletionStatus)
	r.Put("/update", utils.InterceptRequest(new(profile.UpdateProfilePayload)), profile.UpdateProfile)
	r.Put("/update/nik", utils.InterceptRequest(new(profile.CompleteNikPayload)), profile.UpdateNik)
	r.Post("/update/sim", profile.UpdateSim)
	r.Delete("/delete", profile.DeleteAccount)
}
