package user

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type AdminUserParameter struct {
	DB        *gorm.DB
	EchoGroup *echo.Group
}
type adminUserHandler struct {
	DB *gorm.DB
}
