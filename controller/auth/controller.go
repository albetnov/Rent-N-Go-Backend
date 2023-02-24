package auth

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/query"
	"rent-n-go-backend/utils"
	"strconv"
	"strings"
	"time"
)

func Login(c *fiber.Ctx) error {
	u := query.User

	payload := utils.GetPayload[LoginPayload](c)

	user, err := u.Where(u.Email.Eq(payload.Email)).First()

	if err != nil || !utils.ComparePassword(payload.Password, user.Password) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Given Credentials not found",
			"action":  "NOT_FOUND",
		})
	}

	return generateToken(user, c)
}

func Refresh(c *fiber.Ctx) error {
	payload := utils.GetPayload[RefreshPayload](c)
	rt := query.RefreshToken
	u := query.User

	parsedString := strings.Split(payload.RefreshToken, "|")

	userId, err := strconv.ParseUint(parsedString[0], 10, 64)

	var id uint = uint(userId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Wrong nor invalid refresh token.",
			"action":  "INVALID_TOKEN",
		})
	}

	refreshToken, err := rt.Where(rt.UserID.Eq(id)).First()

	if refreshToken.ExpiredAt.Before(time.Now()) {
		rt.Where(rt.ID.Eq(refreshToken.ID)).Delete()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Ups your refresh token has expired.",
			"action":  "RE_LOGIN",
		})
	}

	user, err := u.Where(u.ID.Eq(id)).First()

	if refreshToken.Token != parsedString[1] || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Wrong nor invalid refresh token.",
			"action":  "INVALID_TOKEN",
		})
	}

	rt.Where(rt.ID.Eq(refreshToken.ID)).Delete()

	return generateToken(user, c)
}
