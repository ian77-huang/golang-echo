package shared

import (
	"github.com/labstack/echo/v5"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTranslationHelpersFallbackWithoutLocalizer(t *testing.T) {
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	if T(c, "missing") != "localizer does not exist." {
		t.Fatal("unexpected translation fallback")
	}
	if TFactory(c)("missing") != "localizer does not exist." {
		t.Fatal("unexpected factory fallback")
	}
}
