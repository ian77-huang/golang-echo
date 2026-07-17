package auth

import (
	"github.com/labstack/echo/v5"
)

func (a *Auth[TUser, TSession]) ActionResetPassword(c *echo.Context, userId string, oldPassword string, newPassword string) (*User[TUser], error) {
	config, err := a.getConfig()
	if err != nil {
		return nil, err
	}
	resolver := config.Resolver

	user, err := resolver.GetUser(userId)
	if err != nil || user == nil || user.ID == "" {
		return nil, NewError("error.auth.AccountLookupFailed", "error.auth - Account lookup Failed")
	}

	verify, err := VerifyPassword(oldPassword, user.Password)
	if err != nil {
		return nil, NewError("error.auth.AccountLookupFailed", "error.auth - Account lookup Failed")
	}
	if !verify {
		return nil, NewError("error.auth.AccountLookupFailed", "error.auth - Account lookup Failed")
	}

	passwordHash, err := HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	updateUser, err := a.UpdateUserPassword(user.ID, passwordHash)
	if err != nil {
		return nil, NewError("error.auth.login.failed", "error.auth - Auth login failed")
	}

	_, err = resolver.DeleteSession(user.ID)
	if err != nil {
		return nil, err
	}

	a.deleteAccessToken(c)

	return updateUser, nil
}
