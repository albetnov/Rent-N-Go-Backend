package auth

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"rent-n-go-backend/models/user"
	"rent-n-go-backend/query"
	userRepository "rent-n-go-backend/repositories/user"
	"rent-n-go-backend/utils"
	"time"
)

func generateToken(currentUser *user.User, c *fiber.Ctx) error {
	tokenExpiredAt := time.Now().Add(time.Hour * 24).Unix()
	refreshTokenExpiredAt := time.Now().Add(time.Hour * 720)

	refreshToken := user.RefreshToken{
		Token:     utils.GenerateRandomString(100),
		ExpiredAt: refreshTokenExpiredAt,
		UserID:    currentUser.ID,
	}

	userRepository.RefreshToken.UpdateOrCreateByUserId(currentUser.ID, &refreshToken)

	photo, _ := query.User.Photo.Model(currentUser).Find()
	nik, _ := query.User.Nik.Model(currentUser).Find()
	sim, _ := query.User.Sim.Model(currentUser).Find()

	claims := jwt.MapClaims{
		"name":  currentUser.Name,
		"role":  currentUser.Role,
		"id":    currentUser.ID,
		"exp":   tokenExpiredAt,
		"nik":   nik,
		"sim":   sim,
		"phone": currentUser.PhoneNumber,
		"email": currentUser.Email,
		"photo": photo,
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
