package ServiceRepositories

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gen"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/BasicRepositories"
	"rent-n-go-backend/utils"
)

type tour struct {
}

func (t tour) buildGetQuery(db gen.Dao) gen.Dao {
	qt := query.Tour
	qp := query.Pictures
	qf := query.Features

	return db.Preload(qt.Pictures.On(qp.Associate.Eq(BasicRepositories.Driver))).
		Preload(qt.Features.On(qf.Associate.Eq(BasicRepositories.Driver))).
		Preload(qt.Driver).
		Preload(qt.Car)
}

func (t tour) buildGenericResult(data *models.Tour, features, pictures []fiber.Map) fiber.Map {
	return fiber.Map{
		"id":         data.ID,
		"name":       data.Name,
		"desc":       data.Desc,
		"price":      data.Price,
		"features":   features,
		"pictures":   pictures,
		"car":        data.Car,
		"driver":     data.Driver,
		"stock":      data.Stock,
		"created_at": data.CreatedAt,
		"updated_at": data.UpdatedAt,
		"deleted_at": data.DeletedAt,
	}
}

func (t tour) GetById(id uint) (fiber.Map, error) {
	qt := query.Tour
	result, err := qt.Scopes(t.buildGetQuery).Where(qt.ID.Eq(id)).First()

	if err != nil {
		return nil, err
	}

	return utils.MapToServiceableSingle(result, t.buildGenericResult), nil
}

func (t tour) CheckStock(id uint) (int64, *models.Tour, error) {
	qo := query.Orders
	qt := query.Tour

	total, _ := qo.Scopes(utils.ActiveOrder).Where(qo.TourId.Eq(id)).Count()
	tour, err := qt.Where(qt.ID.Eq(id)).First()

	return int64(tour.Stock) - total, tour, err
}
