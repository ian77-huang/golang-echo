package admin

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type AdminParameter struct {
	DB   *gorm.DB
	Echo *echo.Echo
}
type AdminHandler interface {
}
type adminHandler struct {
	DB *gorm.DB
}
