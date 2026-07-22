package admin

import (
	adminUser "github.com/ian77-huang/golang-echo/internal/handler/admin/user"
)

func New(ap *AdminParameter) AdminHandler {
	h := &adminHandler{DB: ap.DB}

	admin := ap.Echo.Group("/admin")

	admin.GET("", h.GetIndex)

	adminUser.New(&adminUser.AdminUserParameter{
		DB:        ap.DB,
		EchoGroup: admin,
	})

	return h
}
