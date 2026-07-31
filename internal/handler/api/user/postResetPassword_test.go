package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appConfig "github.com/ian77-huang/golang-echo/internal/config"
	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	appValidator "github.com/ian77-huang/golang-echo/pkg/validator"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func resetAuth(t *testing.T) (*gorm.DB, *appAuth.Auth[model.User, model.Session]) {
	t.Helper()
	t.Setenv("SECRET_KEY", "test-secret")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Session{}); err != nil {
		t.Fatal(err)
	}
	hash, err := appAuth.HashPassword("old-password")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	if err := db.Create(&model.User{Id: 1, Account: "tester", Password: hash, IsActive: true, CreatedAt: &now, UpdatedAt: &now}).Error; err != nil {
		t.Fatal(err)
	}
	auth := appConfig.Auth(&appConfig.AuthParameter{
		UserService:    service.NewUserService(db),
		SessionService: service.NewSessionService(db),
	}, nil)
	return db, auth
}

func resetContext(e *echo.Echo, auth *appAuth.Auth[model.User, model.Session], rec *httptest.ResponseRecorder, body string) *echo.Context {
	c := e.NewContext(httptest.NewRequest(http.MethodPost, "/api/user/reset-password", bytes.NewBufferString(body)), rec)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "1"})
	c.Set(appAuth.CONTEXT_KEY_SESSION, &appAuth.Session[model.Session]{ID: "sess-1"})
	c.Set(appAuth.CONTEXT_KEY_AUTH, auth)
	return c
}

func TestPostResetPasswordSuccess(t *testing.T) {
	db, auth := resetAuth(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiUserHandler{DB: db}
	rec := httptest.NewRecorder()
	body := `{"oldPassword":"old-password","newPassword":"new-pass-123","confirmNewPassword":"new-pass-123"}`
	if err := h.PostResetPassword(resetContext(e, auth, rec, body)); err != nil || rec.Code != http.StatusOK {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}

func TestPostResetPasswordRejectsWrongOldPassword(t *testing.T) {
	db, auth := resetAuth(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiUserHandler{DB: db}
	rec := httptest.NewRecorder()
	body := `{"oldPassword":"wrong","newPassword":"new-pass-123","confirmNewPassword":"new-pass-123"}`
	if err := h.PostResetPassword(resetContext(e, auth, rec, body)); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}

func TestPostResetPasswordValidationAndContextErrors(t *testing.T) {
	db, auth := resetAuth(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiUserHandler{DB: db}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodPost, "/api/user/reset-password", bytes.NewBufferString(`{"oldPassword":`)), rec)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := h.PostResetPassword(c); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("bind: status=%d err=%v", rec.Code, err)
	}
	rec = httptest.NewRecorder()
	body := `{"oldPassword":"old-password","newPassword":"new-pass-123","confirmNewPassword":"different"}`
	if err := h.PostResetPassword(resetContext(e, auth, rec, body)); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("validate: status=%d err=%v", rec.Code, err)
	}
	rec = httptest.NewRecorder()
	validBody := `{"oldPassword":"old-password","newPassword":"new-pass-123","confirmNewPassword":"new-pass-123"}`
	c = e.NewContext(httptest.NewRequest(http.MethodPost, "/api/user/reset-password", bytes.NewBufferString(validBody)), rec)
	c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := h.PostResetPassword(c); err == nil {
		t.Fatal("expected missing user error")
	}
	c = resetContext(e, auth, rec, validBody)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "1"})
	c.Set(appAuth.CONTEXT_KEY_SESSION, &appAuth.Session[model.Session]{ID: "sess-1"})
	c.Set(appAuth.CONTEXT_KEY_AUTH, nil)
	if err := h.PostResetPassword(c); err == nil {
		t.Fatal("expected missing auth error")
	}
}
