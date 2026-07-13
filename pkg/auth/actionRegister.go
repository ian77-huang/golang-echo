package auth

import (
	"github.com/labstack/echo/v5"
)

func (a *Auth[TUser, TSession]) ActionRegister(c *echo.Context, account string, password string) (string, error) {
	exists, err := a.IsAccountExist(account)
	if err != nil {
		return "", err
	}
	if exists {
		return "", NewError("error.auth.AccountAlreadyExists", "error.auth - Account already exists")
	}

	user, err := a.CreateUser(account, password)
	if err != nil {
		return "", err
	}

	if _, err := a.createSession(c, user.ID); err != nil {
		return "", err
	}

	return user.ID, nil
}
