package config

import (
	"os"
	"strconv"
	"time"

	"github.com/ian77-huang/golang-echo/internal/locales"
	"github.com/ian77-huang/golang-echo/internal/models/session"
	"github.com/ian77-huang/golang-echo/internal/models/users"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/pkg/cast"
	appi18n "github.com/ian77-huang/golang-echo/pkg/i18n"
	"github.com/ian77-huang/golang-echo/pkg/renderer"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserMenus struct {
	Name string
	Url  string
}

type ConfigDatabases struct {
	Path string
}
type ConfigUsers struct {
	MinLengthAccount  int
	MinLengthPassword int
}

type Config struct {
	Databases ConfigDatabases
	Users     ConfigUsers
}

func Load() Config {
	return Config{
		Databases: ConfigDatabases{
			Path: os.Getenv("DATABASE_PATH"),
		},
		Users: ConfigUsers{
			MinLengthAccount:  cast.Int(os.Getenv("USERS_ACCOUNT_MIN_LENGTH"), 6),
			MinLengthPassword: cast.Int(os.Getenv("USERS_PASSWORD_MIN_LENGTH"), 8),
		},
	}
}

func I18n() (*appi18n.I18n, error) {
	return appi18n.New(appi18n.Config{
		DefaultLang:            "zh-TW",
		SupportedLanguageCodes: []string{"zh-TW", "en"},
		MessageFS:              locales.FS,
		MessageFiles: []string{
			"active.zh-TW.toml",
			"active.en.toml",
			"errors.en.toml",
			"errors.zh-TW.toml",
			"placeholders.en.toml",
			"placeholders.zh-TW.toml",
			"users.en.toml",
			"users.zh-TW.toml",
			"validations.en.toml",
			"validations.zh-TW.toml",
		},
	})
}

func RendererTemplate(options ...renderer.Option) *renderer.TemplateConfig {
	runtime := renderer.RuntimeConfig{}
	for _, option := range options {
		option(&runtime)
	}

	config := &renderer.TemplateConfig{
		BasePath: "internal/views",
		Layouts: map[string]renderer.TemplateNode{
			"frontend": {
				FilePath: "layout.html",
				Layouts: map[string]renderer.TemplateNode{
					"users": {
						FilePath: "layout.html",
					},
				},
			},
			"admin": {
				FilePath: "layout.html",
			},
		},
		SharedTmplPaths: []string{"base.html"},
		Runtime:         runtime,
		SharedData: func(c *echo.Context, layoutNames []string) map[string]any {
			realPath := c.Request().URL.Path

			users := []UserMenus{}
			switch realPath {
			case "/users/login":
				users = append(users, UserMenus{Name: "register", Url: "/users/register"})
			case "/users/register":
				users = append(users, UserMenus{Name: "login", Url: "/users/login"})
			}

			lang, _ := c.Get("lang").(string)
			return map[string]any{
				"Lang":  lang,
				"Users": users,
			}
		},
	}

	return config
}

func DB() *gorm.DB {
	config := Load()

	db, err := gorm.Open(sqlite.Open(config.Databases.Path), &gorm.Config{})
	if err != nil {
		panic("=== Error：Unable to connect to the database. ====")
	}
	return db
}

func Auth(db *gorm.DB) *appAuth.Auth[users.User, session.Session] {
	return appAuth.New(&appAuth.Config[users.User, session.Session]{
		Resolver: appAuth.Resolver[users.User, session.Session]{
			IsAccountExist: func(account string) (bool, error) {
				return users.IsAccountExist(db, account)
			},
			CreateUser: func(account, password string) (*appAuth.User[users.User], error) {
				user, err := users.CreateUser(db, account, password)
				if err != nil {
					return nil, err
				}
				authUser := &appAuth.User[users.User]{ID: strconv.Itoa(user.Id)}
				return authUser, nil
			},
			CreateSession: func(id, userId string, expiresAt time.Time) (*appAuth.Session[session.Session], error) {
				sess, err := session.CreateSession(db, id, userId, expiresAt)
				if err != nil {
					return nil, err
				}

				return (&appAuth.Session[session.Session]{ID: sess.ID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
			},
			// CreateSession: func(id, userId string, expiresAt time.Time) (*session.Session, error) {
			// 	return session.CreateSession(db, id, userId, expiresAt)
			// },
			UpdateSession: func(id string, expiresAt time.Time, sess *session.Session) (*appAuth.Session[session.Session], error) {
				sess, err := session.UpdateSession(db, id, expiresAt, sess)
				if err != nil {
					return nil, err
				}
				return (&appAuth.Session[session.Session]{ID: sess.ID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
			},
			// UpdateSession: func(id string, expiresAt time.Time, sess *session.Session) (*session.Session, error) {
			// 	return session.UpdateSession(db, id, expiresAt, sess)
			// },
			DeleteSession: func(id string) (*appAuth.Session[session.Session], error) {
				sess, err := session.DeleteSession(db, id)
				if err != nil {
					return nil, err
				}
				return (&appAuth.Session[session.Session]{ID: sess.ID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
			},
			// DeleteSession: func(id string) (*session.Session, error) {
			// 	return session.DeleteSession(db, id)
			// },
			GetSession: func(id string) (*appAuth.Session[session.Session], error) {
				sess, err := session.GetSession(db, id)
				if err != nil {
					return nil, err
				}
				return (&appAuth.Session[session.Session]{ID: sess.ID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
			},
			// GetSession: func(id string) (*session.Session, error) {
			// 	return session.GetSession(db, id)
			// },
		},
	})
}
