package config

import (
	"github.com/ian77-huang/golang-echo/service"
)

type Menus struct {
	Name   string
	Url    string
	Childs []MenusChilds
}
type MenusChilds struct {
	Name string
	Url  string
}
type MenuRules struct {
	users []MenusChilds
	T     func(messageID string, pairs ...any) string
}
type MenuUsersRules struct {
	Path       string
	IsSignedIn bool
	IsAdmin    bool
	T          func(messageID string, pairs ...any) string
}

type ConfigDatabases struct {
	Path string
}
type ConfigUsers struct {
	MinLengthAccount  int
	MinLengthPassword int
}

type Config struct {
	SecretKey                string
	Databases                ConfigDatabases
	Users                    ConfigUsers
	AssetsPath               string
	MaxSizeUserProfileAvatar int
}

type AuthParameter struct {
	UserService    *service.UserService
	SessionService *service.SessionService
}
