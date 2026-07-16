package user

func New(ap *UserParameter) {
	h := &UserHandler{DB: ap.DB}

	users := ap.Echo.Group("/user")

	users.GET("", h.GetIndex)
	users.GET("/login", h.GetLogin)
	users.GET("/register", h.GetRegister)
	users.GET("/logout", h.GetLogout)
	users.GET("/profile", h.GetProfile)
	users.GET("/reset-password", h.GetResetPassword)
}
