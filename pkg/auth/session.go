package auth

import (
	"time"

	"github.com/labstack/echo/v5"
)

func (a *Auth[TUser, TSession]) isNeedRefreshSession(expires time.Time) bool {
	config := a.config
	refreshTime := expires.Add(-time.Hour * 24 * time.Duration(config.SessionRefreshAt))

	return time.Now().After(refreshTime) || time.Now().Equal(refreshTime)
}

func (a *Auth[TUser, TSession]) createSession(c *echo.Context, userId string) (bool, error) {
	config, err := a.getConfig()
	if err != nil {
		return false, err
	}
	resolver := config.Resolver

	expiresAt := a.extendDateDays(config.SessionExpiresAt)

	sessionToken, err := a.GenerateSessionToken()
	if err != nil {
		a.deleteAccessToken(c)
		return false, NewError("error.auth.FailedToIssueToken", "failed to generate session token")
	}

	sessionId := a.GenerateID(sessionToken)

	if resolver.CreateSession == nil {
		a.deleteAccessToken(c)
		return false, NewError("error.auth.FailedToWriteSession", "session storage is not configured")
	}

	_, err = resolver.CreateSession(sessionId, userId, expiresAt)
	if err != nil {
		a.deleteAccessToken(c)
		return false, NewError("error.auth.FailedToWriteSession", "failed to write session")
	}

	tokenString, err := a.createToken(sessionId, expiresAt)
	if err != nil {
		a.deleteAccessToken(c)
		return false, NewError("error.auth.FailedToIssueToken", "failed to issue token")
	}

	a.setAccessToken(c, tokenString, expiresAt)

	return true, nil
}

func (a *Auth[TUser, TSession]) refreshSession(c *echo.Context) (bool, error) {
	config, err := a.getConfig()
	if err != nil {
		return false, err
	}
	resolver := config.Resolver
	if resolver.GetSession == nil {
		return false, NewError("error.auth.InvalidConfiguration", "session lookup resolver is not configured")
	}

	expiresAt := a.extendDateDays(config.SessionExpiresAt)

	token, err := a.getAccessToken(c)
	if err != nil {
		a.deleteAccessToken(c)
		return false, err
	}

	claims, err := a.parseToken(token)
	if err != nil {
		a.deleteAccessToken(c)
		return false, err
	}

	sessionId := claims.ID

	sess, err := resolver.GetSession(sessionId)
	if err != nil || sess == nil {
		a.deleteAccessToken(c)
		if err != nil {
			return false, err
		}
		return false, NewError("error.auth.InvalidConfiguration", "session lookup returned no session")
	}

	if a.isNeedRefreshSession(sess.ExpiresAt) {
		if resolver.UpdateSession == nil {
			a.deleteAccessToken(c)
			return false, NewError("error.auth.InvalidConfiguration", "session update resolver is not configured")
		}
		_, err := resolver.UpdateSession(sessionId, expiresAt, sess.Data)
		if err != nil {
			a.deleteAccessToken(c)
			return false, err
		}

		tokenString, err := a.createToken(sessionId, expiresAt)
		if err != nil {
			a.deleteAccessToken(c)
			return false, NewError("error.auth.FailedToIssueToken", "failed to issue token")
		}

		a.setAccessToken(c, tokenString, expiresAt)
	}

	return true, nil
}
