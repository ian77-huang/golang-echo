package database

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestMiddlewareStoresDatabase(t *testing.T) {
	db := NewSqlite(":memory:")
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	err := Middleware(db)(func(c *echo.Context) error {
		got, err := GetDBConnect(c)
		if err != nil || got != db {
			t.Fatalf("GetDBConnect() = %v, %v", got, err)
		}
		return nil
	})(c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDBConnectRequiresValidValue(t *testing.T) {
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	if _, err := GetDBConnect(c); err == nil {
		t.Fatal("expected missing database error")
	}
	c.Set(contextDBKey, "not-a-db")
	if _, err := GetDBConnect(c); err == nil {
		t.Fatal("expected invalid database error")
	}
}

func TestNewSqlitePanicsForUnusablePath(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected sqlite panic")
		}
	}()
	NewSqlite("/dev/null/app.db")
}
