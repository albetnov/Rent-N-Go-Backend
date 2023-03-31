package order

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/UserRepositories"
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
		builder = qo.Preload(qo.User.On(query.User.Name.Like(search))).
			Where(qo.Type.Like(search)).
			Or(qo.Status.Like(search)).
			Or(qo.TotalAmount.Like(searchInt))
	} else {
		builder = qo.Preload(qo.User)
	}

	total, _ := builder.Count()

	order, err := builder.Scopes(utils.Paginate(c)).Order(qo.Status).Find()

	data := fiber.Map{
		"Orders": nil,
	}

	if err == nil {
		data["Orders"] = order
	}

	res := utils.Wrap(data, c, utils.Session.Provide(c)).Pagination(total).Message().Search(search).Csrf()

	return admin.RenderTemplate(c, "order/index", "Orders List", res.Get())
}

func getPicture(value []models.Pictures) *string {
	if len(value) > 0 {
		return &value[0].FileName
	}

	return nil
}

func Show(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/orders")
	}

	order, err := UserRepositories.Order.GetByOrderId(c, uint(orderId))

	res := utils.Wrap(fiber.Map{
		"ID":   order.ID,
		"Name": order.User.Name,
		"Car": fiber.Map{
			"Picture": getPicture(order.Car.Pictures),
			"Name":    order.Car.Name,
			"Desc":    order.Car.Desc,
			"Price":   order.Car.Price,
		},
		"Tour": fiber.Map{
			"Picture": getPicture(order.Tour.Pictures),
			"Name":    order.Tour.Name,
			"Desc":    order.Tour.Desc,
			"Price":   order.Tour.Price,
		},
		"Driver": fiber.Map{
			"Picture": getPicture(order.Driver.Pictures),
			"Name":    order.Driver.Name,
			"Desc":    order.Driver.Desc,
			"Price":   order.Driver.Price,
		},
		"TotalAmount":   order.TotalAmount,
		"Status":        order.Status,
		"Type":          order.Type,
		"PaymentMethod": order.PaymentMethod,
		"StartPeriod":   order.StartPeriod,
		"EndPeriod":     order.EndPeriod,
	}, c, sess).Csrf().Error().Message()

	return admin.RenderTemplate(c, "order/show", "Order Detail", res.Get())
}

func UpdateStatus(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)
	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/orders")
	}

	isComplete := c.FormValue("complete")
	isCancel := c.FormValue("cancel")
	isActive := c.FormValue("active")

	status := UserRepositories.ORDER_ACTIVE

	if isComplete != "" {
		status = UserRepositories.ORDER_COMPLETED
	}

	if isCancel != "" {
		status = UserRepositories.ORDER_CANCELLED
	}

	if isActive != "" {
		status = UserRepositories.ORDER_ACTIVE
	}

	if err := UserRepositories.Order.UpdateOrderStatus(uint(orderId), status); err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/orders")
	}

	sess.SetSession("message", "Order Status Updated Successfully!")
	return c.RedirectBack("/admin/orders")
}

func DeleteOrder(c *fiber.Ctx) error {
	orderId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/orders")
	}

	if err := UserRepositories.Order.DeleteOrder(uint(orderId)); err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/orders")
	}

	sess.SetSession("message", "Order has been deleted!")
	return c.RedirectBack("/admin/orders")
}
