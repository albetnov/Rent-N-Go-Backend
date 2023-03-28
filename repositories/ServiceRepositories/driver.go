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
	c *fiber.Ctx
}

func (d *driver) Ctx(c *fiber.Ctx) *driver {
	d.c = c
	return d
}

func (d driver) buildGetQuery(db gen.Dao) gen.Dao {
	qd := query.Driver
	qp := query.Pictures
	qf := query.Features

	return db.Preload(qd.Pictures.On(qp.Associate.Eq(BasicRepositories.Driver))).
		Preload(qd.Features.On(qf.Associate.Eq(BasicRepositories.Driver)))
}

func (d driver) buildGenericResult(data *models.Driver, features, pictures []fiber.Map) fiber.Map {
	return fiber.Map{
		"id":         data.ID,
		"name":       data.Name,
		"desc":       data.Desc,
		"price":      data.Price,
		"features":   features,
		"pictures":   pictures,
		"created_at": data.CreatedAt,
		"updated_at": data.UpdatedAt,
		"deleted_at": data.DeletedAt,
	}
}

func (d driver) GetById(id uint) (fiber.Map, error) {
	qd := query.Driver
	result, err := qd.Scopes(d.buildGetQuery).Where(qd.ID.Eq(id)).First()

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
