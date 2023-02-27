package auth

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/models"
	"rent-n-go-backend/repositories"
	"rent-n-go-backend/utils"
	"strconv"
	"strings"
	"time"
)

func Login(c *fiber.Ctx) error {
	payload := utils.GetPayload[LoginPayload](c)

	user, err := repositories.GetUserByEmail(payload.Email)

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

	parsedString := strings.Split(payload.RefreshToken, "|")

	userId, err := strconv.ParseUint(parsedString[0], 10, 64)

	var id uint = uint(userId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Wrong nor invalid refresh token.",
			"action":  "INVALID_TOKEN",
		})
	}

	refreshToken, err := repositories.GetTokenByUserId(id)

	repositories.DeleteTokenByTokenId(refreshToken.ID)

	if refreshToken.ExpiredAt.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Ups your refresh token has expired.",
			"action":  "RE_LOGIN",
		})
	}

	user, err := repositories.GetUserById(id)

	if refreshToken.Token != parsedString[1] || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Wrong nor invalid refresh token.",
			"action":  "INVALID_TOKEN",
		})
	}

	return generateToken(user, c)
}

func Register(c *fiber.Ctx) error {
	payload := utils.GetPayload[RegisterPayload](c)
	if _, err := repositories.GetUserByEmailOrPhone(payload.Email, payload.PhoneNumber); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Ups, email already exist",
			"action":  "CHANGE_EMAIL",
		})
	}

	password, err := utils.HashPassword(payload.Password)

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	user := models.User{
		Name:        payload.Name,
		PhoneNumber: payload.PhoneNumber,
		Email:       payload.Email,
		Role:        "user",
		Password:    password,
	}

	repositories.CreateUser(&user)

	return c.JSON(fiber.Map{
		"message": "User created successfully!",
		"user":    user,
	})
}
