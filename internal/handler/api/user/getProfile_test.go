package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetProfileRequiresAuth(t *testing.T) {
	e := echo.New()
	h := &ApiUserHandler{}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/api/user/profile", nil), rec)
	if err := h.GetProfile(c); err != nil || rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}

func TestGetProfileReturnsEmptyProfileWhenMissing(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.UserProfile{}); err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	h := &ApiUserHandler{DB: db}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/api/user/profile", nil), rec)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "9"})
	if err := h.GetProfile(c); err != nil || rec.Code != http.StatusOK {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}

func TestGetProfileReturnsStoredProfile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.UserProfile{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&model.UserProfile{UserID: 1, Name: "Alice", Email: "a@example.com"}).Error; err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	h := &ApiUserHandler{DB: db}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/api/user/profile", nil), rec)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "1"})
	if err := h.GetProfile(c); err != nil || rec.Code != http.StatusOK {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}

func TestGetProfilePropagatesDatabaseError(t *testing.T) {
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
	h := &ApiUserHandler{DB: db}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/api/user/profile", nil), rec)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "1"})
	if err := h.GetProfile(c); err == nil {
		t.Fatal("expected error, got nil")
	}
}
