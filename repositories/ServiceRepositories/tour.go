package ServiceRepositories

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gen"
	"gorm.io/gorm"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/BasicRepositories"
	"rent-n-go-backend/utils"
)

type tour struct {
	c  *fiber.Ctx
	db *gorm.DB
}

func (t *tour) Ctx(c *fiber.Ctx) *tour {
	t.c = c
	return t
}

func (t tour) BuildGetQuery(db gen.Dao) gen.Dao {
	qt := query.Tour
	qp := query.Pictures

	return db.
		Preload(qt.Pictures.On(qp.Associate.Eq(BasicRepositories.Tour))).
		Preload(qt.Driver).
		Preload(qt.Driver.Pictures.On(qp.Associate.Eq(BasicRepositories.Driver))).
		Preload(qt.Car).
		Preload(qt.Car.Pictures.On(qp.Associate.Eq(BasicRepositories.Car)))
}

func (t tour) buildGenericResult(data *models.Tour, features, pictures []fiber.Map) fiber.Map {
	totalPrice := data.Price + data.Driver.Price + data.Car.Price

	return fiber.Map{
		"id":         data.ID,
		"name":       data.Name,
		"desc":       data.Desc,
		"price":      totalPrice,
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
	result, err := qt.Scopes(t.BuildGetQuery).Where(qt.ID.Eq(id)).First()

	if err != nil {
		return nil, err
	}

	return utils.MapToServiceableSingle(t.c, result, t.buildGenericResult), nil
}

func (t tour) CheckStock(id uint) (int64, *models.Tour, error) {
	qo := query.Orders
	qt := query.Tour

	total, _ := qo.Scopes(activeOrder).Where(qo.TourId.Eq(id)).Count()
	tour, err := qt.Where(qt.ID.Eq(id)).First()

	return int64(tour.Stock) - total, tour, err
}

func (t tour) GetTours(c *fiber.Ctx) ([]fiber.Map, error) {
	qt := query.Tour
	results, err := qt.Scopes(t.BuildGetQuery).Find()

	if err != nil {
		return nil, err
	}

	serviceableResults := make([]fiber.Map, len(results))
	for i, result := range results {
		serviceableResults[i] = utils.MapToServiceableSingle(c, result, t.buildGenericResult)
	}

	return serviceableResults, nil
}
