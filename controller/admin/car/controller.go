package car

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/ServiceRepositories"
	"rent-n-go-backend/utils"
	"strconv"
	"strings"
)

func Index(c *fiber.Ctx) error {
	search := strings.ToLower(c.Query("search"))

	cr := query.Cars

	var (
		qry   query.ICarsDo
		cars  []*models.Cars
		err   error
		total int64 = 0
	)

	searchInt, _ := strconv.Atoi(search)

	if search != "" {
		qry = cr.Where(cr.Name.Like(search)).
			Or(cr.Price.Like(searchInt)).
			Or(cr.Desc.Like(search)).
			Or(cr.Stock.Like(searchInt))

		cars, err = qry.Scopes(utils.Paginate(c)).Find()
		total, _ = qry.Count()
	} else {
		cars, err = cr.Scopes(utils.Paginate(c)).Find()
		total, _ = cr.Count()
	}

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	sess := utils.Session.Provide(c)

	carsStock := map[uint]int64{}

	for _, v := range cars {
		stock, _, _ := ServiceRepositories.Car.CheckStock(v.ID)
		carsStock[v.ID] = int64(v.Stock) - stock
	}

	res := utils.Wrap(fiber.Map{
		"Cars":   cars,
		"Stocks": carsStock,
	}, c, sess).Pagination(total).Search(search).Message().Error()

	return admin.RenderTemplate(c, "car/index", "Manage Cars", res.Get())
}
