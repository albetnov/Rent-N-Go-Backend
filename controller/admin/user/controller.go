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
		qry = u.Where(u.Name.Eq(search)).
			Or(u.Email.Eq(search)).
			Or(u.Role.Eq(search)).
			Or(u.PhoneNumber.Eq(search))

		users, err = qry.Scopes(utils.Paginate(c)).Find()
		total, _ = qry.Count()
	} else {
		users, err = u.Scopes(utils.Paginate(c)).Find()
		total, _ = u.Count()
	}

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	return admin.RenderTemplate(c, "users/index", "Manage Users",
		utils.Wrap(fiber.Map{
			"Users": users,
			"Error": utils.Session.Provide(c).GetFlash("error"),
		}).Ctx(c).Pagination(total).Search(search).Get())
}

func showUser(user *UserModels.User) fiber.Map {
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

	return admin.RenderTemplate(c, "users/show", fmt.Sprintf("%s Detail", user.Name), showUser(user))
}

func Create(c *fiber.Ctx) error {
	return admin.RenderTemplate(c, "users/form", "Create",
		utils.Wrap(fiber.Map{}).Validation(utils.Session.Provide(c)).Get())
}

func Store(c *fiber.Ctx) error {
	payload := utils.GetPayload[CreateUserPayload](c)

	simFile, err := utils.SaveFileFromPayload(c, "sim", utils.AssetPath("sim"))

	if err != nil && !strings.Contains(err.Error(), utils.NO_UPLOADED_FILE) {
		fmt.Println(err.Error())
	}

	fmt.Println(payload, simFile)

	return c.JSON(fiber.Map{"message": "test"})
}
