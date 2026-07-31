package frontend

import (
	"github.com/ian77-huang/golang-echo/internal/handler/frontend/user"
	"github.com/ian77-huang/golang-echo/pkg/ws"
	"github.com/labstack/echo/v5"
)

func New(ap *FrontendParameter) {
	h := &FrontendHandler{DB: ap.DB}

	ap.Echo.GET("/", h.GetIndex)

	ap.Echo.GET("/ws", wsHeadler)

	user.New(&user.UserParameter{DB: ap.DB, Echo: ap.Echo})
}

func wsHeadler(c *echo.Context) error {
	hub := ws.LoadWebSocket(c)
	hub.Run()

	// time.Sleep(1 * time.Second)

	hub.New(c)

	return nil
}
