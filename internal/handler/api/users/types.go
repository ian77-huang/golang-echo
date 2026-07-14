package users

type ApiUserHandler struct {
}

type RequestLogin struct {
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type RequestRegister struct {
	Account         string `json:"account" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}
