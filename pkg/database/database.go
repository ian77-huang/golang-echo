package database

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("=== Error：Unable to connect to the database. ====")
	}
	return db
}

func GetDBConnect(c *echo.Context) (*gorm.DB, error) {
	value := c.Get(contextDBKey)
	if value == nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "database connection not found in context")
	}
	db, ok := value.(*gorm.DB)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "invalid database connection type")
	}

	return db, nil
}
