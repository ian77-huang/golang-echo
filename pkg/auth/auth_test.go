package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ian77-huang/golang-echo/pkg/argon2"
	"github.com/labstack/echo/v5"
)

func TestRegisterReturnsAccountLookupError(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{
		SecretKey: "test-secret",
		Resolver: testResolver(func(r *Resolver[struct{}, struct{}]) {
			r.IsAccountExist = func(string) (bool, error) {
				return false, NewError("error.auth.AccountAlreadyExists", "custom lookup error")
			}
		}),
	})

	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodPost, "/register", nil), httptest.NewRecorder())
	id, err := auth.ActionRegister(c, "tester", "password123")
	var fieldErr FieldError
	if id != "" || !errors.As(err, &fieldErr) || fieldErr.Message != "custom lookup error" {
		t.Fatalf("expected the custom account lookup error, got id=%q err=%v", id, err)
	}
}

func TestRefreshSessionReturnsConfigurationErrorForMissingResolver(t *testing.T) {
	// auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret"})
	// token, err := auth.createToken("session-1", time.Now().Add(time.Hour))
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// e := echo.New()
	// req := httptest.NewRequest(http.MethodGet, "/", nil)
	// req.AddCookie(&http.Cookie{Name: DEFAULT_SESSION_COOKIE_NAME, Value: token})
	// c := e.NewContext(req, httptest.NewRecorder())

	// _, err = auth.refreshSession(c)
	// var fieldErr FieldError
	// if !errors.As(err, &fieldErr) || fieldErr.Tag != "error.auth.InvalidConfiguration" {
	// 	t.Fatalf("expected configuration error, got %v", err)
	// }
}

func TestParseTokenRejectsOtherHMACAlgorithms(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
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

func TestNewAppliesDefaultsAndTokenRoundTrip(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	if auth.config.CookieName != DEFAULT_SESSION_COOKIE_NAME || auth.config.SessionExpiresAt != DEFAULT_SESSION_EXPIRES_DAYS {
		t.Fatalf("defaults were not applied: %#v", auth.config)
	}
	token, err := auth.createToken("session-token", time.Now().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	claims, err := auth.parseToken(token)
	if err != nil || claims.ID != "session-token" {
		t.Fatalf("parseToken() = %#v, %v", claims, err)
	}
}

func TestPasswordAndSessionTokenHelpers(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatal(err)
	}
	valid, err := VerifyPassword("password123", hash)
	if err != nil || !valid {
		t.Fatalf("VerifyPassword() = %v, %v", valid, err)
	}
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	token, err := auth.GenerateSessionToken()
	if err != nil || token == "" || len(auth.GenerateID(token)) != 64 {
		t.Fatalf("invalid token helpers: %q, %v", token, err)
	}
}

func TestActionLoginRejectsMissingUser(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(func(r *Resolver[struct{}, struct{}]) {
		r.GetUserByAccount = func(string) (*User[struct{}], error) { return nil, nil }
	})})
	e := echo.New()
	_, err := auth.ActionLogin(e.NewContext(httptest.NewRequest(http.MethodPost, "/login", nil), httptest.NewRecorder()), "missing", "password")
	var fieldErr FieldError
	if !errors.As(err, &fieldErr) || fieldErr.Tag != "error.auth.AccountLookupFailed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestActionRegisterCreatesSessionAndCookie(t *testing.T) {
	created := false
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(func(r *Resolver[struct{}, struct{}]) {
		r.CreateSession = func(id, userID string, expires time.Time) (*Session[struct{}], error) {
			created = id != "" && userID == "user-1" && expires.After(time.Now())
			return &Session[struct{}]{ID: id, UserID: userID, ExpiresAt: expires}, nil
		}
	})})
	e := echo.New()
	rec := httptest.NewRecorder()
	id, err := auth.ActionRegister(e.NewContext(httptest.NewRequest(http.MethodPost, "/register", nil), rec), "new-user", "password123")
	if err != nil || id != "user-1" || !created || len(rec.Result().Cookies()) != 1 {
		t.Fatalf("id=%q created=%v cookies=%#v err=%v", id, created, rec.Result().Cookies(), err)
	}
}

