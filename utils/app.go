package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// GetApp
// Return the application about boilerplate response in Map.
func GetApp() fiber.Map {
	return fiber.Map{
		"name":   "Rent-N-Go Backend",
		"slogan": "Your journey, our priority",
	}
}

// ShouldPanic
// Determine whenever the app should panic or not (Depend on application state)
// Do not panic in production, log them instead.
func ShouldPanic(err error) {
	if IsProduction() {
		log.Fatalf("An error occurred: %v\n", err.Error())
	} else {
		panic(err)
	}
}

// RecordLog
// Smartly record any log if error occur.
func RecordLog(err error) {
	if err != nil {
		log.Fatalf("Something went wrong: %v\n", err.Error())
	}
}

// IsProduction
// Determine the application state if it's running in production mode or not.
func IsProduction() bool {
	return viper.GetString("APP_ENV") == "production"
}

// SafeThrow
// Safely throw an error to end user in production
// Safely throw an error complete with the message in development mode.
func SafeThrow(w *fiber.Ctx, err error) error {
	errorMessage := "Can't proceed your request"

	if !IsProduction() {
		errorMessage = err.Error()
	}

	statusCode := fiber.StatusInternalServerError

	w.Status(statusCode)

	if WantsJson(w) {
		return w.JSON(fiber.Map{
			"message": "Something went wrong",
			"error":   errorMessage,
		})
	}

	return w.Render("error", fiber.Map{
		"Code":    statusCode,
		"Message": errorMessage,
	})
}

// WantsJson
// Determine if the user wanted json response or an html response.
func WantsJson(c *fiber.Ctx) bool {
	return c.Get("Content-Type") == "application/json"
}

// HashPassword
// Hash a given string using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes), err
}

// ComparePassword
// Check whenever the hashed and known password is the same.
func ComparePassword(knownPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(knownPassword))

	return err == nil
}

// GenerateRandomString
// Make a random string from given length
func GenerateRandomString(length int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890!@#$%^&*()_+{}[];:.,/?`~'\"")

	char := make([]rune, length)
	for word := range char {
		char[word] = chars[r.Intn(len(chars))]
	}

	return string(char)
}

func GetUser(c *fiber.Ctx) jwt.MapClaims {
	user := c.Locals("user")

	if user != nil {
		return user.(*jwt.Token).Claims.(jwt.MapClaims)
	}

	return nil
}

func GetUserId(c *fiber.Ctx) uint {
	auth := GetUser(c)

	if auth == nil {
		return 0
	}

	authId := uint(auth["id"].(float64))

	return authId
}
