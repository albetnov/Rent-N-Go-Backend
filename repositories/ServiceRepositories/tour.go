package ServiceRepositories

import (
	"github.com/gofiber/fiber/v2"
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

func (t tour) BuildGetQuery() query.ITourDo {
	qt := query.Tour
	qp := query.Pictures

	return qt.
		Preload(qt.Pictures.On(qp.Associate.Eq(BasicRepositories.Tour))).
		Preload(qt.Driver).
		Preload(qt.Driver.Pictures.On(qp.Associate.Eq(BasicRepositories.Driver))).
		Preload(qt.Car).
		Preload(qt.Car.Pictures.On(qp.Associate.Eq(BasicRepositories.Car)))
}

func (t tour) buildGenericResult(data *models.Tour, features, pictures []fiber.Map) fiber.Map {
	totalPrice := data.Price + data.Driver.Price + data.Car.Price
	carStock, _, _ := Car.CheckStock(data.Car.ID)
	return fiber.Map{
		"id":       data.ID,
		"name":     data.Name,
		"desc":     data.Desc,
		"price":    totalPrice,
		"features": features,
		"pictures": pictures,
		"car": utils.MapToServiceableSingle(t.c, data.Car, func(data models.Cars, features, pictures []fiber.Map) fiber.Map {
			return fiber.Map{
				"id":         data.ID,
				"name":       data.Name,
				"stock":      data.Stock,
				"desc":       data.Desc,
				"price":      data.Price,
				"pictures":   pictures,
				"seats":      data.Seats,
				"baggages":   data.Baggage,
				"created_at": data.CreatedAt,
				"updated_at": data.UpdatedAt,
				"deleted_at": data.DeletedAt,
				"hold_stock": carStock,
			}
		}),
		"driver": utils.MapToServiceableSingle(t.c, data.Driver, func(data models.Driver, features, pictures []fiber.Map) fiber.Map {
			return fiber.Map{
				"id":         data.ID,
				"name":       data.Name,
				"desc":       data.Desc,
				"price":      data.Price,
				"pictures":   pictures,
				"created_at": data.CreatedAt,
				"updated_at": data.UpdatedAt,
				"deleted_at": data.DeletedAt,
			}
		}),
		"stock":      data.Stock,
		"created_at": data.CreatedAt,
		"updated_at": data.UpdatedAt,
		"deleted_at": data.DeletedAt,
	}
}

func (t tour) GetById(id uint) (fiber.Map, error) {
	qt := query.Tour
	result, err := t.BuildGetQuery().Where(qt.ID.Eq(id)).First()

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

func (t tour) GetTours(c *fiber.Ctx, search string, price int) ([]fiber.Map, error) {
	qt := query.Tour

	tourQuery := t.BuildGetQuery().
		Scopes(utils.Paginate(c))

	if search != "%%" {
		tourQuery = tourQuery.Where(qt.Name.Like(search)).Or(qt.Desc.Like(search))
	}

	if price > 0 {
		tourQuery = tourQuery.Where(qt.Price.Gte(price))
	}

	results, err := tourQuery.Find()

	if err != nil {
		return nil, err
	}

	serviceableResults := make([]fiber.Map, len(results))
	for i, result := range results {
		serviceableResults[i] = utils.MapToServiceableSingle(c, result, t.buildGenericResult)
	}

	return serviceableResults, nil
}