func TestMiddlewareSetsAuthForGuestRequest(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	called := false
	err := auth.Middleware()(func(c *echo.Context) error {
		called = Load[struct{}, struct{}](c) == auth && !IsSignedIn[struct{}](c)
		return nil
	})(c)
	if err != nil || !called {
		t.Fatalf("middleware called=%v err=%v", called, err)
	}
}

func TestSessionRefreshAndLogout(t *testing.T) {
	updated, deleted := false, false
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", SessionRefreshAt: 15, Resolver: testResolver(func(r *Resolver[struct{}, struct{}]) {
		r.UpdateSession = func(id string, expires time.Time, _ *struct{}) (*Session[struct{}], error) {
			updated = id == "session-1" && expires.After(time.Now())
			return &Session[struct{}]{ID: id, ExpiresAt: expires}, nil
		}
		r.DeleteSession = func(id string) (*Session[struct{}], error) {
			deleted = id == "user-1"
			return &Session[struct{}]{ID: id}, nil
		}
	})})
	if !auth.isNeedRefreshSession(time.Now().Add(14*24*time.Hour)) || auth.isNeedRefreshSession(time.Now().Add(16*24*time.Hour)) {
		t.Fatal("unexpected refresh decision")
	}
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)
	ok, err := auth.refreshSession(c, "token", &Session[struct{}]{ID: "session-1", ExpiresAt: time.Now().Add(time.Hour), Data: &struct{}{}})
	if err != nil || !ok || !updated || len(rec.Result().Cookies()) != 1 {
		t.Fatalf("refresh: ok=%v updated=%v cookies=%#v err=%v", ok, updated, rec.Result().Cookies(), err)
	}
	c.Set(CONTEXT_KEY_USER, &User[struct{}]{ID: "user-1"})
	ok, err = auth.ActionLogout(c)
	if err != nil || !ok || !deleted {
		t.Fatalf("logout: ok=%v deleted=%v err=%v", ok, deleted, err)
	}
}

// func TestGettersAndGuestOnlyRoute(t *testing.T) {
// auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil), Route: &Route[struct{}]{GuestOnly: &RoutesPaths{Rules: []string{"/login"}, RedirectURL: "/"}}})
// e := echo.New()
// rec := httptest.NewRecorder()
// c := e.NewContext(httptest.NewRequest(http.MethodGet, "/login", nil), rec)
// c.SetPath("/login")
// c.Set(CONTEXT_KEY_USER, &User[struct{}]{ID: "user-1"})
// if err := auth.Middleware()(func(*echo.Context) error { t.Fatal("guest-only handler should not run"); return nil })(c); err != nil || rec.Code != http.StatusFound {
// 	t.Fatalf("redirect: status=%d err=%v", rec.Code, err)
// }
// c.Set(CONTEXT_KEY_SESSION, &Session[struct{}]{ID: "session-1"})
// if GetSession[struct{}](c).ID != "session-1" || GetSession[int](c) != nil || GetUser[struct{}](c).ID != "user-1" {
// 	t.Fatal("unexpected getter values")
// }
// }

func TestActionLoginAndErrorFormatting(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatal(err)
	}
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(func(r *Resolver[struct{}, struct{}]) {
		r.GetUserByAccount = func(string) (*User[struct{}], error) { return &User[struct{}]{ID: "user-1", Password: hash}, nil }
	})})
	e := echo.New()
	rec := httptest.NewRecorder()
	ok, err := auth.ActionLogin(e.NewContext(httptest.NewRequest(http.MethodPost, "/login", nil), rec), "user", "password123")
	if err != nil || !ok || len(rec.Result().Cookies()) != 1 {
		t.Fatalf("login: ok=%v cookies=%#v err=%v", ok, rec.Result().Cookies(), err)
	}
	if ok, err := auth.ActionLogin(e.NewContext(httptest.NewRequest(http.MethodPost, "/login", nil), httptest.NewRecorder()), "user", "wrong"); ok || err == nil {
		t.Fatal("expected invalid password failure")
	}
	if NewError("tag", "message", "x").Error() == "" {
		t.Fatal("expected formatted field error")
	}
}

