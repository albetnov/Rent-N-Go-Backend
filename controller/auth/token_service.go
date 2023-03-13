package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/UserRepositories"
	"rent-n-go-backend/utils"
	"time"
)

func generateToken(currentUser *UserModels.User, c *fiber.Ctx) error {
	tokenExpiredAt := time.Now().Add(time.Hour * 24).Unix()
	refreshTokenExpiredAt := time.Now().Add(time.Hour * 720)

	refreshToken := UserModels.RefreshToken{
		Token:     utils.GenerateRandomString(100),
		ExpiredAt: refreshTokenExpiredAt,
		UserID:    currentUser.ID,
	}

	UserRepositories.RefreshToken.UpdateOrCreateByUserId(currentUser.ID, &refreshToken)

	photo, _ := query.User.Photo.Model(currentUser).Find()

	photo.PhotoPath = utils.FormatUrl(c, photo.PhotoPath, "user")

	claims := jwt.MapClaims{
		"role": currentUser.Role,
		"id":   currentUser.ID,
		"exp":  tokenExpiredAt,
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
		"refresh_token":            fmt.Sprintf("%d|%s", currentUser.ID, refreshToken.Token),
	})
}
