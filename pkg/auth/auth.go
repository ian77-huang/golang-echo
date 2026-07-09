package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ian77-huang/golang-echo/pkg/argon2"
	"github.com/labstack/echo/v5"
)

func GetAuth[TUser, TSession any](c *echo.Context) *Auth[TUser, TSession] {
	if auth, ok := c.Get(CONTEXT_KEY_AUTH).(*Auth[TUser, TSession]); ok {
		return auth
	}
	return nil
}
func New[TUser, TSession any](config *Config[TUser, TSession]) *Auth[TUser, TSession] {
	if config.CookieName == "" {
		config.CookieName = DEFAULT_SESSION_COOKIE_NAME
	}
	if config.SessionExpiresAt == 0 {
		config.SessionExpiresAt = DEFAULT_SESSION_EXPIRES_DAYS
	}
	if config.SessionReflashAt == 0 {
		config.SessionReflashAt = DEFAULT_SESSION_REFLASH_DAYS
	}
	return &Auth[TUser, TSession]{
		config: config,
	}
}

func (a *Auth[TUser, TSession]) Register(c *echo.Context, account string, password string) (string, error) {
	log.Print("\n========= 1 %+v =========\n", a.config)
	confg := a.config
	log.Print("\n========= 1.1 =========\n")
	resolver := confg.Resolver
	log.Print("\n========= 2 =========\n")
	if ok, _ := resolver.IsAccountExist(account); ok {
		return "", NewError("error.auth.AccountAlreadyExists", "account already exists")
	}
	log.Print("\n========= 3 =========\n")
	passwordHash, err := argon2.HashPassword(password)
	if err != nil {
		return "", NewError("error.auth.FailedToSecurePassword", "failed to secure password")
	}
	log.Print("\n========= 4 =========\n")
	user, err := resolver.CreateUser(account, passwordHash)
	if err != nil {
		return "", NewError("error.auth.CreateUserError", "create user error")

	}
	log.Print("\n========= 5 =========\n")
	if _, err := a.createSession(c, user.ID); err != nil {
		return "", err
	}

	return user.ID, nil

	// const { resolver, locale } = options;

	// 			const checkUserAccount = await resolver.getUserByAccount(user.account);
	// 			if (checkUserAccount !== null) {
	// 				throw new Error(locale('register.user-already-exists-1'));
	// 			}

	// 			const { hash } = await hashPassword(user.password);

	// 			const { id: userId } = await resolver.createUser({ ...user, password: hash });

	// 			const { createSession } = session<TUser, TSession>(options);

	// 			await createSession(event, userId);

	// return Promise.resolve(true);
}

type JwtCustomClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func (a *Auth[TUser, TSession]) getJwtSecret() []byte {
	return []byte(a.config.SecretKey)
}
func (a *Auth[TUser, TSession]) createToken(id string) (string, error) {
	claims := &JwtCustomClaims{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(a.getJwtSecret())
}
func (a *Auth[TUser, TSession]) createSession(c *echo.Context, userId string) (bool, error) {
	config := a.config
	resolver := config.Resolver

	sessionToken, err := a.GenerateSessionToken()

	sessionId := a.GenerateID(sessionToken)

	expiresAt := a.CalculateExpiry(config.SessionExpiresAt)

	session, err := resolver.CreateSession(sessionId, userId, expiresAt)
	if err != nil {
		return false, NewError("error.auth.FailedToWriteSession", "failed to write session")
	}
	tokenString, err := a.createToken(sessionId)
	if err != nil {
		return false, NewError("error.auth.FailedToIssueToken", "failed to issue token")
	}

	c.SetCookie(&http.Cookie{
		Name:     config.CookieName,
		Value:    tokenString,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	c.Set(CONTEXT_KEY_SESSION, session)

	return true, nil
}

func (a *Auth[TUser, TSession]) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {

			log.Printf("======= %+v ========", a.config.CookieName)
			c.Set(CONTEXT_KEY_AUTH, a)
			return next(c)
		}
	}
}
