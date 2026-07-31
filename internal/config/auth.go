package config

import (
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/pkg/store"
	"github.com/labstack/echo/v5"
)

const CACHE_KEY_SESSION_ID = "session:"

func Auth(p *AuthParameter, ss *store.StoreServer) *appAuth.Auth[model.User, model.Session] {
	config := Load()
	resolver := &appAuth.Resolver[model.User, model.Session]{
		IsAccountExist:     p.UserService.IsAccountExist,
		CreateUser:         p.UserService.CreateUser,
		GetUser:            p.UserService.GetUser,
		GetUserByAccount:   p.UserService.GetUserByAccount,
		UpdateUserPassword: p.UserService.UpdateUserPassword,
		CreateSession: func(sess *appAuth.Session[model.Session]) (*appAuth.Session[model.Session], error) {
			sess, err := p.SessionService.CreateSession(sess.ID, sess.UserID, sess.ExpiresAt)
			if err != nil {
				return nil, err
			}

			if ss != nil {
				ss.Set(CACHE_KEY_SESSION_ID+sess.ID, sess, time.Until(sess.Data.ExpiresAt))
			}

			return sess, nil
		},
		UpdateSession: func(id string, expiresAt time.Time, sessData *model.Session) (*appAuth.Session[model.Session], error) {
			var sess *appAuth.Session[model.Session]
			var err error

			if sess, err = p.SessionService.UpdateSession(id, expiresAt, sessData); err != nil {
				return nil, err
			}

			if ss != nil {
				if err = ss.Delete(CACHE_KEY_SESSION_ID + id); err != nil {
					return nil, err
				}
			}

			return sess, nil
		},
		DeleteSession: func(id string) (*appAuth.Session[model.Session], error) {
			var sess *appAuth.Session[model.Session]
			var err error

			if sess, err = p.SessionService.DeleteSession(id); err != nil {
				return nil, err
			}
			if ss != nil {
				if err := ss.Delete(CACHE_KEY_SESSION_ID + sess.ID); err != nil {
					return nil, err
				}
			}

			return sess, nil
		},
		DeleteSessionUserId: func(userId string) error {
			if err := p.SessionService.DeleteSessionUserId(userId); err != nil {
				return err
			}
			return nil
		},
		GetSession: func(id string) (*appAuth.Session[model.Session], error) {
			var sess *appAuth.Session[model.Session]
			var err error

			if ss != nil {
				if err = ss.GetByte(CACHE_KEY_SESSION_ID+id, &sess); err != nil {
					return nil, err
				}
			}

			if sess == nil {
				sess, err = p.SessionService.GetSession(id)
				if err != nil {
					return nil, err
				}

				go func() {
					if ss != nil {
						ss.Set(CACHE_KEY_SESSION_ID+id, sess, time.Until(sess.Data.ExpiresAt))
					}
				}()
			}

			return sess, nil
		},
	}
	return appAuth.New(&appAuth.Config[model.User, model.Session]{
		SecretKey: config.SecretKey,
		Resolver:  resolver,
		ValidateRoute: func(c *echo.Context, validateRule *appAuth.ValidateRule[model.User]) (bool, error) {
			rule := validateRule
			if rule.IsSignedIn {
				paths := []string{
					"/user/login",
					"/user/register",
				}
				if slices.Contains(paths, c.Path()) {
					return false, c.Redirect(http.StatusFound, "/user")
				}
			} else {
				paths := []string{
					"/user",
					"/user/profile",
					"/user/reset-password",
					"/api/user/profile",
					"/api/user/reset-password",
					"/api/user/profile/avatar",
				}
				if slices.Contains(paths, c.Path()) {
					return false, c.Redirect(http.StatusFound, "/user/login")
				}
			}

			isAdmin := rule != nil && rule.User != nil && rule.User.Data != nil && rule.User.Data.IsAdmin

			if !isAdmin {
				prefixes := []string{"/admin", "/api/admin"}
				if slices.ContainsFunc(prefixes, func(p string) bool { return strings.HasPrefix(c.Path(), p) }) {
					return false, c.Redirect(http.StatusFound, "/")
				}
			}

			return true, nil
		},
	})
}
