package users

import (
	"bytes"
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

func TestPostRegisterUsesAuthFromMiddleware(t *testing.T) {
	t.Setenv("USERS_ACCOUNT_MIN_LENGTH", "6")
	t.Setenv("USERS_PASSWORD_MIN_LENGTH", "8")

	e := echo.New()
	e.Validator = appValidator.New()

	auth := appAuth.New(&appAuth.Config[userModel.User, sessionModel.Session]{
		SecretKey: "test-secret",
		Resolver: (&appAuth.Resolver[userModel.User, sessionModel.Session]{
			IsAccountExist: func(account string) (bool, error) {
				return false, nil
			},
			CreateUser: func(account string, password string) (*appAuth.User[userModel.User], error) {
				return &appAuth.User[userModel.User]{ID: "user-1"}, nil
			},
			CreateSession: func(id string, userID string, expiresAt time.Time) (*appAuth.Session[sessionModel.Session], error) {
				return &appAuth.Session[sessionModel.Session]{ID: id, ExpiresAt: expiresAt}, nil
			},
			GetSession: func(id string) (*appAuth.Session[sessionModel.Session], error) {
				return &appAuth.Session[sessionModel.Session]{ID: id}, nil
			},
			GetUser: func(id string) (*appAuth.User[userModel.User], error) {
				return &appAuth.User[userModel.User]{ID: id}, nil
			},
			GetUserByAccount: func(account string) (*appAuth.User[userModel.User], error) {
				return &appAuth.User[userModel.User]{ID: "user-1"}, nil
			},
			UpdateSession: func(id string, expiresAt time.Time, _ *sessionModel.Session) (*appAuth.Session[sessionModel.Session], error) {
				return &appAuth.Session[sessionModel.Session]{ID: id, ExpiresAt: expiresAt}, nil
			},
			DeleteSession: func(id string) (*appAuth.Session[sessionModel.Session], error) {
				return &appAuth.Session[sessionModel.Session]{ID: id}, nil
			},
		}),
	})

	body := bytes.NewBufferString(`{"account":"tester","password":"password123","confirmPassword":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/users/register", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := auth.Middleware()(PostRegister)
	if err := handler(c); err != nil {
		e.HTTPErrorHandler(c, err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, rec.Code, rec.Body.String())
	}
}

func TestPostRegisterRequiresAuthContext(t *testing.T) {
	e := echo.New()
	e.Validator = appValidator.New()
	body := bytes.NewBufferString(`{"account":"tester","password":"password123","confirmPassword":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/users/register", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, httptest.NewRecorder())
	if err := PostRegister(c); err == nil {
		t.Fatal("expected missing auth error")
	}
}

func TestPostRegisterRejectsShortValues(t *testing.T) {
	t.Setenv("USERS_ACCOUNT_MIN_LENGTH", "6")
	t.Setenv("USERS_PASSWORD_MIN_LENGTH", "8")
	e := echo.New()
	e.Validator = appValidator.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBufferString(`{"account":"short","password":"password123","confirmPassword":"password123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := PostRegister(e.NewContext(req, rec)); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("account validation: status=%d err=%v", rec.Code, err)
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBufferString(`{"account":"tester","password":"short","confirmPassword":"short"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := PostRegister(e.NewContext(req, rec)); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("password validation: status=%d err=%v", rec.Code, err)
	}
}

func TestPostRegisterRejectsMalformedJSONAndValidationMismatch(t *testing.T) {
	t.Setenv("USERS_ACCOUNT_MIN_LENGTH", "6")
	t.Setenv("USERS_PASSWORD_MIN_LENGTH", "8")
	e := echo.New()
	e.Validator = appValidator.New()
	for _, body := range []string{`{"account":`, `{"account":"tester","password":"password123","confirmPassword":"different"}`} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewBufferString(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		if err := PostRegister(e.NewContext(req, rec)); err != nil || rec.Code != http.StatusBadRequest {
			t.Fatalf("body=%s status=%d err=%v", body, rec.Code, err)
		}
	}
}
