package config

import (
	"time"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/pkg/cast"
	"github.com/ian77-huang/golang-echo/repository"
	"gorm.io/gorm"
)

func Auth(db *gorm.DB) *appAuth.Auth[model.User, model.Session] {
	config := Load()
	userRepository := repository.NewUserRepository(db)
	sessionReopository := repository.NewSessionRepository(db)
	guestOnly := &appAuth.RoutesPaths{
		Rules:       []string{"/user/login", "/user/register"},
		RedirectURL: "/",
	}
	AuthOnly := &appAuth.RoutesPaths{
		Rules:       []string{},
		RedirectURL: "/user/login",
	}
	route := &appAuth.Route[model.User]{
		GuestOnly: guestOnly,
		AuthOnly:  AuthOnly,
	}
	resolver := &appAuth.Resolver[model.User, model.Session]{
		IsAccountExist: userRepository.IsAccountExist,
		CreateUser: func(account, password string) (*appAuth.User[model.User], error) {
			user, err := userRepository.CreateUser(account, password)
			if err != nil {
				return nil, err
			}
			authUser := &appAuth.User[model.User]{ID: cast.IntToString(user.Id), Password: user.Password, Data: user}
			return authUser, nil
		},
		GetUser: func(id string) (*appAuth.User[model.User], error) {
			userId, err := cast.StringToInt(id, 0)
			if err != nil {
				return nil, err
			}
			user, err := userRepository.GetUser(userId)
			if err != nil {
				return nil, err
			}
			authUser := &appAuth.User[model.User]{ID: id, Data: user, Password: user.Password}
			return authUser, nil
		},
		GetUserByAccount: func(account string) (*appAuth.User[model.User], error) {
			user, err := userRepository.GetUserByAccount(account)
			if user == nil {
				return nil, nil
			}
			if err != nil {
				return nil, err
			}
			authUser := &appAuth.User[model.User]{ID: cast.IntToString(user.Id), Data: user, Password: user.Password}
			return authUser, nil
		},
		CreateSession: func(id, userId string, expiresAt time.Time) (*appAuth.Session[model.Session], error) {
			sess, err := sessionReopository.CreateSession(id, userId, expiresAt)
			if err != nil {
				return nil, err
			}

			return (&appAuth.Session[model.Session]{ID: sess.ID, UserID: userId, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
		},
		UpdateSession: func(id string, expiresAt time.Time, sess *model.Session) (*appAuth.Session[model.Session], error) {
			updateSess, err := sessionReopository.UpdateSession(id, expiresAt, sess)
			if err != nil {
				return nil, err
			}
			return (&appAuth.Session[model.Session]{ID: sess.ID, UserID: updateSess.UserID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
		},
		DeleteSession: func(id string) (*appAuth.Session[model.Session], error) {
			sess, err := sessionReopository.DeleteSession(id)
			if err != nil {
				return nil, err
			}
			return (&appAuth.Session[model.Session]{ID: sess.ID, UserID: sess.UserID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
		},
		GetSession: func(id string) (*appAuth.Session[model.Session], error) {
			sess, err := sessionReopository.GetSession(id)
			if err != nil {
				return nil, err
			}
			return (&appAuth.Session[model.Session]{ID: sess.ID, UserID: sess.UserID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
		},
	}
	return appAuth.New(&appAuth.Config[model.User, model.Session]{
		SecretKey: config.SecretKey,
		Resolver:  resolver,
		Route:     route,
	})
}
