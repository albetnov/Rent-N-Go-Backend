package utils

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

const NO_UPLOADED_FILE = "there is no uploaded file associated with the given key"

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

func GetErrorMessage(err error) string {
	errorMessage := "Can't proceed your request"

	if !IsProduction() {
		errorMessage = err.Error()
	}

	return errorMessage
}

// SafeThrow
// Safely throw an error to end UserModels in production
// Safely throw an error complete with the message in development mode.
func SafeThrow(w *fiber.Ctx, err error) error {
	errorMessage := GetErrorMessage(err)

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
// Determine if the UserModels wanted json response or an html response.
func WantsJson(c *fiber.Ctx) bool {
	return c.Get("Accept") == "application/json"
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

// GetUser
// Get the current logged in user from given context
func GetUser(c *fiber.Ctx) jwt.MapClaims {
	user := c.Locals("user")

	if user != nil {
		return user.(*jwt.Token).Claims.(jwt.MapClaims)
	}

	return nil
}

// GetUserId
// Return the UserId of the current user
func GetUserId(c *fiber.Ctx) uint {
	auth := GetUser(c)

	if auth == nil {
		return 0
	}

	authId := uint(auth["id"].(float64))

	return authId
}

func validateAndSaveFile(c *fiber.Ctx, file *multipart.FileHeader, assetDirectory string) (string, error) {
	reader, err := file.Open()
	if err != nil {
		return "", err
	}

	defer reader.Close()

	err = CheckMimes(reader, []string{"image/jpg", "image/png", "image/jpeg"})

	if err != nil {
		return "", err
	}

	salt := uuid.New().String()

	fileName := salt + file.Filename

	basePath := path.Join(AssetPath(assetDirectory))

	if err = os.MkdirAll(basePath, 0700); err != nil {
		return "", err
	}

	c.SaveFile(file, path.Join(AssetPath(assetDirectory), fileName))

	return fileName, nil
}

// SaveFileFromPayload
// Simplify image saving by automatic salt and validate the given file.
func SaveFileFromPayload(c *fiber.Ctx, payload string, assetDirectory string) (string, error) {
	file, err := c.FormFile(payload)

	if err != nil {
		return "", err
	}

	return validateAndSaveFile(c, file, assetDirectory)
}

// SaveMultiFileFromPayload
// Simply multiple image saving by automatically salt and validate each given files from payload
func SaveMultiFilesFromPayload(c *fiber.Ctx, payload string, assetDirectory string) ([]string, error) {
	form, err := c.MultipartForm()

	if err != nil {
		return nil, err
	}

	files := form.File[payload]

	if len(files) <= 0 {
		return nil, errors.New(fmt.Sprintf("Invalid. No %s being provided.", payload))
	}

	fileNames := []string{}

	for _, file := range files {
		fileName, err := validateAndSaveFile(c, file, assetDirectory)
		if err != nil {
			return fileNames, err
		}

		fileNames = append(fileNames, fileName)
	}

	return fileNames, nil
}

func ParseISO8601Date(date string) time.Time {
	layout := "2006-01-02T15:04:05Z07:00"
	result, err := time.Parse(layout, date)

	if err != nil {
		ShouldPanic(err)
	}

	return result
}
