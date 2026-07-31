package ws

import (
	"github.com/labstack/echo/v5"
)

func Middleware(hub WebSocketHub) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set(CONTEXT_WEB_SOCKET, hub)
			return next(c)
		}
	}
}

func LoadWebSocket(c *echo.Context) WebSocketHub {
	return c.Get(CONTEXT_WEB_SOCKET).(WebSocketHub)
}
