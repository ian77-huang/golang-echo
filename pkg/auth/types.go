package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

type User[T any] struct {
	ID string

	Data T
}

type Session[T any] struct {
	ID        string
	ExpiresAt time.Time
	Data      *T
}

type Auth[TUser any, TSession any] struct {
	config *Config[TUser, TSession]
}

type Config[TUser any, TSession any] struct {
	CookieName       string
	SecretKey        string
	SessionExpiresAt int
	SessionRefreshAt int
	Resolver         Resolver[TUser, TSession]
}

type Resolver[TUser any, TSession any] struct {
	IsAccountExist func(account string) (bool, error)
	CreateUser     func(account string, password string) (*User[TUser], error)

	GetUser func(id string) (*User[TUser], error)

	GetUserByAccount func(account string) (*User[TUser], error)

	GetSession func(id string) (*Session[TSession], error)

	CreateSession func(id string, userId string, expiresAt time.Time) (*Session[TSession], error)

	UpdateSession func(id string, expiresAt time.Time, sess *TSession) (*Session[TSession], error)

	DeleteSession func(id string) (*Session[TSession], error)
}

type FieldError struct {
	Tag     string
	Message string
	Params  []interface{}
}
