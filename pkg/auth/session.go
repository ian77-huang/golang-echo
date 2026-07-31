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

	if _, err := resolver.DeleteSession(sessionId); err != nil {
		return false, err
	}

	if resolver.CreateSession == nil {
		a.deleteAccessToken(c)
		return false, NewError("error.auth.FailedToWriteSession", "session storage is not configured")
	}

	// _, err = resolver.CreateSession(sessionId, userId, expiresAt)
	_, err = resolver.CreateSession(&Session[TSession]{ID: sessionId, UserID: userId, ExpiresAt: expiresAt})
	if err != nil {
		a.deleteAccessToken(c)
		return false, NewError("error.auth.FailedToWriteSession", "failed to write session")
	}

	tokenString, err := a.createToken(sessionToken, expiresAt)
	if err != nil {
		a.deleteAccessToken(c)
		return false, NewError("error.auth.FailedToIssueToken", "failed to issue token")
	}

	a.setAccessToken(c, tokenString, expiresAt)

	return true, nil
}

func (a *Auth[TUser, TSession]) refreshSession(c *echo.Context, tokenId string, sess *Session[TSession]) (bool, error) {
	config, err := a.getConfig()
	if err != nil {
		return false, err
	}
	resolver := config.Resolver

	if a.isNeedRefreshSession(sess.ExpiresAt) {
		if resolver.UpdateSession == nil {
			a.deleteAccessToken(c)
			return false, NewError("error.auth.InvalidConfiguration", "session update resolver is not configured")
		}
		expiresAt := a.extendDateDays(config.SessionExpiresAt)

		_, err := resolver.UpdateSession(sess.ID, expiresAt, sess.Data)
		if err != nil {
			return false, err
		}

		tokenString, err := a.createToken(tokenId, expiresAt)
		if err != nil {
			return false, NewError("error.auth.FailedToIssueToken", "failed to issue token")
		}

		a.setAccessToken(c, tokenString, expiresAt)
	}

	return true, nil
}
