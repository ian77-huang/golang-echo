package routing

import (
	"net/http"

	apiHandler "github.com/ian77-huang/golang-echo/internal/handler/api"
	apiUsersHandler "github.com/ian77-huang/golang-echo/internal/handler/api/users"
	"github.com/labstack/echo/v5"
)

func ApiRouting(e *echo.Echo) {
	apiGroup := e.Group("/api")
	apiGroup.POST("/lang", apiHandler.PostChangeLang)

	apiGroup.GET("/ping", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World12!"})
	})

	apiUsersGroup := apiGroup.Group("/users")
	apiUsersGroup.POST("/register", apiUsersHandler.PostRegister)
}
