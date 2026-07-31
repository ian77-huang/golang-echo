package admin

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func testEcho() (*echo.Echo, *captureRenderer) {
	e := echo.New()
	r := &captureRenderer{}
	e.Renderer = r
	return e, r
}

func TestGetIndexRendersTemplate(t *testing.T) {
	e, r := testEcho()
	h := &adminHandler{}
	rec := httptest.NewRecorder()
	if err := h.GetIndex(e.NewContext(httptest.NewRequest(http.MethodGet, "/admin", nil), rec)); err != nil {
		t.Fatal(err)
	}
	if r.name != "admin:index/index.html" || rec.Code != http.StatusOK {
		t.Fatalf("name=%q status=%d", r.name, rec.Code)
	}
}

func TestNewRegistersRoutes(t *testing.T) {
	e, r := testEcho()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	New(&AdminParameter{DB: db, Echo: e})

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /admin status=%d", rec.Code)
	}
	if r.name != "admin:index/index.html" {
		t.Fatalf("unexpected renderer name %q", r.name)
	}
}
