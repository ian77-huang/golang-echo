package api

import (
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type ApiParameter struct {
	DB   *gorm.DB
	Echo *echo.Echo
}
type ApiHandler struct {
	DB *gorm.DB
}

type ChangeLangRequest struct {
	Code string `json:"code" validate:"required,oneof=zh-TW en"`
}
