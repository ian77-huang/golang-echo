package users

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	appValidator "github.com/ian77-huang/golang-echo/pkg/validator"
	"github.com/labstack/echo/v5"
)

func TestPostLoginUsesAuthFromMiddleware(t *testing.T) {
	hash, err := appAuth.HashPassword("password123")
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiUserHandler{}
	auth := appAuth.New(&appAuth.Config[model.User, model.Session]{SecretKey: "test-secret", Resolver: &appAuth.Resolver[model.User, model.Session]{
		IsAccountExist: func(string) (bool, error) { return false, nil }, CreateUser: func(string, string) (*appAuth.User[model.User], error) {
			return &appAuth.User[model.User]{ID: "user-1"}, nil
		},
		GetUser: func(id string) (*appAuth.User[model.User], error) {
			return &appAuth.User[model.User]{ID: id}, nil
		}, GetUserByAccount: func(string) (*appAuth.User[model.User], error) {
			return &appAuth.User[model.User]{ID: "user-1", Password: hash}, nil
		},
		GetSession: func(id string) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id}, nil
		}, CreateSession: func(id, userID string, expires time.Time) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id, UserID: userID, ExpiresAt: expires}, nil
		},
		UpdateSession: func(id string, expires time.Time, _ *model.Session) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id, ExpiresAt: expires}, nil
		}, DeleteSession: func(id string) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id}, nil
		},
	}})
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"account":"tester","password":"password123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := auth.Middleware()(h.PostLogin)(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestPostLoginRejectsInvalidRequestAndMissingAuth(t *testing.T) {
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiUserHandler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"account":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := h.PostLogin(e.NewContext(req, rec)); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("validation status=%d err=%v", rec.Code, err)
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"account":"tester","password":"password123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := h.PostLogin(e.NewContext(req, rec)); err == nil {
		t.Fatal("expected missing auth error")
	}
}

func TestPostLoginRejectsMalformedJSON(t *testing.T) {
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiUserHandler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"account":`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := h.PostLogin(e.NewContext(req, rec)); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}

func TestPostLoginReturnsErrorWhenActionLoginFails(t *testing.T) {
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiUserHandler{}
	auth := appAuth.New(&appAuth.Config[model.User, model.Session]{SecretKey: "test-secret", Resolver: &appAuth.Resolver[model.User, model.Session]{
		IsAccountExist: func(string) (bool, error) { return false, nil }, CreateUser: func(string, string) (*appAuth.User[model.User], error) {
			return &appAuth.User[model.User]{ID: "user-1"}, nil
		},
		GetUser: func(id string) (*appAuth.User[model.User], error) {
			return &appAuth.User[model.User]{ID: id}, nil
		}, GetUserByAccount: func(string) (*appAuth.User[model.User], error) {
			return nil, errors.New("db error")
		},
		GetSession: func(id string) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id}, nil
		}, CreateSession: func(id, userID string, expires time.Time) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id, UserID: userID, ExpiresAt: expires}, nil
		},
		UpdateSession: func(id string, expires time.Time, _ *model.Session) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id, ExpiresAt: expires}, nil
		}, DeleteSession: func(id string) (*appAuth.Session[model.Session], error) {
			return &appAuth.Session[model.Session]{ID: id}, nil
		},
	}})
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"account":"tester","password":"password123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := auth.Middleware()(h.PostLogin)(c)
	if err == nil {
		t.Fatal("expected error from PostLogin when ActionLogin fails, got nil")
	}
}
