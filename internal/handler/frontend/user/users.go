package users

import "github.com/labstack/echo/v5"

func New(e *echo.Echo) {
	h := &UserHandler{}

	users := e.Group("/user")

	users.GET("", h.GetIndex)
	users.GET("/login", h.GetLogin)
	users.GET("/register", h.GetRegister)
	users.GET("/logout", h.GetLogout)
}
