package ServiceRepositories

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/BasicRepositories"
	"rent-n-go-backend/utils"
)

type car struct {
}

func (c car) GetAll(ctx *fiber.Ctx) ([]fiber.Map, error) {
	qc := query.Cars
	qp := query.Pictures
	qf := query.Features

	result, err := qc.
		Preload(qc.Pictures.On(qp.Associate.Eq(BasicRepositories.Car))).
		Preload(qc.Features.On(qf.Associate.Eq(BasicRepositories.Car))).
		Scopes(utils.Paginate(ctx)).Find()

	if err != nil {
		return nil, err
	}

	return utils.MapToServiceable(result, func(data *models.Cars, features, pictures []fiber.Map) fiber.Map {
		return fiber.Map{
			"id":         data.ID,
			"name":       data.Name,
			"stock":      data.Stock,
			"desc":       data.Desc,
			"price":      data.Price,
			"features":   features,
			"pictures":   pictures,
			"created_at": data.CreatedAt,
			"updated_at": data.UpdatedAt,
			"deleted_at": data.DeletedAt,
		}
	}), nil
}