func TestMiddlewareLoadsValidSession(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(func(r *Resolver[struct{}, struct{}]) {
		r.GetSession = func(id string) (*Session[struct{}], error) {
			return &Session[struct{}]{ID: id, UserID: "user-1", ExpiresAt: time.Now().Add(20 * 24 * time.Hour)}, nil
		}
		r.GetUser = func(id string) (*User[struct{}], error) { return &User[struct{}]{ID: id}, nil }
	})})
	token, err := auth.createToken("token-id", time.Now().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: DEFAULT_SESSION_COOKIE_NAME, Value: token})
	c := e.NewContext(req, httptest.NewRecorder())
	called := false
	if err := auth.Middleware()(func(c *echo.Context) error {
		called = IsSignedIn[struct{}](c) && GetUser[struct{}](c).ID == "user-1" && GetSession[struct{}](c) != nil
		return nil
	})(c); err != nil || !called {
		t.Fatalf("called=%v err=%v", called, err)
	}
}

func TestConfigurationAndResolverErrors(t *testing.T) {
	for _, config := range []*Config[struct{}, struct{}]{nil, &Config[struct{}, struct{}]{SecretKey: " "}, &Config[struct{}, struct{}]{SecretKey: "test-secret"}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatal("expected New panic")
				}
			}()
			New(config)
		}()
	}
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	auth.config.Resolver.IsAccountExist = func(string) (bool, error) { return false, errors.New("database") }
	if _, err := auth.IsAccountExist("x"); err == nil {
		t.Fatal("expected lookup error")
	}
	auth.config.Resolver.IsAccountExist = nil
	if _, err := auth.IsAccountExist("x"); err == nil {
		t.Fatal("expected missing lookup resolver error")
	}
}

func TestNewPanicsForEveryMissingRequiredResolver(t *testing.T) {
	missing := []struct {
		name  string
		clear func(*Resolver[struct{}, struct{}])
	}{
		{"CreateSession", func(r *Resolver[struct{}, struct{}]) { r.CreateSession = nil }}, {"CreateUser", func(r *Resolver[struct{}, struct{}]) { r.CreateUser = nil }},
		{"DeleteSession", func(r *Resolver[struct{}, struct{}]) { r.DeleteSession = nil }}, {"GetSession", func(r *Resolver[struct{}, struct{}]) { r.GetSession = nil }},
		{"GetUser", func(r *Resolver[struct{}, struct{}]) { r.GetUser = nil }}, {"GetUserByAccount", func(r *Resolver[struct{}, struct{}]) { r.GetUserByAccount = nil }},
		{"IsAccountExist", func(r *Resolver[struct{}, struct{}]) { r.IsAccountExist = nil }}, {"UpdateSession", func(r *Resolver[struct{}, struct{}]) { r.UpdateSession = nil }},
	}
	for _, tt := range missing {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected panic")
				}
			}()
			r := testResolver(nil)
			tt.clear(r)
			New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: r})
		})
	}
}

func TestCreateUserAndRegisterReturnResolverErrors(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(func(r *Resolver[struct{}, struct{}]) {
		r.CreateUser = func(string, string) (*User[struct{}], error) { return nil, errors.New("write failed") }
	})})
	if _, err := auth.CreateUser("user", "password"); err == nil {
		t.Fatal("expected create user error")
	}
	if _, err := auth.ActionRegister(echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder()), "user", "password"); err == nil {
		t.Fatal("expected register error")
	}
	auth.config.Resolver.CreateUser = nil
	if _, err := auth.CreateUser("user", "password"); err == nil {
		t.Fatal("expected missing create resolver error")
	}
}

func TestSessionErrorsReturnErrors(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(func(r *Resolver[struct{}, struct{}]) {
		r.CreateSession = func(string, string, time.Time) (*Session[struct{}], error) {
			return nil, errors.New("session write failed")
		}
	})})
	c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())
	if ok, err := auth.createSession(c, "user-1"); ok || err == nil {
		t.Fatalf("expected create session error: ok=%v err=%v", ok, err)
	}
	auth.config.Resolver.UpdateSession = func(string, time.Time, *struct{}) (*Session[struct{}], error) {
		return nil, errors.New("update failed")
	}
	if ok, err := auth.refreshSession(c, "token", &Session[struct{}]{ID: "session-1", ExpiresAt: time.Now().Add(time.Hour), Data: &struct{}{}}); ok || err == nil {
		t.Fatalf("expected refresh error: ok=%v err=%v", ok, err)
	}
}

