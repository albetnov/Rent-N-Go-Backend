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

func (c car) buildGetQuery() query.ICarsDo {
	qc := query.Cars
	qp := query.Pictures
	qf := query.Features

	return qc.Preload(qc.Pictures.On(qp.Associate.Eq(BasicRepositories.Car))).
		Preload(qc.Features.On(qf.Associate.Eq(BasicRepositories.Car)))
}

func (c car) buildGenericResult(data *models.Cars, features, pictures []fiber.Map) fiber.Map {
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
}

func (c car) GetAll(ctx *fiber.Ctx) ([]fiber.Map, error) {
	result, err := c.buildGetQuery().
		Scopes(utils.Paginate(ctx)).Find()

	if err != nil {
		return nil, err
	}

	return utils.MapToServiceable(result, c.buildGenericResult), nil
}

func (c car) GetById(id uint) (fiber.Map, error) {
	qc := query.Cars
	result, err := c.buildGetQuery().Where(qc.ID.Eq(id)).First()

	if err != nil {
		return nil, err
	}

	return utils.MapToServiceableSingle(result, c.buildGenericResult), nil
}
