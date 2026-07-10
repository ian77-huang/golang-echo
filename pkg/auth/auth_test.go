package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func TestRegisterReturnsCustomAccountLookupError(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{
		SecretKey: "test-secret",
		Resolver: Resolver[struct{}, struct{}]{
			IsAccountExist: func(string) (bool, error) {
				return false, NewError("error.auth.AccountAlreadyExists", "custom lookup error")
			},
		},
	})

	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodPost, "/register", nil), httptest.NewRecorder())
	_, err := auth.Register(c, "tester", "password123")

	var fieldErr FieldError
	if !errors.As(err, &fieldErr) || fieldErr.Message != "custom lookup error" {
		t.Fatalf("expected the custom account lookup error, got %v", err)
	}
}

func TestRefreshSessionReturnsConfigurationErrorForMissingResolver(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret"})
	token, err := auth.createToken("session-1", time.Now().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: DEFAULT_SESSION_COOKIE_NAME, Value: token})
	c := e.NewContext(req, httptest.NewRecorder())

	_, err = auth.refreshSession(c)
	var fieldErr FieldError
	if !errors.As(err, &fieldErr) || fieldErr.Tag != "error.auth.InvalidConfiguration" {
		t.Fatalf("expected configuration error, got %v", err)
	}
}

func TestParseTokenRejectsOtherHMACAlgorithms(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret"})
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, JwtCustomClaims{
		ID:               "session-1",
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))},
	})
	tokenString, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := auth.parseToken(tokenString); err == nil {
		t.Fatal("expected HS512 token to be rejected")
	}
}
