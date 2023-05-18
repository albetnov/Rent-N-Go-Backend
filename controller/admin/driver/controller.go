package driver

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

	cr := query.Driver

	var (
		qry    query.IDriverDo
		driver []*models.Driver
		err    error
		total  int64 = 0
	)

	searchInt, _ := strconv.Atoi(search)

	if search != "" {
		qry = cr.Where(cr.Name.Like(search)).
			Or(cr.Price.Like(searchInt)).
			Or(cr.Desc.Like(search))

		driver, err = qry.Scopes(utils.Paginate(c)).Find()
		total, _ = qry.Count()
	} else {
		driver, err = cr.Scopes(utils.Paginate(c)).Find()
		total, _ = cr.Count()
	}

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	sess := utils.Session.Provide(c)

	res := utils.Wrap(fiber.Map{
		"Driver": driver,
	}, c, sess).Pagination(total).Search(search).Message().Error()

	return admin.RenderTemplate(c, "driver/index", "Manage Driver", res.Get())
}

func Show(c *fiber.Ctx) error {
	driverId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/driver")
	}

	id := uint(driverId)

	driver, err := ServiceRepositories.Driver.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups driver not found")
		return c.RedirectBack("/admin/driver")
	}

	return admin.RenderTemplate(c, "driver/show", fmt.Sprintf("%s Detail", driver["name"]), driver)
}

func Create(c *fiber.Ctx) error {
	return admin.RenderTemplate(c, "driver/form", "Create",
		utils.Wrap(fiber.Map{}, nil, utils.Session.Provide(c)).Validation().Get())
}

func Store(c *fiber.Ctx) error {
	payload := utils.GetPayload[DriverPayload](c)

	fileNames, err := utils.SaveMultiFilesFromPayload(c, "pictures", "driver")

	detail := "Driver added successfully!"

	if err != nil {
		detail = detail + " But, Some photos failed to upload."
	}

	sess := utils.Session.Provide(c)
	sess.SetSession("message", detail)

	driver := &models.Driver{
		Name:  payload.Name,
		Price: payload.Price,
		Desc:  payload.Desc,
	}

	if err := query.Driver.Create(driver); err != nil {
		return utils.SafeThrow(c, err)
	}

	for _, fileName := range fileNames {
		if err := BasicRepositories.Pictures.Insert(BasicRepositories.Driver, driver.ID, fileName); err != nil {
			return utils.SafeThrow(c, err)
		}
	}

	return c.Redirect("/admin/driver")
}

func Edit(c *fiber.Ctx) error {
	driverId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/driver")
	}

	id := uint(driverId)

	driver, err := ServiceRepositories.Driver.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups driver not found")
		return c.RedirectBack("/admin/driver")
	}

	response := utils.Wrap(driver, nil, utils.Session.Provide(c)).Error().Validation()

	return admin.RenderTemplate(c, "driver/form", fmt.Sprintf("%s Edit", driver["name"]), response.Get())
}

func Update(c *fiber.Ctx) error {
	driverId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/driver")
	}

	id := uint(driverId)

	driver, err := ServiceRepositories.Driver.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups driver not found")
		return c.RedirectBack("/admin/driver")
	}

	payload := utils.GetPayload[DriverPayload](c)

	qd := query.Driver
	if _, err := qd.Where(qd.ID.Eq(driver["id"].(uint))).Updates(&models.Driver{
		Name:  payload.Name,
		Price: payload.Price,
		Desc:  payload.Desc,
	}); err != nil {
		return utils.SafeThrow(c, err)
	}

	sess.SetSession("message", "Driver updated successfully")

	return c.Redirect("/admin/driver")
}

func Delete(c *fiber.Ctx) error {
	driverId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/driver")
	}

	id := uint(driverId)

	driver, err := ServiceRepositories.Driver.Ctx(c).GetById(id)

	if err != nil {
		sess.SetSession("error", "Ups driver not found")
		return c.RedirectBack("/admin/driver")
	}

	if _, err := BasicRepositories.Features.DeleteByModuleId(BasicRepositories.Driver, driver["id"].(uint)); err != nil {
		return utils.SafeThrow(c, err)
	}

	if _, err := BasicRepositories.Pictures.DeleteByModuleId(BasicRepositories.Driver, driver["id"].(uint)); err != nil {
		return utils.SafeThrow(c, err)
	}

	qd := query.Driver

	if _, err := qd.Where(qd.ID.Eq(driver["id"].(uint))).Delete(); err != nil {
		return utils.SafeThrow(c, err)
	}

	sess.SetSession("message", "Driver deleted successfully")
	return c.RedirectBack("/admin/driver")
}