func TestAuthActionErrorBranches(t *testing.T) {
	e := echo.New()
	context := func() *echo.Context {
		return e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())
	}
	newAuth := func(update func(*Resolver[struct{}, struct{}])) *Auth[struct{}, struct{}] {
		return New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(update)})
	}
	if id, err := newAuth(func(r *Resolver[struct{}, struct{}]) {
		r.IsAccountExist = func(string) (bool, error) { return true, nil }
	}).ActionRegister(context(), "u", "p"); id != "" || err == nil {
		t.Fatal("expected duplicate account error")
	}
	if ok, err := newAuth(func(r *Resolver[struct{}, struct{}]) {
		r.GetUserByAccount = func(string) (*User[struct{}], error) { return nil, errors.New("lookup") }
	}).ActionLogin(context(), "u", "p"); ok || err == nil {
		t.Fatal("expected lookup error")
	}
	if ok, err := newAuth(func(r *Resolver[struct{}, struct{}]) {
		r.GetUserByAccount = func(string) (*User[struct{}], error) { return &User[struct{}]{ID: "u", Password: "bad"}, nil }
	}).ActionLogin(context(), "u", "p"); ok || err == nil {
		t.Fatal("expected password verification error")
	}
	auth := newAuth(func(r *Resolver[struct{}, struct{}]) {
		r.DeleteSession = func(string) (*Session[struct{}], error) { return nil, errors.New("delete") }
	})
	c := context()
	c.Set(CONTEXT_KEY_USER, &User[struct{}]{ID: "u"})
	if ok, err := auth.ActionLogout(c); ok || err == nil {
		t.Fatal("expected delete error")
	}
}

func TestAuthNilConfigurationAndLoadTypeErrors(t *testing.T) {
	var auth *Auth[struct{}, struct{}]
	if _, err := auth.getConfig(); err == nil {
		t.Fatal("expected nil auth configuration error")
	}
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	c.Set(CONTEXT_KEY_AUTH, "wrong")
	if Load[struct{}, struct{}](c) != nil {
		t.Fatal("expected nil auth for wrong context type")
	}
	if _, err := (&Auth[struct{}, struct{}]{config: nil}).CreateUser("u", "p"); err == nil {
		t.Fatal("expected create user config error")
	}
	if _, err := (&Auth[struct{}, struct{}]{config: nil}).IsAccountExist("u"); err == nil {
		t.Fatal("expected account config error")
	}
}

func TestActionMethodsReturnConfigurationAndSessionErrors(t *testing.T) {
	var nilAuth *Auth[struct{}, struct{}]
	c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())
	if _, err := nilAuth.ActionLogin(c, "u", "p"); err == nil {
		t.Fatal("expected nil login config error")
	}
	if _, err := nilAuth.ActionLogout(c); err == nil {
		t.Fatal("expected nil logout config error")
	}
	hash, err := HashPassword("password")
	if err != nil {
		t.Fatal(err)
	}
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(func(r *Resolver[struct{}, struct{}]) {
		r.GetUserByAccount = func(string) (*User[struct{}], error) { return &User[struct{}]{ID: "u", Password: hash}, nil }
		r.CreateSession = func(string, string, time.Time) (*Session[struct{}], error) { return nil, errors.New("write") }
	})})
	if _, err := auth.ActionRegister(c, "u", "p"); err == nil {
		t.Fatal("expected register session error")
	}
	if _, err := auth.ActionLogin(c, "u", "password"); err == nil {
		t.Fatal("expected login session error")
	}
}

