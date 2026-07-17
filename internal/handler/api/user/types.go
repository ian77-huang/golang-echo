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
type RequestPutProfile struct {
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required,numeric,len=10"`
	Email string `json:"email" validate:"required,email"`
	Bio   string `json:"bio" validate:"required"`
}
type RequestResetPassword struct {
	OldPassword        string `json:"oldPassword" validate:"required"`
	NewPassword        string `json:"newPassword" validate:"required,nefield=oldPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword" validate:"required,eqfield=NewPassword"`
}
