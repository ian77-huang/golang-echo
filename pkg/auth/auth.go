package auth

import (
	"net/http"

	"github.com/ian77-huang/golang-echo/pkg/argon2"
	"github.com/labstack/echo/v5"
)

func GetAuth[TUser, TSession any](c *echo.Context) *Auth[TUser, TSession] {
	if auth, ok := c.Get(CONTEXT_KEY_AUTH).(*Auth[TUser, TSession]); ok {
		return auth
	}
	return nil
}

func (a *Auth[TUser, TSession]) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {

			c.Set(CONTEXT_KEY_AUTH, a)
			return next(c)
		}
	}
}

func New[TUser, TSession any](config *Config[TUser, TSession]) *Auth[TUser, TSession] {
	if config.CookieName == "" {
		config.CookieName = DEFAULT_SESSION_COOKIE_NAME
	}
	if config.SessionExpiresAt == 0 {
		config.SessionExpiresAt = DEFAULT_SESSION_EXPIRES_DAYS
	}
	if config.SessionReflashAt == 0 {
		config.SessionReflashAt = DEFAULT_SESSION_REFLASH_DAYS
	}
	if config.ErrorMessage == nil {
		// 給予預設邏輯
		config.ErrorMessage = func(tag string) string {
			return ""
		}
	}
	return &Auth[TUser, TSession]{
		config: config,
	}
}

func (a *Auth[TUser, TSession]) Register(c *echo.Context, account string, password string) (string, error) {
	confg := a.config
	resolver := confg.Resolver
	if ok, _ := resolver.IsAccountExist(account); ok {
		return "", NewError("error.auth.AccountAlreadyExists", "account already exists")
	}
	passwordHash, err := argon2.HashPassword(password)
	if err != nil {
		return "", NewError("error.auth.FailedToSecurePassword", "failed to secure password")
	}
	user, err := resolver.CreateUser(account, passwordHash)
	if err != nil {
		return "", NewError("error.auth.CreateUserError", "create user error")

	}
	if _, err := a.createSession(c, user.ID); err != nil {
		return "", err
	}

	return user.ID, nil
}

func (a *Auth[TUser, TSession]) createSession(c *echo.Context, userId string) (bool, error) {
	config := a.config
	resolver := config.Resolver

	sessionToken, err := a.GenerateSessionToken()

	sessionId := a.GenerateID(sessionToken)

	expiresAt := a.CalculateExpiry(config.SessionExpiresAt)

	session, err := resolver.CreateSession(sessionId, userId, expiresAt)
	if err != nil {
		return false, NewError("error.auth.FailedToWriteSession", "failed to write session")
	}
	tokenString, err := a.createToken(sessionId)
	if err != nil {
		return false, NewError("error.auth.FailedToIssueToken", "failed to issue token")
	}

	c.SetCookie(&http.Cookie{
		Name:     config.CookieName,
		Value:    tokenString,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	c.Set(CONTEXT_KEY_SESSION, session)

	return true, nil
}
