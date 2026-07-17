package config

import (
	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
)

func Auth(p *AuthParameter) *appAuth.Auth[model.User, model.Session] {
	config := Load()

	guestOnly := &appAuth.RoutesPaths{
		Rules:       []string{"/user/login", "/user/register"},
		RedirectURL: "/",
	}
	AuthOnly := &appAuth.RoutesPaths{
		Rules:       []string{"/user/profile", "/user/reset-password"},
		RedirectURL: "/user/login",
	}
	route := &appAuth.Route[model.User]{
		GuestOnly: guestOnly,
		AuthOnly:  AuthOnly,
	}
	resolver := &appAuth.Resolver[model.User, model.Session]{
		IsAccountExist:     p.UserService.IsAccountExist,
		CreateUser:         p.UserService.CreateUser,
		GetUser:            p.UserService.GetUser,
		GetUserByAccount:   p.UserService.GetUserByAccount,
		UpdateUserPassword: p.UserService.UpdateUserPassword,
		CreateSession:      p.SessionService.CreateSession,
		UpdateSession:      p.SessionService.UpdateSession,
		DeleteSession:      p.SessionService.DeleteSession,
		GetSession:         p.SessionService.GetSession,
	}
	return appAuth.New(&appAuth.Config[model.User, model.Session]{
		SecretKey: config.SecretKey,
		Resolver:  resolver,
		Route:     route,
	})
}
