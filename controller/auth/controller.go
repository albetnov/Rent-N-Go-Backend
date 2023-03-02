package auth

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/repositories/UserRepositories"
	"rent-n-go-backend/utils"
	"strconv"
	"strings"
	"time"
)

func Login(c *fiber.Ctx) error {
	payload := utils.GetPayload[LoginPayload](c)

	currentUser, err := UserRepositories.User.GetByEmail(payload.Email)

	if err != nil || !utils.ComparePassword(payload.Password, currentUser.Password) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Given Credentials not found",
			"action":  "NOT_FOUND",
		})
	}

	return generateToken(currentUser, c)
}

func Refresh(c *fiber.Ctx) error {
	payload := utils.GetPayload[RefreshPayload](c)

	parsedString := strings.Split(payload.RefreshToken, "|")

	userId, err := strconv.ParseUint(parsedString[0], 10, 64)

	var id = uint(userId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Wrong nor invalid refresh token.",
			"action":  "INVALID_TOKEN",
		})
	}

	refreshToken, err := UserRepositories.RefreshToken.GetByUserId(id)

	UserRepositories.RefreshToken.DeleteByTokenId(refreshToken.ID)

	if refreshToken.ExpiredAt.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Ups your refresh token has expired.",
			"action":  "RE_LOGIN",
		})
	}

	currentUser, err := UserRepositories.User.GetById(id)

	if refreshToken.Token != parsedString[1] || err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Wrong nor invalid refresh token.",
			"action":  "INVALID_TOKEN",
		})
	}

	return generateToken(currentUser, c)
}

func Register(c *fiber.Ctx) error {
	payload := utils.GetPayload[RegisterPayload](c)
	if _, err := UserRepositories.User.GetByEmailOrPhone(payload.Email, payload.PhoneNumber); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Ups, email already exist",
			"action":  "CHANGE_EMAIL",
		})
	}

	password, err := utils.HashPassword(payload.Password)

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	currentUser := UserModels.User{
		Name:        payload.Name,
		PhoneNumber: payload.PhoneNumber,
		Email:       payload.Email,
		Role:        "currentUser",
		Password:    password,
	}

	UserRepositories.User.Create(&currentUser)

	return c.JSON(fiber.Map{
		"message":     "User created successfully!",
		"currentUser": currentUser,
	})
}
