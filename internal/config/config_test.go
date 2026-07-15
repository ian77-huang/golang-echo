package config

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/repository"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLoadUsesEnvironmentAndDefaults(t *testing.T) {
	t.Setenv("USER_ACCOUNT_MIN_LENGTH", "9")
	t.Setenv("USER_PASSWORD_MIN_LENGTH", "12")
	t.Setenv("SECRET_KEY", "secret")
	t.Setenv("DATABASE_PATH", "test.db")
	c := Load()
	if c.Users.MinLengthAccount != 9 || c.Users.MinLengthPassword != 12 || c.SecretKey != "secret" || c.Databases.Path != "test.db" {
		t.Fatalf("unexpected config: %#v", c)
	}
}

func TestLoadPanicsForInvalidLengthEnvironment(t *testing.T) {
	for _, key := range []string{"USER_ACCOUNT_MIN_LENGTH", "USERS_PASSWORD_MIN_LENGTH"} {
		t.Run(key, func(t *testing.T) {
			t.Setenv(key, "invalid")
			defer func() {
				if recover() == nil {
					t.Fatal("expected panic")
				}
			}()
			Load()
		})
	}
}

func TestI18nAndRendererTemplate(t *testing.T) {
	translator, err := I18n()
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	if err := translator.Middleware()(func(c *echo.Context) error { return nil })(c); err != nil {
		t.Fatal(err)
	}
	template := RendererTemplate()
	data := template.SharedData(c, []string{"frontend"})
	if template.BasePath != "internal/views" || data["Lang"] == "" || data["Menus"] == nil {
		t.Fatalf("template=%#v data=%#v", template, data)
	}
}

func TestAuthBuildsWorkingResolvers(t *testing.T) {
	t.Setenv("SECRET_KEY", "test-secret")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Session{}); err != nil {
		t.Fatal(err)
	}
	auth := Auth(&AuthParameter{
		UserService:    service.NewUserService(db),
		SessionService: service.NewSessionService(db),
	})
	e := echo.New()
	rec := httptest.NewRecorder()
	id, err := auth.ActionRegister(e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), rec), "config-user", "password")
	if err != nil || id == "" || len(rec.Result().Cookies()) != 1 {
		t.Fatalf("id=%q cookies=%#v err=%v", id, rec.Result().Cookies(), err)
	}
	loginRec := httptest.NewRecorder()
	loggedIn, err := auth.ActionLogin(e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), loginRec), "config-user", "password")
	if err != nil || !loggedIn || len(loginRec.Result().Cookies()) != 1 {
		t.Fatalf("login=%v cookies=%#v err=%v", loggedIn, loginRec.Result().Cookies(), err)
	}
	var saved model.Session
	if err := db.Where("userId = ?", id).First(&saved).Error; err != nil {
		t.Fatal(err)
	}
	if err := auth.Middleware()(func(c *echo.Context) error { return nil })(e.NewContext(requestWithCookie(loginRec.Result().Cookies()[0]), httptest.NewRecorder())); err != nil {
		t.Fatal(err)
	}
	if _, err := auth.ActionLogout(contextWithUser(e, auth, id)); err != nil {
		t.Fatal(err)
	}
	sessionRepository := repository.NewSessionRepository(db)
	if _, err := sessionRepository.UpdateSession(saved.ID, time.Now().Add(time.Hour), &saved); err != nil {
		t.Fatal(err)
	}
}

func TestAuthResolversPropagateDatabaseErrors(t *testing.T) {
	t.Setenv("SECRET_KEY", "test-secret")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Session{}); err != nil {
		t.Fatal(err)
	}
	auth := Auth(&AuthParameter{
		UserService:    service.NewUserService(db),
		SessionService: service.NewSessionService(db),
	})
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := auth.IsAccountExist("user"); err == nil {
		t.Fatal("expected account lookup error")
	}
	if _, err := auth.CreateUser("user", "password"); err == nil {
		t.Fatal("expected create user error")
	}
	if _, err := auth.ActionLogin(echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder()), "user", "password"); err == nil {
		t.Fatal("expected login lookup error")
	}
	c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "1"})
	if _, err := auth.ActionLogout(c); err == nil {
		t.Fatal("expected delete session error")
	}
}

func requestWithCookie(cookie *http.Cookie) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookie)
	return req
}
func contextWithUser(e *echo.Echo, auth *appAuth.Auth[model.User, model.Session], id string) *echo.Context {
	c := e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: id})
	c.Set(appAuth.CONTEXT_KEY_AUTH, auth)
	return c
}

func TestDBOpensConfiguredSqliteFile(t *testing.T) {
	t.Setenv("DATABASE_PATH", filepath.Join(t.TempDir(), "app.db"))
	if db := DB(); db == nil {
		t.Fatal("expected database")
	}
}

func TestDBPanicsForUnusablePath(t *testing.T) {
	t.Setenv("DATABASE_PATH", "/dev/null/app.db")
	defer func() {
		if recover() == nil {
			t.Fatal("expected DB panic")
		}
	}()
	DB()
}
