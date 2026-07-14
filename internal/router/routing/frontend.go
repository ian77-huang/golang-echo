package routing

import (
	"github.com/ian77-huang/golang-echo/internal/handler"
	handlerUsers "github.com/ian77-huang/golang-echo/internal/handler/frontend/users"
	"github.com/labstack/echo/v5"
)

func (h *Routing) Frontend(e *echo.Echo) {
	e.GET("/", handler.GetIndex)

	handlerUsers.New(e)
}
