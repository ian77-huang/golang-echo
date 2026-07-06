package router

import (
	"net/http"

	"github.com/ian77-huang/golang-echo/internal/router/routing"
	"github.com/labstack/echo/v5"
)

func New() *echo.Echo {
	e := echo.New()

	routing.FrontendRouting(e)
	routing.ApiRouting(e)

	e.GET("/.well-known/appspecific/com.chrome.devtools.json", func(c *echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})

	return e
}
