package tour

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

	cr := query.Tour

	var (
		tours []*models.Tour
		err   error
		total int64 = 0
	)

	qry := ServiceRepositories.Tour.BuildGetQuery()
	searchInt, _ := strconv.Atoi(search)

	if search != "" {
		qry = qry.Where(cr.Name.Like(search)).
			Or(cr.Price.Like(searchInt)).
			Or(cr.Desc.Like(search)).
			Or(cr.Stock.Like(searchInt))

		tours, err = qry.Scopes(utils.Paginate(c)).Find()
		total, _ = qry.Count()
	} else {
		tours, err = cr.Scopes(utils.Paginate(c)).Find()
		total, _ = cr.Count()
	}

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	sess := utils.Session.Provide(c)

	tourStock := map[uint]int64{}

	res := utils.Wrap(fiber.Map{
		"Tours":  tours,
		"Stocks": tourStock,
	}, c, sess).Pagination(total).Search(search).Message().Error()

	return admin.RenderTemplate(c, "tour/index", "Manage Tours", res.Get())
}

func Show(c *fiber.Ctx) error {
	tourId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/tours")
	}

	id := uint(tourId)

	tour, err := ServiceRepositories.Tour.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups tour not found")
		return c.RedirectBack("/admin/tours")
	}

	return admin.RenderTemplate(c, "tour/show", fmt.Sprintf("%s Detail", tour["name"]), tour)
}

func Create(c *fiber.Ctx) error {
	// Retrieve the car data from the database
	cars, err := query.Cars.Find()
	if err != nil {
		return utils.SafeThrow(c, err)
	}

	// Retrieve the driver data from the database
	drivers, err := query.Driver.Find()
	if err != nil {
		return utils.SafeThrow(c, err)
	}

	// Pass the car and driver data to the template
	return admin.RenderTemplate(c, "tour/form", "Create",
		utils.Wrap(fiber.Map{
			"Cars":    cars,
			"Drivers": drivers,
		}, nil, utils.Session.Provide(c)).Validation().Get())
}

func Store(c *fiber.Ctx) error {
	payload := utils.GetPayload[TourPayload](c)

	fileNames, err := utils.SaveMultiFilesFromPayload(c, "pictures", "tour")

	detail := "Tour added successfully!"

	if err != nil {
		detail = detail + " But, Some photos failed to upload."
	}

	sess := utils.Session.Provide(c)
	sess.SetSession("message", detail)

	tour := &models.Tour{
		Name:     payload.Name,
		Price:    payload.Price,
		Stock:    payload.Stock,
		Desc:     payload.Desc,
		CarId:    uint(payload.CarID),
		DriverId: uint(payload.DriverID),
	}

	if err := query.Tour.Create(tour); err != nil {
		return utils.SafeThrow(c, err)
	}

	for _, fileName := range fileNames {
		if err := BasicRepositories.Pictures.Insert(BasicRepositories.Tour, tour.ID, fileName); err != nil {
			return utils.SafeThrow(c, err)
		}
	}

	return c.Redirect("/admin/tours")
}

func Edit(c *fiber.Ctx) error {
	tourId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/tours")
	}

	id := uint(tourId)

	tour, err := ServiceRepositories.Tour.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups tour not found")
		return c.RedirectBack("/admin/tours")
	}

	// Retrieve the car data from the database
	cars, err := query.Cars.Find()
	if err != nil {
		return utils.SafeThrow(c, err)
	}

	// Retrieve the driver data from the database
	drivers, err := query.Driver.Find()
	if err != nil {
		return utils.SafeThrow(c, err)
	}

	response := utils.Wrap(tour, nil, utils.Session.Provide(c)).Error().Validation()

	data := response.Get()
	data["Cars"] = cars
	data["Drivers"] = drivers

	return admin.RenderTemplate(c, "tour/form", fmt.Sprintf("%s Edit", tour["name"]), data)
}

func Update(c *fiber.Ctx) error {
	tourId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/tours")
	}

	id := uint(tourId)

	tour, err := ServiceRepositories.Tour.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups tour not found")
		return c.RedirectBack("/admin/tours")
	}

	payload := utils.GetPayload[TourPayload](c)

	carId := c.FormValue("carId")
	driverId := c.FormValue("driverId")

	// Convert carId and driverId to uint
	carIDUint, err := strconv.ParseUint(carId, 10, 64)
	if err != nil {
		// Handle error
	}
	driverIDUint, err := strconv.ParseUint(driverId, 10, 64)
	if err != nil {
		// Handle error
	}

	qc := query.Tour
	if _, err := qc.Where(qc.ID.Eq(tour["id"].(uint))).Updates(&models.Tour{
		Name:     payload.Name,
		Price:    payload.Price,
		Stock:    payload.Stock,
		Desc:     payload.Desc,
		CarId:    uint(carIDUint),
		DriverId: uint(driverIDUint),
	}); err != nil {
		return utils.SafeThrow(c, err)
	}

	sess.SetSession("message", "Tour updated successfully")

	return c.Redirect("/admin/tours")
}

func Delete(c *fiber.Ctx) error {
	tourId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/tours")
	}

	id := uint(tourId)

	tour, err := ServiceRepositories.Tour.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups tour not found")
		return c.RedirectBack("/admin/tours")
	}

	if _, err := BasicRepositories.Features.DeleteByModuleId(BasicRepositories.Tour, tour["id"].(uint)); err != nil {
		return utils.SafeThrow(c, err)
	}

	if _, err := BasicRepositories.Pictures.DeleteByModuleId(BasicRepositories.Tour, tour["id"].(uint)); err != nil {
		return utils.SafeThrow(c, err)
	}

	qc := query.Tour

	if _, err := qc.Where(qc.ID.Eq(tour["id"].(uint))).Delete(); err != nil {
		return utils.SafeThrow(c, err)
	}

	sess.SetSession("message", "Tour deleted successfully")
	return c.RedirectBack("/admin/tours")
}
