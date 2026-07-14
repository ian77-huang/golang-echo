package routing

import (
	apiHandler "github.com/ian77-huang/golang-echo/internal/handler/api"
	"github.com/labstack/echo/v5"
)

func (h *Routing) Api(e *echo.Echo) {
	apiHandler.New(e)
}
