package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func NewValidator() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func FormatValidationError(errs validator.ValidationErrors) map[string]string {
	errorMessages := make(map[string]string)
	for _, err := range errs {
		var message string
		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", err.Field())
		case "min":
			message = fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
		case "max":
			message = fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
		case "email":
			message = fmt.Sprintf("%s must be a valid email", err.Field())
		}

		errorMessages[err.Field()] = message
	}
	return errorMessages
}
