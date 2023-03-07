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
)

type ErrorResponse struct {
	FailedFields string
	Tag          string
	Value        string
}

var validate = validator.New()

const BODY_DATA = "body_data"

// validateStruct
// Validate an payload to a given struct, return error if something wrong, and return empty if all passed.
func validateStruct(data any) []*ErrorResponse {
	var errorResponses []*ErrorResponse

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
