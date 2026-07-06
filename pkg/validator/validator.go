package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func (cv *CustomValidator) Messages(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return map[string]string{
			"_": err.Error(),
		}
	}

	messages := make(map[string]string)

	for _, fieldErr := range validationErrors {
		field := fieldErr.Field()

		switch fieldErr.Tag() {
		case "required":
			messages[field] = field + " is required"
		case "oneof":
			messages[field] = field + " must be one of: " + fieldErr.Param()
		default:
			messages[field] = field + " is invalid"
		}
	}

	return messages
}
