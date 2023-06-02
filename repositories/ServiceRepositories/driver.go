package ServiceRepositories

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gen"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/BasicRepositories"
	"rent-n-go-backend/utils"
)

type driver struct {
	c  *fiber.Ctx
	db gen.Dao // Add db field of type gen.Dao
}

func (d *driver) Ctx(c *fiber.Ctx) *driver {
	d.c = c
	return d
}

func (d driver) buildGetQuery() query.IDriverDo {
	qd := query.Driver
	qp := query.Pictures

	return qd.Preload(qd.Pictures.On(qp.Associate.Eq(BasicRepositories.Driver)))
}

func (d driver) buildGenericResult(data *models.Driver, features, pictures []fiber.Map) fiber.Map {
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
}

func (d driver) GetById(id uint) (fiber.Map, error) {
	qd := query.Driver
	result, err := d.buildGetQuery().Where(qd.ID.Eq(id)).First()

	if err != nil {
		return nil, err
	}

	return utils.MapToServiceableSingle(d.c, result, d.buildGenericResult), nil
}

func (d driver) CheckAvailability(id uint) bool {
	qo := query.Orders

	total, _ := qo.Scopes(activeOrder).Where(qo.DriverId.Eq(id)).Count()

	return total > 0
}
func (d driver) GetAll(c *fiber.Ctx, search string, price int) ([]fiber.Map, error) {
	qd := query.Driver

	driverQuery := d.buildGetQuery()

	if search != "%%" {
		driverQuery = driverQuery.Where(qd.Name.Like(search)).Or(qd.Desc.Like(search))
	}

	if price > 0 {
		driverQuery = driverQuery.Where(qd.Price.Gte(price))
	}
	results, err := driverQuery.Find()

	if err != nil {
		return nil, err
	}

	serviceableResults := make([]fiber.Map, len(results))
	for i, result := range results {
		serviceableResults[i] = utils.MapToServiceableSingle(c, result, d.buildGenericResult)
	}

	return serviceableResults, nil
}
