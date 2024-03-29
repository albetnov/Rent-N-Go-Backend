package ServiceRepositories

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/BasicRepositories"
	"rent-n-go-backend/utils"
)

type car struct {
	c *fiber.Ctx
}

func (c *car) Ctx(ctx *fiber.Ctx) *car {
	c.c = ctx
	return c
}

func (c car) buildGetQuery() query.ICarsDo {
	qc := query.Cars
	qp := query.Pictures

	return qc.Preload(qc.Pictures.On(qp.Associate.Eq(BasicRepositories.Car)))
}

func (c car) buildGenericResult(data *models.Cars, features, pictures []fiber.Map) fiber.Map {
	stock, _, _ := c.CheckStock(data.ID)

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
		"hold_stock": stock,
	}
}

func (c car) GetRandom() ([]fiber.Map, error) {
	result, err := c.buildGetQuery().RandomizeWithLimit(6)

	if err != nil {
		return nil, err
	}

	return utils.MapToServiceable(c.c, result, c.buildGenericResult), nil
}

func (c car) GetAll(ctx *fiber.Ctx, search string, seats, price int) ([]fiber.Map, error) {
	qc := query.Cars
	carQuery := c.buildGetQuery().
		Scopes(utils.Paginate(ctx))

	if search != "%%" {
		carQuery = carQuery.Where(qc.Name.Like(search)).Or(qc.Desc.Like(search))
	}

	if seats > 0 {
		carQuery = carQuery.Where(qc.Seats.Eq(seats))
	}

	if price > 0 {
		carQuery = carQuery.Where(qc.Price.Gte(price))
	}

	result, err := carQuery.Find()

	if err != nil {
		return nil, err
	}

	return utils.MapToServiceable(c.c, result, c.buildGenericResult), nil
}

func (c car) GetById(id uint) (fiber.Map, error) {
	qc := query.Cars
	result, err := c.buildGetQuery().Where(qc.ID.Eq(id)).First()

	if err != nil {
		return nil, err
	}

	return utils.MapToServiceableSingle(c.c, result, c.buildGenericResult), nil
}

func (c car) CheckStock(id uint) (int64, *models.Cars, error) {
	qo := query.Orders
	qc := query.Cars

	stock, err := qc.Where(qc.ID.Eq(id)).First()
	if err != nil {
		return 0, nil, err
	}

	totalOrder, err := qo.Scopes(activeOrder).Where(qo.CarId.Eq(id)).Count()

	if err != nil {
		return 0, nil, err
	}

	result := int64(stock.Stock) - totalOrder

	if result < 0 {
		result = 0
	}

	return result, stock, nil
}
