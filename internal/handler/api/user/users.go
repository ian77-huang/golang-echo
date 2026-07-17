package user

func New(a *ApiUserParameter) {
	h := &ApiUserHandler{DB: a.DB}

	users := a.EchoGroup.Group("/user")
	users.POST("/register", h.PostRegister)
	users.POST("/login", h.PostLogin)
	users.GET("/profile", h.GetProfile)
	users.PUT("/profile", h.PutProfile)
	users.POST("/profile/avatar", h.PostProfileUploadAvatar)
	users.POST("/reset-password", h.PostResetPassword)
}
