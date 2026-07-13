package auth

import (
	"errors"

	"github.com/labstack/echo/v5"
)

func IsSignedIn[TUser any](c *echo.Context) bool {
	return GetUser[TUser](c) != nil
}
func (a *Auth[TUser, TSession]) IsAccountExist(account string) (bool, error) {
	config, err := a.getConfig()
	if err != nil {
		return false, err
	}
	resolver := config.Resolver

	if resolver.IsAccountExist == nil {
		return false, NewError("error.auth.InvalidConfiguration", "account lookup resolver is not configured")
	}

	exists, err := resolver.IsAccountExist(account)
	if err != nil {
		if fieldErr, ok := errors.AsType[FieldError](err); ok {
			return false, fieldErr
		}
		return false, NewError("error.auth.AccountLookupFailed", "failed to look up account")
	}
	return exists, nil
}
func checkRoute[TUser any, TSession any](config *Config[TUser, TSession]) {
	if config.Route == nil {
		config.Route = &Route[TUser]{}
	}
	if config.Route.AuthOnly == nil {
		config.Route.AuthOnly = &RoutesPaths{}
	}
	if config.Route.GuestOnly == nil {
		config.Route.GuestOnly = &RoutesPaths{}
	}
	if config.Route.AuthOnly.Rules == nil {
		config.Route.AuthOnly.Rules = []string{}
	}
	if config.Route.GuestOnly.Rules == nil {
		config.Route.GuestOnly.Rules = []string{}
	}
	if config.Route.SpecialValidate == nil {
		config.Route.SpecialValidate = func(c *echo.Context, validateRule *ValidateRule[TUser]) {}
	}
}
func checkResolver[TUser any, TSession any](config *Config[TUser, TSession]) error {
	if config.Resolver == nil {
		return errors.New("Auth - config.Resolver property not set")
	}
	if config.Resolver.CreateSession == nil {
		return errors.New("Auth - config.Resolver.CreateSession property not set")
	}
	if config.Resolver.CreateUser == nil {
		return errors.New("Auth - config.Resolver.CreateUser property not set")
	}
	if config.Resolver.DeleteSession == nil {
		return errors.New("Auth - config.Resolver.DeleteSession property not set")
	}
	if config.Resolver.GetSession == nil {
		return errors.New("Auth - config.Resolver.GetSession property not set")
	}
	if config.Resolver.GetUser == nil {
		return errors.New("Auth - config.Resolver.GetUser property not set")
	}
	if config.Resolver.GetUserByAccount == nil {
		return errors.New("Auth - config.Resolver.GetUserByAccount property not set")
	}
	if config.Resolver.IsAccountExist == nil {
		return errors.New("Auth - config.Resolver.IsAccountExist property not set")
	}
	if config.Resolver.UpdateSession == nil {
		return errors.New("Auth - config.Resolver.UpdateSession property not set")
	}
	return nil
}
