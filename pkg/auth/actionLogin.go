package auth

import (
	"github.com/labstack/echo/v5"
)

func (a *Auth[TUser, TSession]) ActionLogin(c *echo.Context, account string, password string) (bool, error) {
	config, err := a.getConfig()
	if err != nil {
		return false, err
	}
	resolver := config.Resolver

	user, err := resolver.GetUserByAccount(account)
	if err != nil || user == nil || user.ID == "" {
		return false, NewError("error.auth.AccountLookupFailed", "error.auth - Account lookup Failed")
	}

	verify, err := VerifyPassword(password, user.Password)

	if err != nil {
		return false, NewError("error.auth.AccountLookupFailed", "error.auth - Account lookup Failed")
	}
	if !verify {
		return false, NewError("error.auth.AccountLookupFailed", "error.auth - Account lookup Failed")
	}

	status, err := a.createSession(c, user.ID)
	if err != nil {
		return false, NewError("error.auth.login.failed", "error.auth - Auth login failed")
	}

	return status, nil
}
