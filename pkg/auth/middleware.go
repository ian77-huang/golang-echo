package auth

import (
	"github.com/labstack/echo/v5"
)

func (a *Auth[TUser, TSession]) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			isValid := true
			config, err := a.getConfig()
			if err != nil {
				return err
			}
			resolver := config.Resolver

			token, err := a.getAccessToken(c)
			if err != nil {
				isValid = false
			} else {
				claims, err := a.parseToken(token)
				if err != nil {
					isValid = false
				} else {
					sessionId := a.GenerateID(claims.ID)
					sess, err := resolver.GetSession(sessionId)
					if err != nil {
						isValid = false
					} else {
						if a.isNeedRefreshSession(sess.ExpiresAt) {
							if ok, _ := a.refreshSession(c, claims.ID, sess); !ok {
								isValid = false
							}
						}
						// sess.UserID
						c.Set(CONTEXT_KEY_SESSION, sess)

						user, err := resolver.GetUser(sess.UserID)
						if err != nil || user == nil || user.ID == "" {
							return NewError("error.auth.CreateUserError", "create user error")
						}

						c.Set(CONTEXT_KEY_USER, user)
					}
				}
			}

			if !isValid {
				a.deleteAccessToken(c)
			}

			validateRule := &ValidateRule[TUser]{IsSignedIn: IsSignedIn[TUser](c), User: nil}
			if validateRule.IsSignedIn {
				validateRule.User = GetUser[TUser](c)
			}

			handled, err := config.ValidateRoute(c, validateRule)
			if !handled {
				return err
			}

			c.Set(CONTEXT_KEY_AUTH, a)

			return next(c)
		}
	}
}
