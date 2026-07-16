package response

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	appi18n "github.com/ian77-huang/golang-echo/pkg/i18n"
	"github.com/labstack/echo/v5"
)

type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type FieldError struct {
	Field  string
	Tag    string
	Params []interface{}
}

func (fe FieldError) Error() string {
	return fmt.Sprintf("field: %s, tag: %s, params: %v", fe.Field, fe.Tag, fe.Params)
}

func NewFieldError(field, tag string, params ...interface{}) FieldError {
	return FieldError{
		Field:  field,
		Tag:    tag,
		Params: params,
	}
}

func ErrorInternalServerError(c *echo.Context, code string) error {
	return echo.NewHTTPError(http.StatusInternalServerError, appi18n.T(c, "error."+code))
}
func ErrorBadRequest(c *echo.Context, code string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{
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
func ValidationCustomError(c *echo.Context, err FieldError) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    "validation_failed",
		Message: appi18n.T(c, "error.validation_failed"),
		Errors:  customValidationMessages(c, err),
	})
}
func ValidationErrorAuth(c *echo.Context, err error) error {
	if fe, ok := err.(appAuth.FieldError); ok {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "validation_failed",
			Message: appi18n.T(c, fe.Tag),
			Errors:  customValidationMessages(c, FieldError{Field: "", Tag: fe.Tag, Params: fe.Params}),
		})

	}
	return errors.New("auth error: Type assertion failed")
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
		messages[fieldErr.Field()] = translateFieldError(c, NewFieldError(fieldErr.Field(), fieldErr.Tag(), fieldErr.Param()))
	}

	return messages
}
func customValidationMessages(c *echo.Context, fe FieldError) map[string]string {
	messages := make(map[string]string)

	messages[fe.Field] = translateFieldError(c, fe)

	return messages
}

func translateFieldError(c *echo.Context, fieldErr FieldError) string {
	var param interface{}
	if len(fieldErr.Params) > 0 {
		param = fieldErr.Params[0]
	}
	errorMessage := ""
	switch fieldErr.Tag {
	case "required":
		errorMessage = appi18n.T(c, "validation.required", appi18n.KV("Field", fieldErr.Field))
	case "oneof":
		errorMessage = appi18n.T(c, "validation.oneof", appi18n.KV("Field", fieldErr.Field), appi18n.KV("Choices", param))
	case "min":
		errorMessage = appi18n.T(c, "validation.min", appi18n.KV("Field", fieldErr.Field), appi18n.KV("Length", param))
	case "eqfield":
		errorMessage = appi18n.T(c, "validation.eqfield", appi18n.KV("Field", fieldErr.Field), appi18n.KV("Target", param))
	case "invalid":
		errorMessage = appi18n.T(c, "validation.invalid", appi18n.KV("Field", fieldErr.Field))
	default:
		errorMessage = appi18n.T(c, fieldErr.Tag)
	}
	return errorMessage
}
