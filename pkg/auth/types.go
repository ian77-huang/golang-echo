package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type JwtCustomClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

type User[T any] struct {
	ID       string
	Password string
	Data     *T
}

type Session[T any] struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
	Data      *T
}
type ValidateRule[TUser any] struct {
	Path       string
	IsSignedIn bool
	User       *User[TUser]
}
type RoutesPaths struct {
	Rules       []string
	RedirectURL string
}
type Route[TUser any] struct {
	GuestOnly       *RoutesPaths
	AuthOnly        *RoutesPaths
	SpecialValidate func(c *echo.Context, validateRule *ValidateRule[TUser])
}

type Auth[TUser any, TSession any] struct {
	config *Config[TUser, TSession]
}

type Config[TUser any, TSession any] struct {
	CookieName       string
	SecretKey        string
	SessionExpiresAt int
	SessionRefreshAt int
	Resolver         *Resolver[TUser, TSession]
	Route            *Route[TUser]
}

type Resolver[TUser any, TSession any] struct {
	IsAccountExist     func(account string) (bool, error)
	CreateUser         func(account string, password string) (*User[TUser], error)
	GetUser            func(id string) (*User[TUser], error)
	GetUserByAccount   func(account string) (*User[TUser], error)
	UpdateUserPassword func(id string, passwordHash string) (*User[TUser], error)

	GetSession    func(id string) (*Session[TSession], error)
	CreateSession func(id string, userId string, expiresAt time.Time) (*Session[TSession], error)
	UpdateSession func(id string, expiresAt time.Time, sess *TSession) (*Session[TSession], error)
	DeleteSession func(id string) (*Session[TSession], error)
}

type FieldError struct {
	Tag     string
	Message string
	Params  []interface{}
}
