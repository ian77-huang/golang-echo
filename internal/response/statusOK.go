package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func JsonOk(c *echo.Context, i any) error {
	return c.JSON(http.StatusOK, i)
}

func JsonCreated(c *echo.Context, i any) error {
	return c.JSON(http.StatusCreated, i)
}

func Render(c *echo.Context, name string, data any) (err error) {
	return c.Render(http.StatusOK, name, data)
}
