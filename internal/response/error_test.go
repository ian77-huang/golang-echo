package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	appValidator "github.com/ian77-huang/golang-echo/pkg/validator"
	"github.com/labstack/echo/v5"
)

func testContext() (*echo.Echo, *echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Validator = appValidator.New()
	rec := httptest.NewRecorder()
	return e, e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), rec), rec
}

func TestResponseHelpers(t *testing.T) {
	_, c, rec := testContext()
	if err := JsonOk(c, map[string]string{"ok": "yes"}); err != nil || rec.Code != http.StatusOK {
		t.Fatalf("JSON: status=%d err=%v", rec.Code, err)
	}
	_, c, rec = testContext()
	if err := ErrorBadRequest(c, "forbidden"); err != nil || rec.Code != http.StatusForbidden {
		t.Fatalf("Error: status=%d err=%v", rec.Code, err)
	}
}

func TestValidationResponses(t *testing.T) {
	_, c, rec := testContext()
	err := c.Validate(struct {
		Name string `validate:"required"`
	}{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if err := ValidationError(c, err); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("ValidationError: %v status=%d", err, rec.Code)
	}
	_, c, rec = testContext()
	if err := ValidationCustomError(c, NewFieldError("name", "min", 3)); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("ValidationCustomError: %v", err)
	}
	_, c, rec = testContext()
	if err := ValidationErrorAuth(c, appAuth.NewError("error.auth.failed", "failed")); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("ValidationErrorAuth: %v", err)
	}
	if err := ValidationErrorAuth(c, errors.New("plain")); err == nil {
		t.Fatal("expected type assertion error")
	}
}

func TestFieldErrorAndMessageFallbacks(t *testing.T) {
	_, c, _ := testContext()
	if NewFieldError("name", "required").Error() == "" {
		t.Fatal("expected field error text")
	}
	if validationMessages(c, errors.New("bad"))["_"] == "" {
		t.Fatal("expected fallback message")
	}
	for _, tag := range []string{"required", "oneof", "min", "eqfield", "invalid", "custom"} {
		if translateFieldError(c, NewFieldError("name", tag, "x")) == "" {
			t.Fatalf("missing message for %s", tag)
		}
	}
}
