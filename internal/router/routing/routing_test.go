package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewRegistersAllRouteGroups(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	New(&RoutingParameter{DB: db, Echo: e})
	routes := e.Router().Routes()
	paths := make(map[string]bool)
	for _, route := range routes {
		paths[route.Method+" "+route.Path] = true
	}
	for _, expected := range []string{"GET /", "GET /ws", "GET /api/ping", "POST /api/lang", "GET /admin", "PUT /api/admin/user"} {
		if !paths[expected] {
			t.Fatalf("missing route %s; got %v", expected, paths)
		}
	}
}

func TestRoutingRegistersEndpoints(t *testing.T) {
	e := echo.New()
	h := &Routing{Echo: e}
	h.Frontend()
	h.Api()
	for _, tt := range []struct {
		method, path string
		want         int
	}{{http.MethodGet, "/api/ping", http.StatusOK}, {http.MethodPost, "/api/lang", http.StatusBadRequest}} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))
		if rec.Code != tt.want {
			t.Fatalf("%s %s: got %d", tt.method, tt.path, rec.Code)
		}
	}
}
