package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	appValidator "github.com/ian77-huang/golang-echo/pkg/validator"
	"github.com/labstack/echo/v5"
)

func TestPostChangeLangSetsCookie(t *testing.T) {
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiHandler{}
	req := httptest.NewRequest(http.MethodPost, "/api/lang", bytes.NewBufferString(`{"code":"en"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	if err := h.PostChangeLang(e.NewContext(req, rec)); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK || rec.Result().Cookies()[0].Value != "en" {
		t.Fatalf("status=%d cookies=%#v", rec.Code, rec.Result().Cookies())
	}
}

func TestPostChangeLangRejectsInvalidCode(t *testing.T) {
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiHandler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/lang", bytes.NewBufferString(`{"code":"fr"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := h.PostChangeLang(e.NewContext(req, rec)); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d", rec.Code)
	}
}

func TestPostChangeLangRejectsMalformedJSON(t *testing.T) {
	e := echo.New()
	e.Validator = appValidator.New()
	h := &ApiHandler{}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/lang", bytes.NewBufferString(`{"code":`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if err := h.PostChangeLang(e.NewContext(req, rec)); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d", rec.Code)
	}
}
