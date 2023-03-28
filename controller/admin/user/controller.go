package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/controller/admin"
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/UserRepositories"
	"rent-n-go-backend/utils"
	"strconv"
	"strings"
)

func Dashboard(c *fiber.Ctx) error {
	return admin.RenderTemplate(c, "dashboard", "Dashboard", nil)
}

func Index(c *fiber.Ctx) error {
	search := strings.ToLower(c.Query("search"))

	u := query.User

	var (
		qry   query.IUserDo
		users []*UserModels.User
		err   error
		total int64 = 0
	)

	if search != "" {
		qry = u.Where(u.Name.Like(search)).
			Or(u.Email.Like(search)).
			Or(u.Role.Like(search)).
			Or(u.PhoneNumber.Like(search))

		users, err = qry.Scopes(utils.Paginate(c)).Find()
		total, _ = qry.Count()
	} else {
		users, err = u.Scopes(utils.Paginate(c)).Find()
		total, _ = u.Count()
	}

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	sess := utils.Session.Provide(c)

	return admin.RenderTemplate(c, "users/index", "Manage Users",
		utils.Wrap(fiber.Map{
			"Users": users,
		}, c, sess).Pagination(total).Search(search).Message().Error().Get())
}

func showUser(c *fiber.Ctx, user *UserModels.User) fiber.Map {
	user.Sim.FilePath = utils.FormatUrl(c, user.Sim.FilePath, "sim")
	user.Photo.PhotoPath = utils.FormatUrl(c, user.Photo.PhotoPath, "user")

	return fiber.Map{
		"Name":        user.Name,
		"Email":       user.Email,
		"PhoneNumber": user.PhoneNumber,
		"Role":        user.Role,
		"SIM":         user.Sim,
		"NIK":         user.Nik,
		"Photo":       user.Photo,
		"ID":          user.ID,
	}
}

func Show(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/users")
	}

	id := uint(userId)

	user, err := UserRepositories.User.GetAllById(id)

	if err != nil {
		sess.SetSession("error", "Ups user not found")
		return c.RedirectBack("/admin/users")
	}

	return admin.RenderTemplate(c, "users/show", fmt.Sprintf("%s Detail", user.Name), showUser(c, user))
}

func Create(c *fiber.Ctx) error {
	return admin.RenderTemplate(c, "users/form", "Create",
		utils.Wrap(fiber.Map{}, nil, utils.Session.Provide(c)).Validation().Get())
}

func Store(c *fiber.Ctx) error {
	payload := utils.GetPayload[CreateUserPayload](c)

	hashed, err := utils.HashPassword(payload.Password)

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	userData := UserModels.User{
		Role:        payload.Role,
		Email:       payload.Email,
		Name:        payload.Name,
		PhoneNumber: strconv.FormatInt(int64(payload.PhoneNumber), 10),
		Password:    hashed,
	}

	if err := UserRepositories.User.Create(&userData); err != nil {
		return utils.SafeThrow(c, err)
	}

	sess := utils.Session.Provide(c)

	if err := UserRepositories.Sim.OptionalCreate(c,
		"sim",
		sess,
		userData.ID,
		"/admin/users"); err != nil {
		return err
	}

	UserRepositories.Nik.OptionalCreate(payload.Nik, userData.ID)

	if err := UserRepositories.User.OptionalCreatePhoto(c,
		sess,
		"photo",
		userData.ID,
		"/admin/users"); err != nil {
		return err
	}

	sess.SetSession("message", "User created successfully.")
	return c.Redirect("/admin/users")
}

func Edit(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))

	sess := utils.Session.Provide(c)

	if err != nil {
		sess.SetSession("error", err.Error())
		return c.RedirectBack("/admin/users")
	}

	id := uint(userId)

	user, err := UserRepositories.User.GetAllById(id)

	if err != nil {
		sess.SetSession("error", "Ups user not found")
		return c.RedirectBack("/admin/users")
	}

	response := utils.Wrap(showUser(c, user), nil, utils.Session.Provide(c)).Error().Validation().Get()

	return admin.RenderTemplate(c, "users/form", fmt.Sprintf("Edit %s", user.Name), response)
}

func Update(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	formattedUserId := uint(userId)

	if _, err = UserRepositories.User.GetById(formattedUserId); err != nil {
		return c.Status(fiber.StatusNotFound).Render("error", fiber.Map{
			"Code":    fiber.StatusNotFound,
			"Message": "Ups, user not found",
		})
	}

	payload := utils.GetPayload[UpdateUserPayload](c)

	userData := UserModels.User{
		Role:        payload.Role,
		Email:       payload.Email,
		Name:        payload.Name,
		PhoneNumber: strconv.FormatInt(int64(payload.PhoneNumber), 10),
	}

	if payload.Password != "" {
		hashed, err := utils.HashPassword(payload.Password)
		if err != nil {
			return utils.SafeThrow(c, err)
		}

		userData.Password = hashed
	}

	if err := UserRepositories.User.UpdateById(c, formattedUserId, &userData); err != nil {
		return utils.SafeThrow(c, err)
	}

	sess := utils.Session.Provide(c)

	if err := UserRepositories.Sim.OptionalCreate(c,
		"sim",
		sess,
		userData.ID,
		"/admin/users"); err != nil {
		return err
	}

	UserRepositories.Nik.OptionalCreate(payload.Nik, userData.ID)

	if err := UserRepositories.User.OptionalCreatePhoto(c,
		sess,
		"photo",
		userData.ID,
		"/admin/users"); err != nil {
		return err
	}

	sess.SetSession("message", "User edited successfully.")
	return c.Redirect("/admin/users")
}

func Destroy(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	parsedId := uint(userId)

	if err := UserRepositories.User.DeleteById(parsedId); err != nil {
		return c.Status(fiber.StatusNotFound).Render("error", fiber.Map{
			"Code":    fiber.StatusNotFound,
			"Message": "Corresponding user not found",
		})
	}

	utils.Session.Provide(c).SetSession("message", "User deleted successfully")
	return c.RedirectBack("/admin/users")
}
