package auth

import "github.com/labstack/echo/v5"

func GetUser[TUser any](c *echo.Context) *User[TUser] {
	if sess, ok := c.Get(CONTEXT_KEY_USER).(*User[TUser]); ok {
		return sess
	}
	return nil
}
func GetSession[TSession any](c *echo.Context) *Session[TSession] {
	if sess, ok := c.Get(CONTEXT_KEY_SESSION).(*Session[TSession]); ok {
		return sess
	}
	return nil
}
func (a *Auth[TUser, TSession]) getConfig() (*Config[TUser, TSession], error) {
	if a == nil || a.config == nil {
		return nil, NewError("error.auth.InvalidConfiguration", "auth config is not configured")
	}
	return a.config, nil
}
