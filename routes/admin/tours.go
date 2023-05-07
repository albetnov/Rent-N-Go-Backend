package admin

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin/tour"
	"rent-n-go-backend/utils"
)

func ToursRoutes(r fiber.Router) {
	r.Get("/", tour.Index)
	r.Get("/create", tour.Create)
	r.Post("/create", utils.InterceptWebRequest(new(tour.TourPayload)), tour.Store)
	r.Get("/edit/:id<int>", tour.Edit)
	r.Post("/edit/:id<int>", utils.InterceptWebRequest(new(tour.TourPayload)), tour.Update)
	r.Get("/delete/:id<int>", tour.Delete)
	r.Get("/:id<int>", tour.Show)
}
