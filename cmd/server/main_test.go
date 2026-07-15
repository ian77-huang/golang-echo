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
