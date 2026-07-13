package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterRegistersUtilityAndPingRoutes(t *testing.T) {
	e := New()
	for _, tt := range []struct {
		path string
		want int
	}{{"/api/ping", http.StatusOK}, {"/.well-known/appspecific/com.chrome.devtools.json", http.StatusNotFound}} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))
		if rec.Code != tt.want {
			t.Fatalf("%s: got %d", tt.path, rec.Code)
		}
	}
}
