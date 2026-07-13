package users

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sessionModel "github.com/ian77-huang/golang-echo/internal/models/session"
	userModel "github.com/ian77-huang/golang-echo/internal/models/users"
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
	auth := appAuth.New(&appAuth.Config[userModel.User, sessionModel.Session]{SecretKey: "test-secret", Resolver: &appAuth.Resolver[userModel.User, sessionModel.Session]{
		IsAccountExist: func(string) (bool, error) { return false, nil }, CreateUser: func(string, string) (*appAuth.User[userModel.User], error) {
			return &appAuth.User[userModel.User]{ID: "user-1"}, nil
		},
		GetUser: func(id string) (*appAuth.User[userModel.User], error) {
			return &appAuth.User[userModel.User]{ID: id}, nil
		}, GetUserByAccount: func(string) (*appAuth.User[userModel.User], error) {
			return &appAuth.User[userModel.User]{ID: "user-1", Password: hash}, nil
		},
		GetSession: func(id string) (*appAuth.Session[sessionModel.Session], error) {
			return &appAuth.Session[sessionModel.Session]{ID: id}, nil
		}, CreateSession: func(id, userID string, expires time.Time) (*appAuth.Session[sessionModel.Session], error) {
			return &appAuth.Session[sessionModel.Session]{ID: id, UserID: userID, ExpiresAt: expires}, nil
		},
		UpdateSession: func(id string, expires time.Time, _ *sessionModel.Session) (*appAuth.Session[sessionModel.Session], error) {
			return &appAuth.Session[sessionModel.Session]{ID: id, ExpiresAt: expires}, nil
		}, DeleteSession: func(id string) (*appAuth.Session[sessionModel.Session], error) {
			return &appAuth.Session[sessionModel.Session]{ID: id}, nil
		},
	}})
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"account":"tester","password":"password123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := auth.Middleware()(PostLogin)(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestPostLoginRejectsInvalidRequestAndMissingAuth(t *testing.T) {
	e := echo.New()
	e.Validator = appValidator.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"account":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := PostLogin(e.NewContext(req, rec)); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("validation status=%d err=%v", rec.Code, err)
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"account":"tester","password":"password123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := PostLogin(e.NewContext(req, rec)); err == nil {
		t.Fatal("expected missing auth error")
	}
}

func TestPostLoginRejectsMalformedJSON(t *testing.T) {
	e := echo.New()
	e.Validator = appValidator.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"account":`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := PostLogin(e.NewContext(req, rec)); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}

func TestPostLoginReturnsErrorWhenActionLoginFails(t *testing.T) {
	e := echo.New()
	e.Validator = appValidator.New()
	auth := appAuth.New(&appAuth.Config[userModel.User, sessionModel.Session]{SecretKey: "test-secret", Resolver: &appAuth.Resolver[userModel.User, sessionModel.Session]{
		IsAccountExist: func(string) (bool, error) { return false, nil }, CreateUser: func(string, string) (*appAuth.User[userModel.User], error) {
			return &appAuth.User[userModel.User]{ID: "user-1"}, nil
		},
		GetUser: func(id string) (*appAuth.User[userModel.User], error) {
			return &appAuth.User[userModel.User]{ID: id}, nil
		}, GetUserByAccount: func(string) (*appAuth.User[userModel.User], error) {
			return nil, errors.New("db error")
		},
		GetSession: func(id string) (*appAuth.Session[sessionModel.Session], error) {
			return &appAuth.Session[sessionModel.Session]{ID: id}, nil
		}, CreateSession: func(id, userID string, expires time.Time) (*appAuth.Session[sessionModel.Session], error) {
			return &appAuth.Session[sessionModel.Session]{ID: id, UserID: userID, ExpiresAt: expires}, nil
		},
		UpdateSession: func(id string, expires time.Time, _ *sessionModel.Session) (*appAuth.Session[sessionModel.Session], error) {
			return &appAuth.Session[sessionModel.Session]{ID: id, ExpiresAt: expires}, nil
		}, DeleteSession: func(id string) (*appAuth.Session[sessionModel.Session], error) {
			return &appAuth.Session[sessionModel.Session]{ID: id}, nil
		},
	}})
	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewBufferString(`{"account":"tester","password":"password123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := auth.Middleware()(PostLogin)(c)
	if err == nil {
		t.Fatal("expected error from PostLogin when ActionLogin fails, got nil")
	}
}
