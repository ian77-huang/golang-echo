package admin

import (
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
)

func (f *adminHandler) GetIndex(c *echo.Context) error {

	bibleService := service.NewBibleService(f.DB)
	bible, err := bibleService.GetBibleByDate()

	if err != nil {
		return response.ErrorInternalServerError(c, "invalid_request")
	}

	return response.Render(c, "admin:index/index.html", map[string]any{
		"bible": bible,
	})
}
