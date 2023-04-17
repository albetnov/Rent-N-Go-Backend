package car

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/BasicRepositories"
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

func Show(c *fiber.Ctx) error {
	carId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/cars")
	}

	id := uint(carId)

	car, err := ServiceRepositories.Car.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups car not found")
		return c.RedirectBack("/admin/cars")
	}

	return admin.RenderTemplate(c, "car/show", fmt.Sprintf("%s Detail", car["name"]), car)
}

func Create(c *fiber.Ctx) error {
	return admin.RenderTemplate(c, "car/form", "Create",
		utils.Wrap(fiber.Map{}, nil, utils.Session.Provide(c)).Validation().Get())
}

func Store(c *fiber.Ctx) error {
	payload := utils.GetPayload[CarPayload](c)

	fileNames, err := utils.SaveMultiFilesFromPayload(c, "pictures", "cars")

	detail := "Car added successfully!"

	if err != nil {
		detail = detail + " But, Some photos failed to upload."
	}

	sess := utils.Session.Provide(c)
	sess.SetSession("message", detail)

	car := &models.Cars{
		Name:  payload.Name,
		Price: payload.Price,
		Stock: payload.Stock,
		Desc:  payload.Desc,
	}

	if err := query.Cars.Create(car); err != nil {
		return utils.SafeThrow(c, err)
	}

	for _, fileName := range fileNames {
		if err := BasicRepositories.Pictures.Insert(BasicRepositories.Car, car.ID, fileName); err != nil {
			return utils.SafeThrow(c, err)
		}
	}

	for i, featureIcon := range payload.FeaturesIcon {
		if err := BasicRepositories.Features.Insert(BasicRepositories.Car, car.ID, featureIcon, payload.FeaturesLabel[i]); err != nil {
			return utils.SafeThrow(c, err)
		}
	}

	return c.RedirectBack("/admin/cars")
}

func Edit(c *fiber.Ctx) error {
	carId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/cars")
	}

	id := uint(carId)

	car, err := ServiceRepositories.Car.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups car not found")
		return c.RedirectBack("/admin/cars")
	}

	response := utils.Wrap(car, nil, utils.Session.Provide(c)).Error().Validation()

	return admin.RenderTemplate(c, "car/form", fmt.Sprintf("%s Edit", car["name"]), response.Get())
}

func Update(c *fiber.Ctx) error {
	carId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/cars")
	}

	id := uint(carId)

	car, err := ServiceRepositories.Car.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups car not found")
		return c.RedirectBack("/admin/cars")
	}

	payload := utils.GetPayload[CarPayload](c)

	qc := query.Cars
	if _, err := qc.Where(qc.ID.Eq(car["id"].(uint))).Updates(&models.Cars{
		Name:  payload.Name,
		Price: payload.Price,
		Stock: payload.Stock,
		Desc:  payload.Desc,
	}); err != nil {
		return utils.SafeThrow(c, err)
	}

	sess.SetSession("message", "Car updated successfully")

	return c.RedirectBack("/admin/cars")
}

func Delete(c *fiber.Ctx) error {
	carId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/cars")
	}

	id := uint(carId)

	car, err := ServiceRepositories.Car.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups car not found")
		return c.RedirectBack("/admin/cars")
	}

	if _, err := BasicRepositories.Features.DeleteByModuleId(BasicRepositories.Car, car["id"].(uint)); err != nil {
		return utils.SafeThrow(c, err)
	}

	if _, err := BasicRepositories.Pictures.DeleteByModuleId(BasicRepositories.Car, car["id"].(uint)); err != nil {
		return utils.SafeThrow(c, err)
	}

	qc := query.Cars

	if _, err := qc.Where(qc.ID.Eq(car["id"].(uint))).Delete(); err != nil {
		return utils.SafeThrow(c, err)
	}

	sess.SetSession("message", "Car deleted successfully")
	return c.RedirectBack("/admin/cars")
}
