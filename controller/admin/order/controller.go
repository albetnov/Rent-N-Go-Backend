package order

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin"
	"rent-n-go-backend/query"
	"rent-n-go-backend/utils"
	"strconv"
	"strings"
)

func Index(c *fiber.Ctx) error {
	qo := query.Orders

	search := strings.ToLower(c.Query("search"))
	searchInt, _ := strconv.Atoi(search)

	var builder query.IOrdersDo

	if search != "" {
		builder = builder.Preload(qo.User.On(query.User.Name.Like(search))).
			Where(qo.Type.Like(search)).
			Or(qo.Status.Like(search)).
			Or(qo.TotalAmount.Like(searchInt))
	} else {
		builder = builder.Preload(qo.User)
	}

	total, _ := builder.Count()

	order, err := builder.Find()

	data := fiber.Map{
		"Orders": nil,
	}

	if err == nil {
		data["Orders"] = order
	}

	res := utils.Wrap(data, c, utils.Session.Provide(c)).Pagination(total).Message().Search(search)

	return admin.RenderTemplate(c, "order/index", "Orders List", res.Get())
}

func Show() {

}
