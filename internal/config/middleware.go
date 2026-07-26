package config

import (
	"github.com/labstack/echo/v5"
)

const CONTEXT_SERVER_NAME = "contextServerName"

type ConfigMiddleware struct {
	ServerName string
}

func Middleware(cm *ConfigMiddleware) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if cm.ServerName == "" {
				cm.ServerName = "echo"
			}
			c.Set(CONTEXT_SERVER_NAME, cm.ServerName)
			return next(c)
		}
	}
}

func LoadServerName(c *echo.Context) string {
	return c.Get(CONTEXT_SERVER_NAME).(string)
}
