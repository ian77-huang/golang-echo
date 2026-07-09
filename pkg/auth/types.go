package auth

import (
	"time"

	"github.com/labstack/echo/v5"
)

type User[T any] struct {
	ID string

	Data T
}
type Session[T any] struct {
	Data T
}
type Auth[TUser any, TSession any] struct {
	config *Config[TUser, TSession]
}
type Config[TUser any, TSession any] struct {
	CookieName       string
	SecretKey        string
	SessionExpiresAt int
	SessionReflashAt int
	Resolver         Resolver[TUser, TSession]
}
type CustomContext[TUser, TSession any] struct {
	*echo.Context
	Auth *Auth[TUser, TSession] // 💡 直接用強型別欄位存取
}

type Resolver[TUser any, TSession any] struct {
	IsAccountExist func(account string) (bool, error)
	CreateUser     func(account string, password string) (*User[TUser], error)

	GetUser func(id string) (*User[TUser], error)

	GetUserByAccount func(account string) (*User[TUser], error)

	GetSession func(id string) (*TSession, error)

	CreateSession func(id string, userId string, expiresAt time.Time) (*TSession, error)

	UpdateSession func(id string, sess *TSession) (*TSession, error)

	DeleteSession func(id string) (*TSession, error)
}

// export type Resolver<
// 	TUser extends Record<string, any> = Record<string, any>,
// 	TSession extends Record<string, any> = Record<string, any>
// > = {
// 	getUser: (id: string) => Awaitable<User<TUser> | null>;
// 	getUserByAccount: (account: string) => Awaitable<User<TUser> | null>;
// 	getSession: (sessionId: string) => Awaitable<Session<TSession> | null>;
// 	createSession: (
// 		sessionId: string,
// 		userId: string,
// 		expiresAt: Date
// 	) => Awaitable<Session<TSession>>;
// 	updateSession: (session: Session<TSession>) => Awaitable<Session<TSession>>;
// 	createUser: (user: Omit<User<TUser>, 'id'>) => Awaitable<User<TUser>>;
// 	deleteSession: (sessionId: string) => Awaitable<boolean>;
// };

type FieldError struct {
	Tag     string
	Message string
	Params  []interface{}
}

// getUser: (id: string) => Awaitable<User<TUser> | null>;
// getUserByAccount: (account: string) => Awaitable<User<TUser> | null>;
// getSession: (sessionId: string) => Awaitable<Session<TSession> | null>;
// createSession: (sessionId: string, userId: string, expiresAt: Date) => Awaitable<Session<TSession>>;
// updateSession: (session: Session<TSession>) => Awaitable<Session<TSession>>;
// createUser: (user: Omit<User<TUser>, 'id'>) => Awaitable<User<TUser>>;
// deleteSession: (sessionId: string) => Awaitable<boolean>;