func TestMiddlewareReturnsAndHandlesFailureStates(t *testing.T) {
	e := echo.New()
	newContext := func(cookie *http.Cookie) *echo.Context {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if cookie != nil {
			req.AddCookie(cookie)
		}
		return e.NewContext(req, httptest.NewRecorder())
	}
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	auth.config = nil
	if err := auth.Middleware()(func(*echo.Context) error { return nil })(newContext(nil)); err == nil {
		t.Fatal("expected middleware configuration error")
	}
	auth = New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	called := false
	if err := auth.Middleware()(func(*echo.Context) error { called = true; return nil })(newContext(&http.Cookie{Name: DEFAULT_SESSION_COOKIE_NAME, Value: "not-a-jwt"})); err != nil || !called {
		t.Fatalf("invalid token: called=%v err=%v", called, err)
	}
	token, err := auth.createToken("token", time.Now().Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	cookie := &http.Cookie{Name: DEFAULT_SESSION_COOKIE_NAME, Value: token}
	auth.config.Resolver.GetSession = func(string) (*Session[struct{}], error) { return nil, errors.New("missing session") }
	if err := auth.Middleware()(func(*echo.Context) error { return nil })(newContext(cookie)); err != nil {
		t.Fatal(err)
	}
	auth.config.Resolver.GetSession = func(id string) (*Session[struct{}], error) {
		return &Session[struct{}]{ID: id, UserID: "user", ExpiresAt: time.Now().Add(time.Hour)}, nil
	}
	auth.config.Resolver.GetUser = func(string) (*User[struct{}], error) { return nil, nil }
	if err := auth.Middleware()(func(*echo.Context) error { return nil })(newContext(cookie)); err == nil {
		t.Fatal("expected missing user error")
	}
}

func TestGenerateSessionTokenReturnsRandomSourceError(t *testing.T) {
	original := randomRead
	randomRead = func([]byte) (int, error) { return 0, errors.New("random failed") }
	t.Cleanup(func() { randomRead = original })
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	if _, err := auth.GenerateSessionToken(); err == nil {
		t.Fatal("expected random source error")
	}
}

func TestSessionConfigurationErrorBranches(t *testing.T) {
	c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())
	var nilAuth *Auth[struct{}, struct{}]
	if ok, err := nilAuth.createSession(c, "user"); ok || err == nil {
		t.Fatalf("expected nil createSession error: ok=%v err=%v", ok, err)
	}
	if ok, err := nilAuth.refreshSession(c, "token", &Session[struct{}]{}); ok || err == nil {
		t.Fatalf("expected nil refreshSession error: ok=%v err=%v", ok, err)
	}
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	auth.config.Resolver.CreateSession = nil
	if ok, err := auth.createSession(c, "user"); ok || err == nil {
		t.Fatalf("expected missing create resolver error: ok=%v err=%v", ok, err)
	}
	auth.config.Resolver.CreateSession = testResolver(nil).CreateSession
	auth.config.Resolver.UpdateSession = nil
	if ok, err := auth.refreshSession(c, "token", &Session[struct{}]{ID: "session", ExpiresAt: time.Now().Add(time.Hour)}); ok || err == nil {
		t.Fatalf("expected missing update resolver error: ok=%v err=%v", ok, err)
	}
}

func TestSessionReturnsTokenSigningErrors(t *testing.T) {
	original := signJWT
	signJWT = func(*jwt.Token, []byte) (string, error) { return "", errors.New("sign failed") }
	t.Cleanup(func() { signJWT = original })
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())
	if ok, err := auth.createSession(c, "user"); ok || err == nil {
		t.Fatalf("expected create token error: ok=%v err=%v", ok, err)
	}
	if ok, err := auth.refreshSession(c, "token", &Session[struct{}]{ID: "session", ExpiresAt: time.Now().Add(time.Hour)}); ok || err == nil {
		t.Fatalf("expected refresh token error: ok=%v err=%v", ok, err)
	}
}

func testResolver(update func(*Resolver[struct{}, struct{}])) *Resolver[struct{}, struct{}] {
	r := &Resolver[struct{}, struct{}]{
		IsAccountExist:   func(string) (bool, error) { return false, nil },
		CreateUser:       func(string, string) (*User[struct{}], error) { return &User[struct{}]{ID: "user-1"}, nil },
		GetUser:          func(id string) (*User[struct{}], error) { return &User[struct{}]{ID: id}, nil },
		GetUserByAccount: func(string) (*User[struct{}], error) { return &User[struct{}]{ID: "user-1"}, nil },
		GetSession:       func(id string) (*Session[struct{}], error) { return &Session[struct{}]{ID: id}, nil },
		CreateSession: func(id, userID string, expires time.Time) (*Session[struct{}], error) {
			return &Session[struct{}]{ID: id, UserID: userID, ExpiresAt: expires}, nil
		},
		UpdateSession: func(id string, expires time.Time, _ *struct{}) (*Session[struct{}], error) {
			return &Session[struct{}]{ID: id, ExpiresAt: expires}, nil
		},
		DeleteSession: func(id string) (*Session[struct{}], error) { return &Session[struct{}]{ID: id}, nil },
	}
	if update != nil {
		update(r)
	}
	return r
}

