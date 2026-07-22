package admin

func New(a *ApiAminParameter) {
	h := &ApiAminHandler{DB: a.DB}

	admin := a.EchoGroup.Group("/admin")

	admin.PUT("/user", h.PutUser)
}
