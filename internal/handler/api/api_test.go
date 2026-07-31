package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewRegistersAPIRoutes(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	New(&ApiParameter{DB: db, Echo: e})
	routes := e.Router().Routes()
	paths := make(map[string]bool)
	for _, route := range routes {
		paths[route.Method+" "+route.Path] = true
	}
	for _, expected := range []string{"GET /api/ping", "POST /api/lang", "PUT /api/admin/user"} {
		if !paths[expected] {
			t.Fatalf("missing route %s; got %v", expected, paths)
		}
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/ping", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("ping status=%d", rec.Code)
	}
}
