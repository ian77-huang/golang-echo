package api

import (
	"net/http"

	users "github.com/ian77-huang/golang-echo/internal/handler/api/user"
	"github.com/labstack/echo/v5"
)

func New(ap *ApiParameter) {
	h := &ApiHandler{DB: ap.DB}

	api := ap.Echo.Group("/api")
	api.POST("/lang", h.PostChangeLang)

	api.GET("/ping", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World12!"})
	})

	users.New(&users.ApiUserParameter{DB: ap.DB, EchoGroup: api})
}
