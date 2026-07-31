package store

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestMiddlewareStoresAndLoadsServer(t *testing.T) {
	ss := New(nil)
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

	called := false
	err := Middleware(ss)(func(c *echo.Context) error {
		called = LoadStore(c) == ss
		return nil
	})(c)

	if err != nil || !called {
		t.Fatalf("middleware called=%v err=%v", called, err)
	}
}
