package user

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type ApiUserParameter struct {
	DB        *gorm.DB
	EchoGroup *echo.Group
}
type ApiUserHandler struct {
	DB *gorm.DB
}

type RequestLogin struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type RequestRegister struct {
	Account         string `json:"account" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}
