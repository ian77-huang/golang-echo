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

	// user := GetUser[TUser](c)
	sess := GetSession[TSession](c)
	sesss, err := resolver.DeleteSession(sess.ID)
	if sesss == nil && err == nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	a.deleteAccessToken(c)

	return true, nil
}
