package users

import "github.com/labstack/echo/v5"

func New(e *echo.Group) {
	h := &ApiUserHandler{}

	users := e.Group("/user")
	users.POST("/register", h.PostRegister)
	users.POST("/login", h.PostLogin)
}
