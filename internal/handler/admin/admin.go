package admin

func New(ap *AdminParameter) AdminHandler {
	h := &adminHandler{DB: ap.DB}

	// ap.Echo.GET("/", h.GetIndex)

	// user.New(&user.UserParameter{DB: ap.DB, Echo: ap.Echo})

	return h
}
