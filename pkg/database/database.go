package database

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

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

func New(cfg *DBConfig) (*gorm.DB, error) {
	switch cfg.Driver {
	case Sqlite:
		if cfg.Sqlite == nil {
			return nil, fmt.Errorf("Sqlite configuration is not provided")
		}
		return NewSqlite(cfg.Sqlite)
	case PostgreSQL:
		if cfg.Pgsql == nil {
			return nil, fmt.Errorf("Pgsql configuration is not provided")
		}
		return NewPostgresSql(cfg.Pgsql)
	case Mysql:
		if cfg.Mysql == nil {
			return nil, fmt.Errorf("Mysql configuration is not provided")
		}
		return NewMysql(cfg.Mysql)

	default:
		return nil, fmt.Errorf("unsupported driver")
	}
}
