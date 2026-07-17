package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

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
