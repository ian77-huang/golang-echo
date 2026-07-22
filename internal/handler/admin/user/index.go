package user

import (
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/pkg/cast"
	"github.com/ian77-huang/golang-echo/pkg/utils"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
)

func (f *adminUserHandler) GetIndex(c *echo.Context) error {

	page, _ := cast.StringToInt(c.Param("page"), 1)
	pageSize, _ := cast.StringToInt(c.Param("pageSize"), 20)

	userService := service.NewUserService(f.DB)
	users, err := userService.GetPaginate(page, pageSize)

	totalRows, err := userService.CountUserAll()
	if err != nil {
		return response.ErrorInternalServerError(c, "ServerError", map[string]any{"No": 1})
	}

	var maxPage int = 0
	if totalRows > 0 {
		maxPage = (totalRows + pageSize - 1) / pageSize
	}

	pagination := utils.GetPagination(page, maxPage, 5)

	return response.Render(c, "admin:user/index.html", map[string]any{
		"data":       users,
		"total":      totalRows,
		"page":       page,
		"maxPage":    maxPage,
		"pagination": pagination,
	})
}
