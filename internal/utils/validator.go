package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func NewValidator() *validator.Validate {
	validate := validator.New()

	// Custom validation for uuid.UUID fields.
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if _, err := uuid.Parse(field); err != nil {
			return true
		}
		return false
	})

	return validate
}

func ValidatorErrors(err error) map[string]string {
	fields := map[string]string{}

	for _, err := range err.(validator.ValidationErrors) {
		// lowercase the first letter and remove the struct name from the field

		fields[strings.ToLower(err.Field())] = customMessage(err)
	}

	return fields
}

func customMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "The " + fe.Field() + " field is required."
	case "email":
		return "The " + fe.Field() + " field must be a valid email address."
	case "uuid":
		return "The " + fe.Field() + " field must be a valid UUID."
	case "min":
		return "The " + fe.Field() + " field must be at least " + fe.Param() + " characters."
	case "max":
		return "The " + fe.Field() + " field must be at most " + fe.Param() + " characters."
	case "eqfield":
		return "The " + fe.Field() + " field must be equal to the " + fe.Param() + " field."
	default:
		return "The " + fe.Field() + " field is invalid."
	}
}
