package routing

import "github.com/labstack/echo/v5"

func New(e *echo.Echo) {
	h := &Routing{}

	h.Frontend(e)
	h.Api(e)
}
