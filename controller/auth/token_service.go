package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"rent-n-go-backend/models"
	"rent-n-go-backend/repositories"
	"rent-n-go-backend/utils"
	"time"
)

func generateToken(user *models.User, c *fiber.Ctx) error {
	tokenExpiredAt := time.Now().Add(time.Hour * 24).Unix()
	refreshTokenExpiredAt := time.Now().Add(time.Hour * 720)

	refreshToken := models.RefreshToken{
		Token:     utils.GenerateRandomString(100),
		ExpiredAt: refreshTokenExpiredAt,
		UserID:    user.ID,
	}

	repositories.UpdateOrCreateTokenByUserId(user.ID, &refreshToken)

	claims := jwt.MapClaims{
		"name": user.Name,
		"role": user.Role,
		"id":   user.ID,
		"exp":  tokenExpiredAt,
		"nik":  user.Nik,
		"sim":  user.Sim,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(viper.GetString("APP_KEY")))
	if err != nil {
		return utils.SafeThrow(c, err)
	}

	return c.JSON(fiber.Map{
		"message":                  "Authenticated Successfully!",
		"token":                    t,
		"token_expired_at":         tokenExpiredAt,
		"refresh_token_expired_at": refreshTokenExpiredAt.Unix(),
		"refresh_token":            fmt.Sprintf("%d|%s", user.ID, refreshToken.Token),
	})
}
