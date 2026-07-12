package auth

import (
	"errors"

	"github.com/ian77-huang/golang-echo/pkg/argon2"
	"github.com/labstack/echo/v5"
)

func (a *Auth[TUser, TSession]) ActionRegister(c *echo.Context, account string, password string) (string, error) {
	config, err := a.getConfig()
	if err != nil {
		return "", err
	}
	resolver := config.Resolver
	if resolver.IsAccountExist == nil {
		return "", NewError("error.auth.InvalidConfiguration", "account lookup resolver is not configured")
	}

	exists, err := resolver.IsAccountExist(account)
	if err != nil {
		if fieldErr, ok := errors.AsType[FieldError](err); ok {
			return "", fieldErr
		}
		return "", NewError("error.auth.AccountLookupFailed", "failed to look up account")
	}
	if exists {
		return "", NewError("error.auth.AccountAlreadyExists", "account already exists")
	}

	if resolver.CreateUser == nil {
		return "", NewError("error.auth.InvalidConfiguration", "user creation resolver is not configured")
	}
	passwordHash, err := argon2.HashPassword(password)
	if err != nil {
		return "", NewError("error.auth.FailedToSecurePassword", "failed to secure password")
	}
	user, err := resolver.CreateUser(account, passwordHash)
	if err != nil || user == nil || user.ID == "" {
		return "", NewError("error.auth.CreateUserError", "create user error")
	}
	if _, err := a.createSession(c, user.ID); err != nil {
		return "", err
	}

	return user.ID, nil
}
