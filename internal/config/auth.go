package config

import (
	"net/http"
	"slices"
	"strings"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/labstack/echo/v5"
)

func Auth(p *AuthParameter) *appAuth.Auth[model.User, model.Session] {
	config := Load()
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
