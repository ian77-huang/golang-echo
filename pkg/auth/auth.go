package auth

import (
	"strings"

	"github.com/labstack/echo/v5"
)

func Load[TUser any, TSession any](c *echo.Context) *Auth[TUser, TSession] {
	if auth, ok := c.Get(CONTEXT_KEY_AUTH).(*Auth[TUser, TSession]); ok {
		return auth
	}
	return nil
}

func New[TUser any, TSession any](config *Config[TUser, TSession]) *Auth[TUser, TSession] {
	if config == nil {
		panic("auth config cannot be nil")
	}
	if strings.TrimSpace(config.SecretKey) == "" {
		panic("auth secret key cannot be empty")
	}
	if err := checkResolver(config); err != nil {
		panic(err)
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

	checkRoute(config)

	return &Auth[TUser, TSession]{
		config: config,
	}
}
