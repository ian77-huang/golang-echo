package database

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func Middleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set(contextDBKey, db)
			return next(c)
		}
	}
}
