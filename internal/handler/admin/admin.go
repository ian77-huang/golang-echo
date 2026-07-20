package admin

func New(ap *AdminParameter) AdminHandler {
	h := &adminHandler{DB: ap.DB}

	admin := ap.Echo.Group("/admin")

	admin.GET("", h.GetIndex)

	// user.New(&user.UserParameter{DB: ap.DB, Echo: ap.Echo})

	return h
}
