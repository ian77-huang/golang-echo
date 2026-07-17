package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestNewServerBuildsRunnableApplication(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	t.Setenv("SECRET_KEY", "test-secret")
	t.Setenv("DATABASE_PATH", filepath.Join(t.TempDir(), "app.db"))
	t.Setenv("USER_ACCOUNT_MIN_LENGTH", "6")
	t.Setenv("USER_PASSWORD_MIN_LENGTH", "8")
	e, err := newServer()
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/ping", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("ping status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestNewServerRedirectsGuestFromProtectedAPIs(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	t.Setenv("SECRET_KEY", "test-secret")
	t.Setenv("DATABASE_PATH", filepath.Join(t.TempDir(), "app.db"))
	t.Setenv("USER_ACCOUNT_MIN_LENGTH", "6")
	t.Setenv("USER_PASSWORD_MIN_LENGTH", "8")
	e, err := newServer()
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/user/profile"},
		{http.MethodPut, "/api/user/profile"},
		{http.MethodPost, "/api/user/profile/avatar"},
		{http.MethodPost, "/api/user/reset-password"},
	} {
		t.Run(testCase.method+" "+testCase.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, httptest.NewRequest(testCase.method, testCase.path, nil))
			if rec.Code != http.StatusFound || rec.Header().Get("Location") != "/user/login" {
				t.Fatalf("status=%d location=%q", rec.Code, rec.Header().Get("Location"))
			}
		})
	}
}
