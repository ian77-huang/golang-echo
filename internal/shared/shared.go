package shared

import (
	appi18n "github.com/ian77-huang/golang-echo/pkg/i18n"
	"github.com/labstack/echo/v5"
)

func T(c *echo.Context, messageID string, pairs ...any) string {
	return appi18n.T(c, messageID, pairs...)
}

func TFactory(c *echo.Context) func(messageID string, pairs ...any) string {
	return func(messageID string, pairs ...any) string {
		return appi18n.T(c, messageID, pairs...)
	}
}
