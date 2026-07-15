package routing

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Routing struct {
	DB   *gorm.DB
	Echo *echo.Echo
}
type RoutingParameter struct {
	DB   *gorm.DB
	Echo *echo.Echo
}
