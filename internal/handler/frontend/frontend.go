package frontend

import (
	"github.com/ian77-huang/golang-echo/internal/handler/frontend/user"
)

func New(ap *FrontendParameter) {
	h := &FrontendHandler{DB: ap.DB}

	ap.Echo.GET("/", h.GetIndex)

	user.New(&user.UserParameter{DB: ap.DB, Echo: ap.Echo})
}
