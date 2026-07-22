package admin

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type ApiAminParameter struct {
	DB        *gorm.DB
	EchoGroup *echo.Group
}

type ApiAminHandler struct {
	DB *gorm.DB
}

type RequestPutUser struct {
	Id       string `json:"id" validate:"required"`
	Account  string `json:"account" validate:"required"`
	IsActive *bool  `json:"isActive" validate:"required"`
	IsAdmin  *bool  `json:"isAdmin" validate:"required"`
}