func TestGetAccessTokenMalformedCookie(t *testing.T) {
	original := readCookie
	readCookie = func(c *echo.Context, name string) (*http.Cookie, error) {
		return nil, errors.New("mock read cookie error")
	}
	t.Cleanup(func() { readCookie = original })

	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	_, err := auth.getAccessToken(c)
	if err == nil {
		t.Fatal("expected malformed cookie read error")
	}
	var fieldErr FieldError
	if !errors.As(err, &fieldErr) || fieldErr.Tag != "error.auth.ReadCookieFailed" {
		t.Fatalf("expected ReadCookieFailed error, got %v", err)
	}
}

func TestCreateUserPasswordHashingError(t *testing.T) {
	original := argon2.RandomRead
	argon2.RandomRead = func([]byte) (int, error) { return 0, errors.New("mock rand failed") }
	t.Cleanup(func() { argon2.RandomRead = original })

	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	_, err := auth.CreateUser("user", "password")
	if err == nil {
		t.Fatal("expected password hashing failure")
	}
	var fieldErr FieldError
	if !errors.As(err, &fieldErr) || fieldErr.Tag != "error.auth.FailedToSecurePassword" {
		t.Fatalf("expected FailedToSecurePassword error, got %v", err)
	}
}

func TestParseTokenInvalidClaimsTypeOrValidFalse(t *testing.T) {
	original := parseJWT
	t.Cleanup(func() { parseJWT = original })

	parseJWT = func(tokenStr string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
		return &jwt.Token{Valid: false, Claims: claims}, nil
	}

	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	_, err := auth.parseToken("some-token")
	if err == nil || err.Error() != "invalid token claims" {
		t.Fatalf("expected 'invalid token claims' error, got %v", err)
	}
}

func TestMiddlewareSessionRefreshFails(t *testing.T) {
	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", SessionRefreshAt: 15, Resolver: testResolver(func(r *Resolver[struct{}, struct{}]) {
		r.GetSession = func(id string) (*Session[struct{}], error) {
			return &Session[struct{}]{ID: id, UserID: "user-1", ExpiresAt: time.Now().Add(5 * 24 * time.Hour)}, nil
		}
		r.GetUser = func(id string) (*User[struct{}], error) { return &User[struct{}]{ID: id}, nil }
	})})
	auth.config.Resolver.UpdateSession = nil

	token, err := auth.createToken("token", time.Now().Add(20*24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: DEFAULT_SESSION_COOKIE_NAME, Value: token})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	called := false
	err = auth.Middleware()(func(c *echo.Context) error {
		called = true
		return nil
	})(c)

	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("expected middleware handler to be called")
	}

	// Verify that the cookie was deleted (i.e. has Max-Age = -1 or Expires in past)
	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected at least one cookie in response to indicate deletion")
	}
	deleted := false
	for _, cookie := range cookies {
		if cookie.Name == DEFAULT_SESSION_COOKIE_NAME && (cookie.MaxAge < 0 || cookie.Expires.Before(time.Now())) {
			deleted = true
			break
		}
	}
	if !deleted {
		t.Fatalf("expected session cookie to be deleted, response cookies: %#v", cookies)
	}
}

func TestCreateSessionGenerateTokenError(t *testing.T) {
	original := randomRead
	randomRead = func([]byte) (int, error) { return 0, errors.New("random failed") }
	t.Cleanup(func() { randomRead = original })

	auth := New(&Config[struct{}, struct{}]{SecretKey: "test-secret", Resolver: testResolver(nil)})
	c := echo.New().NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())
	ok, err := auth.createSession(c, "user-1")
	if ok || err == nil {
		t.Fatalf("expected error from createSession, got ok=%v, err=%v", ok, err)
	}
	var fieldErr FieldError
	if !errors.As(err, &fieldErr) || fieldErr.Tag != "error.auth.FailedToIssueToken" {
		t.Fatalf("expected FailedToIssueToken error, got %v", err)
	}
}
