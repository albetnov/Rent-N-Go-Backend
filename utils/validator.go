package utils

import (
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
	"io"
	"strings"
	"time"
)

type ErrorResponse struct {
	FailedFields string
	Tag          string
	Value        string
}

var validate = validator.New()

//func registerValidator() {
//validate.RegisterValidation()
//}

const BODY_DATA = "body_data"

// customValidator
// set the custom validator.
func customValidator() {
	// Allow password to empty, while not add minimum 8 to it's length.
	validate.RegisterValidation("passwordable", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "" {
			return true
		}

		return len(fl.Field().String()) >= 8
	})

	// Check if date comply ISO8601 standard
	validate.RegisterValidation("ISO8601date", func(fl validator.FieldLevel) bool {
		layout := "2006-01-02T15:04:05Z07:00"
		_, err := time.Parse(layout, fl.Field().String())
		return err == nil
	})

	// Validate date after (only works with ISO8601 format)
	validate.RegisterValidation("afteriso", func(fl validator.FieldLevel) bool {
		layout := "2006-01-02T15:04:05Z07:00"
		fieldValue, _, _, _ := fl.GetStructFieldOK2()
		startDate, err := time.Parse(layout, fieldValue.String())
		if err != nil {
			return false
		}
		endDate, err := time.Parse(layout, fl.Field().String())
		if err != nil {
			return false
		}
		// exactly one day after start. No matter the hours difference.
		return endDate.YearDay() == startDate.YearDay()+1 && endDate.Year() == startDate.Year()
	})
}

// validateStruct
// Validate an payload to a given struct, return error if something wrong, and return empty if all passed.
func validateStruct(data any) []*ErrorResponse {
	var errorResponses []*ErrorResponse
	customValidator()

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedFields = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errorResponses = append(errorResponses, &element)
		}
	}

	return errorResponses
}

// validateWebStruct
// Validate given payload but instead of producing ErrorResponse pointer, instead this function return array
// of string.
func validateWebStruct(data any) []string {
	var errorResponses []string
	customValidator()

	if err := validate.Struct(data); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errorResponses = append(errorResponses, fmt.Sprintf("%s is %s", err.Field(), err.Tag()))
		}
	}

	return errorResponses
}

// InterceptRequest
// Check and validate the payload, will intercept if the validation fails.
// If success, a locals will be set with given payload which can decrease unnecessary use of another
// BodyParser. Read more info about locals: https://docs.gofiber.io/api/ctx#locals
func InterceptRequest(data any) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		errorResponses := validateStruct(data)

		if errorResponses != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Given payload is invalid!",
				"errors":  errorResponses,
				"action":  "INVALID_PAYLOAD",
			})
		}

		c.Locals(BODY_DATA, data)

		return c.Next()
	}
}

// InterceptWebRequest
// Just like InterceptRequest, but instead of JSON, this function will intercept form based inputs.
func InterceptWebRequest(ref any) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sess := Session.Provide(c)

		if err := c.BodyParser(ref); err != nil {
			sess.SetSession("error", err.Error())
			return c.RedirectBack("/")
		}

		errorResponses := validateWebStruct(ref)

		if errorResponses != nil {
			sess.SetSession("validation_error", errorResponses)
			return c.RedirectBack("/")
		}

		c.Locals(BODY_DATA, ref)

		return c.Next()
	}
}

// GetPayload
// Smartly get a payload from locals and map them into given struct.
func GetPayload[T comparable](c *fiber.Ctx) T {
	payload := *c.Locals(BODY_DATA).(*T)

	return payload
}

// CheckMimes
// Check uploaded files mimes to avoid file injection.
func CheckMimes(reader io.Reader, acceptedTypes []string) error {
	mtype, err := mimetype.DetectReader(reader)

	if !slices.Contains(acceptedTypes, mtype.String()) || err != nil {
		return errors.New("ups, file not allowed... only allow: " + strings.Join(acceptedTypes, ", "))
	}

	return nil
}

// GetFailedValidation
// Get the validated validation from Session Store
func GetFailedValidation(store SessionStore) (any, any) {
	return store.GetFlash("error"), store.GetFlash("validation_error")
}
