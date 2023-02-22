package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"rent-n-go-backend/query"
	"rent-n-go-backend/utils"
	"time"
)

func Login(c *fiber.Ctx) error {
	u := query.User

	payload := utils.GetPayload[RequestPayload](c)

	if _, err := u.Where(u.Username.Eq(payload.Username)).First(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Given Credentials not found",
		})
	}

	claims := jwt.MapClaims{
		"name":  payload.Username,
		"admin": false,
		"exp":   time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(viper.GetString("APP_KEY")))
	if err != nil {
		return utils.SafeThrow(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "Authenticated Successfully!",
		"token":   t,
	})
}
