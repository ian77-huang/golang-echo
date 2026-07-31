package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	appValidator "github.com/ian77-huang/golang-echo/pkg/validator"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPutProfileUpdatesProfile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.UserProfile{}); err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiUserHandler{DB: db}
	rec := httptest.NewRecorder()
	body := `{"name":"Alice","phone":"0912345678","email":"alice@example.com","bio":"hello"}`
	c := e.NewContext(httptest.NewRequest(http.MethodPut, "/api/user/profile", bytes.NewBufferString(body)), rec)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "1"})
	if err := h.PutProfile(c); err != nil || rec.Code != http.StatusOK {
		t.Fatalf("status=%d err=%v body=%s", rec.Code, err, rec.Body.String())
	}
}

func TestPutProfileValidationAndBindErrors(t *testing.T) {
	e := echo.New()
	h := &ApiUserHandler{}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodPut, "/api/user/profile", bytes.NewBufferString(`{"name":`)), rec)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "1"})
	if err := h.PutProfile(c); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("bind: status=%d err=%v", rec.Code, err)
	}
	rec = httptest.NewRecorder()
	body := `{"name":"Alice","phone":"abc","email":"bad","bio":"x"}`
	c = e.NewContext(httptest.NewRequest(http.MethodPut, "/api/user/profile", bytes.NewBufferString(body)), rec)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "1"})
	if err := h.PutProfile(c); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("validate: status=%d err=%v", rec.Code, err)
	}
}

func TestPutProfilePropagatesDatabaseError(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiUserHandler{DB: db}
	rec := httptest.NewRecorder()
	body := `{"name":"Alice","phone":"0912345678","email":"alice@example.com","bio":"hello"}`
	c := e.NewContext(httptest.NewRequest(http.MethodPut, "/api/user/profile", bytes.NewBufferString(body)), rec)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "1"})
	if err := h.PutProfile(c); err == nil {
		t.Fatal("expected error, got nil")
	}
}
