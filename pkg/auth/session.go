package auth

import (
	"time"

	"github.com/labstack/echo/v5"
)

func (a *Auth[TUser, TSession]) isNeedRefreshSession(expires time.Time) bool {
	config := a.config
	refreshTime := expires.Add(-time.Hour * 24 * time.Duration(config.SessionReflashAt))

	return time.Now().After(refreshTime) || time.Now().Equal(refreshTime)
}

func (a *Auth[TUser, TSession]) createSession(c *echo.Context, userId string) (bool, error) {
	config := a.config
	resolver := config.Resolver

	expiresAt := a.extendDateDays(config.SessionExpiresAt)

	sessionToken, err := a.GenerateSessionToken()
	if err != nil {
		a.deleteAccessToken(c)
		return false, NewError("error.auth.FailedToIssueToken", "failed to generate session token")
	}

	sessionId := a.GenerateID(sessionToken)

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
	config := a.config
	resolver := config.Resolver

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
	if err != nil {
		a.deleteAccessToken(c)
		return false, err
	}

	if a.isNeedRefreshSession(sess.ExpiresAt) {
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
