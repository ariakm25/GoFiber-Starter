package utils

import (
	database "GoFiber-API/external/database/postgres"
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

	// Custom validation for unique fields on the database.
	_ = validate.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
		// value format is table.column
		value := strings.Split(fl.Param(), ".")
		table := value[0]
		column := value[1]
		except := ""

		if len(value) == 3 {
			except = value[2]
		}

		var count int64

		if except == "" {
			_ = database.Connection.Table(table).Where(column+" = ?", fl.Field().String()).Count(&count)
			return count == 0
		} else {
			_ = database.Connection.Table(table).Where(column+" = ? AND "+column+" != ?", fl.Field().String(), except).Count(&count)
			return count == 0
		}
	})

	// Custom validation for check existing field.
	_ = validate.RegisterValidation("exist", func(fl validator.FieldLevel) bool {
		// value format is table.column
		value := strings.Split(fl.Param(), ".")
		table := value[0]
		column := value[1]

		var count int64

		_ = database.Connection.Table(table).Where(column+" = ?", fl.Field().String()).Count(&count)

		return count > 0
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
	case "unique":
		return "The " + fe.Field() + " is already taken."
	case "exist":
		return "The " + fe.Field() + " is not exist."
	default:
		return "The " + fe.Field() + " field is invalid."
	}
}
