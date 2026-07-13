package auth

import (
	"github.com/labstack/echo/v5"
)

func (a *Auth[TUser, TSession]) ActionLogout(c *echo.Context) (bool, error) {
	config, err := a.getConfig()
	if err != nil {
		return false, err
	}
	resolver := config.Resolver

	user := GetUser[TUser](c)

	_, err = resolver.DeleteSession(user.ID)
	if err != nil {
		return false, err
	}

	a.deleteAccessToken(c)

	return true, nil
}
