package user

func New(aup *ApiUserParameter) {
	h := &ApiUserHandler{DB: aup.DB}

	users := aup.EchoGroup.Group("/user")
	users.POST("/register", h.PostRegister)
	users.POST("/login", h.PostLogin)
	users.GET("/profile", h.GetProfile)
	users.PUT("/profile", h.PutProfile)
	users.POST("/profile/avatar", h.PostProfileUploadAvatar)
}
