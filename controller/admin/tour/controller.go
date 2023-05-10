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
		qry   query.ITourDo
		tours []*models.Tour
		err   error
		total int64 = 0
	)

	searchInt, _ := strconv.Atoi(search)

	if search != "" {
		qry = cr.Where(cr.Name.Like(search)).
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

	carNames := make(map[uint]string)
	for _, tour := range tours {
		tourCar, err := ServiceRepositories.Tour.GetByCarIdForName(tour.CarId)
		if err != nil {
			return utils.SafeThrow(c, err)
		}
		carNames[tour.CarId] = tourCar.Car.Name
	}

	driverNames := make(map[uint]string)
	for _, tour := range tours {
		tourDriver, err := ServiceRepositories.Tour.GetByDriverIdForName(tour.DriverId)
		if err != nil {
			return utils.SafeThrow(c, err)
		}
		driverNames[tour.DriverId] = tourDriver.Driver.Name
	}

	res := utils.Wrap(fiber.Map{
		"Tours":       tours,
		"Stocks":      tourStock,
		"CarNames":    carNames,
		"DriverNames": driverNames,
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
	return admin.RenderTemplate(c, "tour/form", "Create",
		utils.Wrap(fiber.Map{}, nil, utils.Session.Provide(c)).Validation().Get())
}

func Store(c *fiber.Ctx) error {
	payload := utils.GetPayload[TourPayload](c)

	fileNames, err := utils.SaveMultiFilesFromPayload(c, "pictures", "tours")

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

	response := utils.Wrap(tour, nil, utils.Session.Provide(c)).Error().Validation()

	return admin.RenderTemplate(c, "tour/form", fmt.Sprintf("%s Edit", tour["name"]), response.Get())
}

func Update(c *fiber.Ctx) error {
	tourId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/cars")
	}

	id := uint(tourId)

	tour, err := ServiceRepositories.Tour.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups car not found")
		return c.RedirectBack("/admin/cars")
	}

	payload := utils.GetPayload[TourPayload](c)

	// Get the carId and driverId values from the form data
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
