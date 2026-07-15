package frontend

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type FrontendParameter struct {
	DB   *gorm.DB
	Echo *echo.Echo
}
type FrontendHandler struct {
	DB *gorm.DB
}
