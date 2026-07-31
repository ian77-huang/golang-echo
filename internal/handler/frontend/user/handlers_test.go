package user

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type captureRenderer struct {
	name string
	data map[string]any
}

func (r *captureRenderer) Render(_ *echo.Context, _ io.Writer, name string, data any) error {
	r.name = name
	r.data, _ = data.(map[string]any)
	return nil
}

func TestRenderHandlers(t *testing.T) {
	e := echo.New()
	r := &captureRenderer{}
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.UserProfile{}); err != nil {
		t.Fatal(err)
	}
	user := &model.User{Account: "tester"}
	if err := db.Create(user).Error; err != nil {
		t.Fatal(err)
	}
	h := &UserHandler{DB: db}
	e.Renderer = r
	for _, tt := range []struct {
		name    string
		handler echo.HandlerFunc
	}{{"index", h.GetIndex}, {"login", h.GetLogin}, {"register", h.GetRegister}} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "register" {
				t.Setenv("USER_ACCOUNT_MIN_LENGTH", "6")
				t.Setenv("USER_PASSWORD_MIN_LENGTH", "8")
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/user/"+tt.name, nil), rec)
			if tt.name == "index" {
				c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: strconv.Itoa(user.Id)})
			}
			err := tt.handler(c)
			if err != nil || rec.Code != http.StatusOK || r.name == "" {
				t.Fatalf("render: name=%q status=%d err=%v", r.name, rec.Code, err)
			}
		})
	}
}

func TestGetProfileAndResetPasswordRender(t *testing.T) {
	e := echo.New()
	r := &captureRenderer{}
	e.Renderer = r
	h := &UserHandler{}
	for _, tt := range []struct {
		name    string
		handler echo.HandlerFunc
		want    string
	}{{"profile", h.GetProfile, "frontend:user:/profile.html"}, {"reset-password", h.GetResetPassword, "frontend:user:/reset-password.html"}} {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/user/"+tt.name, nil), rec)
			if err := tt.handler(c); err != nil || rec.Code != http.StatusOK || r.name != tt.want {
				t.Fatalf("render: name=%q want=%q status=%d err=%v", r.name, tt.want, rec.Code, err)
			}
		})
	}
}

func TestGetLogoutRequiresAuth(t *testing.T) {
	e := echo.New()
	h := &UserHandler{}
	err := h.GetLogout(e.NewContext(httptest.NewRequest(http.MethodGet, "/user/logout", nil), httptest.NewRecorder()))
	he, ok := err.(*echo.HTTPError)
	if !ok || he.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestGetLogoutDeletesSessionAndRedirects(t *testing.T) {
	deleted := false
	h := &UserHandler{}
	auth := appAuth.New(&appAuth.Config[model.User, model.Session]{SecretKey: "test-secret", Resolver: &appAuth.Resolver[model.User, model.Session]{
		IsAccountExist: func(string) (bool, error) { return false, nil }, CreateUser: func(string, string) (*appAuth.User[model.User], error) {
			return &appAuth.User[model.User]{ID: "user-1"}, nil
		},
		GetUser: func(id string) (*appAuth.User[model.User], error) {
			return &appAuth.User[model.User]{ID: id}, nil
		}, GetUserByAccount: func(string) (*appAuth.User[model.User], error) {
			return &appAuth.User[model.User]{ID: "user-1"}, nil
		},
		GetSession: func(id string) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id}, nil
		},
		CreateSession: func(sess *appAuth.Session[model.Session]) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: sess.ID, UserID: sess.UserID, ExpiresAt: sess.ExpiresAt}, nil
		},
		UpdateSession: func(id string, expires time.Time, _ *model.Session) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id, ExpiresAt: expires}, nil
		}, DeleteSession: func(id string) (*appAuth.Session[model.Session], error) {
			deleted = id == "user-1"
			return &appAuth.Session[model.Session]{ID: id}, nil
		},
	}})
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/user/logout", nil), rec)
	c.Set(appAuth.CONTEXT_KEY_AUTH, auth)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "user-1"})
	c.Set(appAuth.CONTEXT_KEY_SESSION, &appAuth.Session[model.Session]{ID: "user-1"})
	if err := h.GetLogout(c); err != nil || !deleted || rec.Code != http.StatusFound {
		t.Fatalf("deleted=%v status=%d err=%v", deleted, rec.Code, err)
	}
}

func TestGetLogoutActionLogoutError(t *testing.T) {
	h := &UserHandler{}
	auth := appAuth.New(&appAuth.Config[model.User, model.Session]{SecretKey: "test-secret", Resolver: &appAuth.Resolver[model.User, model.Session]{
		IsAccountExist: func(string) (bool, error) { return false, nil }, CreateUser: func(string, string) (*appAuth.User[model.User], error) {
			return &appAuth.User[model.User]{ID: "user-1"}, nil
		},
		GetUser: func(id string) (*appAuth.User[model.User], error) {
			return &appAuth.User[model.User]{ID: id}, nil
		}, GetUserByAccount: func(string) (*appAuth.User[model.User], error) {
			return &appAuth.User[model.User]{ID: "user-1"}, nil
		},
		GetSession: func(id string) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id}, nil
		},
		CreateSession: func(sess *appAuth.Session[model.Session]) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: sess.ID, UserID: sess.UserID, ExpiresAt: sess.ExpiresAt}, nil
		},
		UpdateSession: func(id string, expires time.Time, _ *model.Session) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id, ExpiresAt: expires}, nil
		}, DeleteSession: func(id string) (*appAuth.Session[model.Session], error) {
			return nil, errors.New("delete session failed")
		},
	}})
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/user/logout", nil), rec)
	c.Set(appAuth.CONTEXT_KEY_AUTH, auth)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "user-1"})
	c.Set(appAuth.CONTEXT_KEY_SESSION, &appAuth.Session[model.Session]{ID: "user-1"})

	err := h.GetLogout(c)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
