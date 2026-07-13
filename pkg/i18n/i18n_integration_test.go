package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/labstack/echo/v5"
)

func testI18n(t *testing.T) *I18n {
	t.Helper()
	translator, err := New(Config{DefaultLang: "en", SupportedLanguageCodes: []string{"en", "zh-TW"}, MessageFS: fstest.MapFS{
		"messages.en.toml":    {Data: []byte("[welcome]\nother = 'Hello {{.Name}}'\n")},
		"messages.zh-TW.toml": {Data: []byte("[welcome]\nother = '你好 {{.Name}}'\n")},
	}, MessageFiles: []string{"messages.en.toml", "messages.zh-TW.toml"}})
	if err != nil {
		t.Fatal(err)
	}
	return translator
}

func TestNewValidatesConfiguration(t *testing.T) {
	if _, err := New(Config{}); err == nil {
		t.Fatal("expected configuration error")
	}
	if _, err := New(Config{DefaultLang: "en", SupportedLanguageCodes: []string{"en"}, MessageFiles: []string{"x.toml"}}); err == nil {
		t.Fatal("expected message fs error")
	}
	if _, err := New(Config{DefaultLang: "!", SupportedLanguageCodes: []string{"en"}}); err == nil {
		t.Fatal("expected invalid default language error")
	}
	if _, err := New(Config{DefaultLang: "en", SupportedLanguageCodes: []string{"!"}}); err == nil {
		t.Fatal("expected invalid supported language error")
	}
	if _, err := New(Config{DefaultLang: "en"}); err == nil {
		t.Fatal("expected missing supported language error")
	}
	if _, err := New(Config{DefaultLang: "en", SupportedLanguageCodes: []string{"en"}, MessageFS: fstest.MapFS{}, MessageFiles: []string{"missing.toml"}}); err == nil {
		t.Fatal("expected message load error")
	}
}

func TestMiddlewareSelectsLanguageAndTranslation(t *testing.T) {
	translator := testI18n(t)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "zh-TW, en;q=0.8")
	c := e.NewContext(req, httptest.NewRecorder())
	err := translator.Middleware()(func(c *echo.Context) error {
		if c.Get("lang") != "zh-TW" || T(c, "welcome", KV("Name", "Yien")) != "你好 Yien" {
			t.Fatalf("language=%v translation=%q", c.Get("lang"), T(c, "welcome", KV("Name", "Yien")))
		}
		return nil
	})(c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTemplateFuncsUseFallbackWithoutContext(t *testing.T) {
	funcs := testI18n(t).TemplateFuncs(nil, nil)
	if funcs["t"].(func(string, ...any) string)("missing") != "missing" {
		t.Fatal("expected fallback key")
	}
	if T(echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder()), "missing") != "localizer does not exist." {
		t.Fatal("expected missing localizer message")
	}
}

func TestTemplateFuncsUseRequestLocalizer(t *testing.T) {
	translator := testI18n(t)
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	if err := translator.Middleware()(func(*echo.Context) error { return nil })(c); err != nil {
		t.Fatal(err)
	}
	funcs := translator.TemplateFuncs(c, nil)
	if funcs["t"].(func(string, ...any) string)("welcome", KV("Name", "Yien")) != "Hello Yien" {
		t.Fatal("expected translated template function")
	}
	if _, ok := translator.Localizer(c); !ok {
		t.Fatal("expected localizer")
	}
}

func TestTemplateFuncsFallbackForMissingLocalizerAndMessage(t *testing.T) {
	translator := testI18n(t)
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	if translator.TemplateFuncs(c, nil)["t"].(func(string, ...any) string)("missing") != "missing" {
		t.Fatal("expected fallback template function")
	}
	if err := translator.Middleware()(func(c *echo.Context) error {
		if T(c, "missing") != "missing" {
			t.Fatal("expected missing message id")
		}
		return nil
	})(c); err != nil {
		t.Fatal(err)
	}
}

func TestMiddlewarePrefersLanguageCookie(t *testing.T) {
	translator := testI18n(t)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "en")
	req.AddCookie(&http.Cookie{Name: "lang", Value: "zh-TW"})
	c := e.NewContext(req, httptest.NewRecorder())
	if err := translator.Middleware()(func(c *echo.Context) error {
		if c.Get("lang") != "zh-TW" {
			t.Fatalf("got %v", c.Get("lang"))
		}
		return nil
	})(c); err != nil {
		t.Fatal(err)
	}
}
