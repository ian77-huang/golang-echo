package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func Redirect(c *echo.Context, url string) error {
	return c.Redirect(http.StatusFound, url)
}

func RedirectAfterPost(c *echo.Context, url string) error {
	return c.Redirect(http.StatusSeeOther, url)
}

func RedirectPermanent(c *echo.Context, url string) error {
	return c.Redirect(http.StatusMovedPermanently, url)
}
