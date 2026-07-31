package frontend

import (
	"errors"
	"log"

	"github.com/ian77-huang/golang-echo/internal/handler/frontend/user"
	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
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
	auth := appAuth.GetUser[model.User](c)
	if auth == nil {
		return errors.New("auth is error")
	}

	hub := ws.LoadWebSocket(c)
	hub.Run()

	// time.Sleep(1 * time.Second)

	hub.New(c, "User/"+auth.ID, func(messageType int, p []byte) {
		log.Printf("\n === %+v %+v === \n", messageType, p)
		hub.Single(&ws.MessagePacket{ID: "User/1", Message: p})
	})

	if auth.ID == "2" {
		log.Printf("\n==== 123456 ======\n")
		// hub.Single(&ws.MessagePacket{ID: "User/1", Message: []byte("12323")})
	}

	return nil
}
