package response

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	appi18n "github.com/ian77-huang/golang-echo/pkg/i18n"
	"github.com/labstack/echo/v5"
)

type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func Error(c *echo.Context, status int, code string) error {
	return c.JSON(status, ErrorResponse{
		Code:    code,
		Message: appi18n.T(c, "error."+code),
	})
}

func ValidationError(c *echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    "validation_failed",
		Message: appi18n.T(c, "error.validation_failed"),
		Errors:  validationMessages(c, err),
	})
}

func validationMessages(c *echo.Context, err error) map[string]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return map[string]string{
			"_": appi18n.T(c, "error.invalid_request"),
		}
	}

	messages := make(map[string]string)

	for _, fieldErr := range validationErrors {
		field := fieldErr.Field()
		jsonField := fieldErr.Field()

		switch fieldErr.Tag() {
		case "required":
			messages[jsonField] = appi18n.T(c, "validation.required", appi18n.KV("Field", field))
		case "oneof":
			messages[jsonField] = appi18n.T(c, "validation.oneof", appi18n.KV("Field", field), appi18n.KV("Choices", fieldErr.Param()))
		default:
			messages[jsonField] = appi18n.T(c, "validation.invalid", appi18n.KV("Field", field))
		}
	}

	return messages
}
