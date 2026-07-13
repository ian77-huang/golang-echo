package config

import (
	"time"

	"github.com/ian77-huang/golang-echo/internal/models/session"
	"github.com/ian77-huang/golang-echo/internal/models/users"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/pkg/cast"
	"gorm.io/gorm"
)

func Auth(db *gorm.DB) *appAuth.Auth[users.User, session.Session] {
	config := Load()
	guestOnly := &appAuth.RoutesPaths{
		Rules:       []string{"/users/login", "/users/register"},
		RedirectURL: "/",
	}
	AuthOnly := &appAuth.RoutesPaths{
		Rules:       []string{},
		RedirectURL: "/users/login",
	}
	route := &appAuth.Route[users.User]{
		GuestOnly: guestOnly,
		AuthOnly:  AuthOnly,
	}
	resolver := &appAuth.Resolver[users.User, session.Session]{
		IsAccountExist: func(account string) (bool, error) {
			return users.IsAccountExist(db, account)
		},
		CreateUser: func(account, password string) (*appAuth.User[users.User], error) {
			user, err := users.CreateUser(db, account, password)
			if err != nil {
				return nil, err
			}
			authUser := &appAuth.User[users.User]{ID: cast.IntToString(user.Id), Password: user.Password}
			return authUser, nil
		},
		GetUser: func(id string) (*appAuth.User[users.User], error) {
			userId, err := cast.StringToInt(id, 0)
			if err != nil {
				return nil, err
			}
			user, err := users.GetUser(db, userId)
			if err != nil {
				return nil, err
			}
			authUser := &appAuth.User[users.User]{ID: id, Data: user, Password: user.Password}
			return authUser, nil
		},
		GetUserByAccount: func(account string) (*appAuth.User[users.User], error) {
			user, err := users.GetUserByAccount(db, account)
			if err != nil {
				return nil, err
			}
			authUser := &appAuth.User[users.User]{ID: cast.IntToString(user.Id), Data: user, Password: user.Password}
			return authUser, nil
		},
		CreateSession: func(id, userId string, expiresAt time.Time) (*appAuth.Session[session.Session], error) {
			sess, err := session.CreateSession(db, id, userId, expiresAt)
			if err != nil {
				return nil, err
			}

			return (&appAuth.Session[session.Session]{ID: sess.ID, UserID: userId, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
		},
		UpdateSession: func(id string, expiresAt time.Time, sess *session.Session) (*appAuth.Session[session.Session], error) {
			updateSess, err := session.UpdateSession(db, id, expiresAt, sess)
			if err != nil {
				return nil, err
			}
			return (&appAuth.Session[session.Session]{ID: sess.ID, UserID: updateSess.UserID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
		},
		DeleteSession: func(id string) (*appAuth.Session[session.Session], error) {
			sess, err := session.DeleteSession(db, id)
			if err != nil {
				return nil, err
			}
			return (&appAuth.Session[session.Session]{ID: sess.ID, UserID: sess.UserID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
		},
		GetSession: func(id string) (*appAuth.Session[session.Session], error) {
			sess, err := session.GetSession(db, id)
			if err != nil {
				return nil, err
			}
			return (&appAuth.Session[session.Session]{ID: sess.ID, UserID: sess.UserID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
		},
	}
	return appAuth.New(&appAuth.Config[users.User, session.Session]{
		SecretKey: config.SecretKey,
		Resolver:  resolver,
		Route:     route,
	})
}
