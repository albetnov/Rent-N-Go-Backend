package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	FailedFields string
	Tag          string
	Value        string
}

var validate = validator.New()

const BODY_DATA = "body_data"

func validateStruct(data any) []*ErrorResponse {
	var errors []*ErrorResponse

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedFields = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}

func InterceptRequest(data any) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		errors := validateStruct(data)

		if errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Given payload is invalid!",
				"errors":  errors,
			})
		}

		c.Locals(BODY_DATA, data)

		return c.Next()
	}
}

func GetPayload[T comparable](c *fiber.Ctx) T {
	payload := *c.Locals(BODY_DATA).(*T)

	return payload
}
