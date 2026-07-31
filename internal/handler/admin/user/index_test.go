package user

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/pkg/utils"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type captureRenderer struct {
	name string
	data map[string]any
}

func (r *captureRenderer) Render(_ *echo.Context, _ io.Writer, name string, data any) error {
	r.name = name
	r.data, _ = data.(map[string]any)
	return nil
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestGetIndexRendersUserList(t *testing.T) {
	db := newTestDB(t)
	now := time.Now()
	for i := 1; i <= 3; i++ {
		if err := db.Create(&model.User{Account: "user" + string(rune('0'+i)), CreatedAt: &now, UpdatedAt: &now}).Error; err != nil {
			t.Fatal(err)
		}
	}

	e := echo.New()
	r := &captureRenderer{}
	e.Renderer = r
	h := &adminUserHandler{DB: db}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/user/1/10", nil), rec)
	c.SetPathValues(echo.PathValues{
		{Name: "page", Value: "1"},
		{Name: "pageSize", Value: "10"},
	})

	if err := h.GetIndex(c); err != nil {
		t.Fatal(err)
	}
	if r.name != "admin:user/index.html" || rec.Code != http.StatusOK {
		t.Fatalf("name=%q status=%d", r.name, rec.Code)
	}
	if r.data["total"] != 3 {
		t.Fatalf("total=%#v", r.data["total"])
	}
	if r.data["page"] != 1 || r.data["maxPage"] != 1 {
		t.Fatalf("page=%#v maxPage=%#v", r.data["page"], r.data["maxPage"])
	}
	pagination, ok := r.data["pagination"].(*utils.Pagination)
	if !ok || pagination.CurrentPage != 1 || pagination.TotalPages != 1 {
		t.Fatalf("pagination=%#v", r.data["pagination"])
	}
}

func TestGetIndexPaginates(t *testing.T) {
	db := newTestDB(t)
	now := time.Now()
	for i := 1; i <= 25; i++ {
		account := "user_" + string(rune('a'+i))
		if err := db.Create(&model.User{Account: account, CreatedAt: &now, UpdatedAt: &now}).Error; err != nil {
			t.Fatal(err)
		}
	}

	e := echo.New()
	r := &captureRenderer{}
	e.Renderer = r
	h := &adminUserHandler{DB: db}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/user/2/10", nil), rec)
	c.SetPathValues(echo.PathValues{
		{Name: "page", Value: "2"},
		{Name: "pageSize", Value: "10"},
	})

	if err := h.GetIndex(c); err != nil {
		t.Fatal(err)
	}
	if r.data["total"] != 25 || r.data["maxPage"] != 3 {
		t.Fatalf("total=%#v maxPage=%#v", r.data["total"], r.data["maxPage"])
	}
	if len(r.data["data"].([]service.UserOmitPassword)) != 10 {
		t.Fatalf("expected 10 users on page 2, got %d", len(r.data["data"].([]service.UserOmitPassword)))
	}
}

func TestGetIndexUsesDefaultsWithoutParams(t *testing.T) {
	db := newTestDB(t)
	e := echo.New()
	r := &captureRenderer{}
	e.Renderer = r
	h := &adminUserHandler{DB: db}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/user", nil), rec)

	if err := h.GetIndex(c); err != nil {
		t.Fatal(err)
	}
	if r.data["page"] != 1 || r.data["maxPage"] != 0 {
		t.Fatalf("page=%#v maxPage=%#v", r.data["page"], r.data["maxPage"])
	}
}

func TestGetIndexReturnsServerErrorOnDatabaseFailure(t *testing.T) {
	db := newTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	r := &captureRenderer{}
	e.Renderer = r
	h := &adminUserHandler{DB: db}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/user", nil), rec)

	err = h.GetIndex(c)
	he, ok := err.(*echo.HTTPError)
	if !ok || he.Code != http.StatusInternalServerError {
		t.Fatalf("expected HTTPError 500, got %#v", err)
	}
}
