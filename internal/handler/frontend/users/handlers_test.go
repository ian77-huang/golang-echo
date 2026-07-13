package users

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sessionModel "github.com/ian77-huang/golang-echo/internal/models/session"
	userModel "github.com/ian77-huang/golang-echo/internal/models/users"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/labstack/echo/v5"
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
	e.Renderer = r
	for _, tt := range []struct {
		name    string
		handler echo.HandlerFunc
	}{{"index", GetIndex}, {"login", GetLogin}, {"register", GetRegister}} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "register" {
				t.Setenv("USERS_ACCOUNT_MIN_LENGTH", "6")
				t.Setenv("USERS_PASSWORD_MIN_LENGTH", "8")
			}
			rec := httptest.NewRecorder()
			err := tt.handler(e.NewContext(httptest.NewRequest(http.MethodGet, "/users/"+tt.name, nil), rec))
			if err != nil || rec.Code != http.StatusOK || r.name == "" {
				t.Fatalf("render: name=%q status=%d err=%v", r.name, rec.Code, err)
			}
		})
	}
}

func TestGetLogoutRequiresAuth(t *testing.T) {
	e := echo.New()
	err := GetLogout(e.NewContext(httptest.NewRequest(http.MethodGet, "/users/logout", nil), httptest.NewRecorder()))
	he, ok := err.(*echo.HTTPError)
	if !ok || he.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected error: %#v", err)
	}
}

func TestGetLogoutDeletesSessionAndRedirects(t *testing.T) {
	deleted := false
	auth := appAuth.New(&appAuth.Config[userModel.User, sessionModel.Session]{SecretKey: "test-secret", Resolver: &appAuth.Resolver[userModel.User, sessionModel.Session]{
		IsAccountExist: func(string) (bool, error) { return false, nil }, CreateUser: func(string, string) (*appAuth.User[userModel.User], error) {
			return &appAuth.User[userModel.User]{ID: "user-1"}, nil
		},
		GetUser: func(id string) (*appAuth.User[userModel.User], error) {
			return &appAuth.User[userModel.User]{ID: id}, nil
		}, GetUserByAccount: func(string) (*appAuth.User[userModel.User], error) {
			return &appAuth.User[userModel.User]{ID: "user-1"}, nil
		},
		GetSession: func(id string) (*appAuth.Session[sessionModel.Session], error) {
			return &appAuth.Session[sessionModel.Session]{ID: id}, nil
		}, CreateSession: func(id, userID string, expires time.Time) (*appAuth.Session[sessionModel.Session], error) {
			return &appAuth.Session[sessionModel.Session]{ID: id, UserID: userID, ExpiresAt: expires}, nil
		},
		UpdateSession: func(id string, expires time.Time, _ *sessionModel.Session) (*appAuth.Session[sessionModel.Session], error) {
			return &appAuth.Session[sessionModel.Session]{ID: id, ExpiresAt: expires}, nil
		}, DeleteSession: func(id string) (*appAuth.Session[sessionModel.Session], error) {
			deleted = id == "user-1"
			return &appAuth.Session[sessionModel.Session]{ID: id}, nil
		},
	}})
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/users/logout", nil), rec)
	c.Set(appAuth.CONTEXT_KEY_AUTH, auth)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[userModel.User]{ID: "user-1"})
	if err := GetLogout(c); err != nil || !deleted || rec.Code != http.StatusSeeOther {
		t.Fatalf("deleted=%v status=%d err=%v", deleted, rec.Code, err)
	}
}
