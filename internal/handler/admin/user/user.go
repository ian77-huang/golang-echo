package user

func New(ap *AdminUserParameter) {
	h := &adminUserHandler{DB: ap.DB}

	user := ap.EchoGroup.Group("/user")

	user.GET("", h.GetIndex)
	user.GET("/:page", h.GetIndex)
	user.GET("/:page/:pageSize", h.GetIndex)
}
