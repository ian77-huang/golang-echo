package store

import (
	"github.com/labstack/echo/v5"
)

const CONTEXT_STORE = "contextStore"

func Middleware(server *StoreServer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set(CONTEXT_STORE, server)
			return next(c)
		}
	}
}

func LoadStore(c *echo.Context) *StoreServer {
	return c.Get(CONTEXT_STORE).(*StoreServer)
}
