package auth

import (
	"errors"
	"strings"

	"github.com/ian77-huang/golang-echo/pkg/argon2"
	"github.com/labstack/echo/v5"
)

func Load[TUser, TSession any](c *echo.Context) *Auth[TUser, TSession] {
	if auth, ok := c.Get(CONTEXT_KEY_AUTH).(*Auth[TUser, TSession]); ok {
		return auth
	}
	return nil
}

func (a *Auth[TUser, TSession]) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			// config := a.config
			// resolver := config.Resolver

			if ok, _ := a.refreshSession(c); !ok {
				// return err
			}

			c.Set(CONTEXT_KEY_AUTH, a)
			return next(c)
		}
	}
}

func New[TUser, TSession any](config *Config[TUser, TSession]) *Auth[TUser, TSession] {
	if config == nil {
		panic("auth config cannot be nil")
	}
	if strings.TrimSpace(config.SecretKey) == "" {
		panic("auth secret key cannot be empty")
	}
	if config.CookieName == "" {
		config.CookieName = DEFAULT_SESSION_COOKIE_NAME
	}
	if config.SessionExpiresAt == 0 {
		config.SessionExpiresAt = DEFAULT_SESSION_EXPIRES_DAYS
	}
	if config.SessionRefreshAt == 0 {
		config.SessionRefreshAt = DEFAULT_SESSION_REFLASH_DAYS
	}

	return &Auth[TUser, TSession]{
		config: config,
	}
}

func (a *Auth[TUser, TSession]) Register(c *echo.Context, account string, password string) (string, error) {
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
		var fieldErr FieldError
		if errors.As(err, &fieldErr) {
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

func (a *Auth[TUser, TSession]) getConfig() (*Config[TUser, TSession], error) {
	if a == nil || a.config == nil {
		return nil, NewError("error.auth.InvalidConfiguration", "auth config is not configured")
	}
	return a.config, nil
}
