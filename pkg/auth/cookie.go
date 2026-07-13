package auth

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

var readCookie = func(c *echo.Context, name string) (*http.Cookie, error) {
	return c.Cookie(name)
}

func (a *Auth[TUser, TSession]) getAccessToken(c *echo.Context) (string, error) {
	config := a.config

	token, err := readCookie(c, config.CookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return "", NewError("error.auth.CookieNotFound", "Specified cookie not found")
		}
		return "", NewError("error.auth.ReadCookieFailed", "Failed to read cookie")
	}
	return token.Value, nil
}

func (a *Auth[TUser, TSession]) setAccessToken(c *echo.Context, tokenString string, expiresAt time.Time) {
	config := a.config

	c.SetCookie(&http.Cookie{
		Name:     config.CookieName,
		Value:    tokenString,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	})
}

func (a *Auth[TUser, TSession]) deleteAccessToken(c *echo.Context) {
	config := a.config

	c.SetCookie(&http.Cookie{
		Name:     config.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	})
}
