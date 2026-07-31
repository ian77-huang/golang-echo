package admin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	appValidator "github.com/ian77-huang/golang-echo/pkg/validator"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Session{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func seedUser(t *testing.T, db *gorm.DB, account string) *model.User {
	t.Helper()
	now := time.Now()
	user := &model.User{Account: account, CreatedAt: &now, UpdatedAt: &now}
	if err := db.Create(user).Error; err != nil {
		t.Fatal(err)
	}
	return user
}

func newContext(t *testing.T, body string) (*echo.Context, *httptest.ResponseRecorder) {
	t.Helper()
	e := echo.New()
	e.Validator = appValidator.New()
	req := httptest.NewRequest(http.MethodPut, "/api/admin/user", bytes.NewBufferString(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestPutUserUpdatesUser(t *testing.T) {
	db := newTestDB(t)
	user := seedUser(t, db, "tester")

	h := &ApiAminHandler{DB: db}
	body := `{"id":"` + string(rune('0'+user.Id)) + `","account":"tester","isActive":true,"isAdmin":true}`
	c, rec := newContext(t, body)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: string(rune('0' + user.Id))})

	if err := h.PutUser(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}

	var saved model.User
	if err := db.First(&saved, user.Id).Error; err != nil {
		t.Fatal(err)
	}
	if !saved.IsAdmin || !saved.IsActive {
		t.Fatalf("expected user to be updated: %#v", saved)
	}
}

func TestPutUserDisablesUser(t *testing.T) {
	db := newTestDB(t)
	user := seedUser(t, db, "tester")

	h := &ApiAminHandler{DB: db}
	body := `{"id":"` + string(rune('0'+user.Id)) + `","account":"tester","isActive":false,"isAdmin":false}`
	c, rec := newContext(t, body)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "999"})

	if err := h.PutUser(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestPutUserRejectsSelfDisable(t *testing.T) {
	db := newTestDB(t)
	user := seedUser(t, db, "tester")

	h := &ApiAminHandler{DB: db}
	body := `{"id":"` + string(rune('0'+user.Id)) + `","account":"tester","isActive":false,"isAdmin":false}`
	c, _ := newContext(t, body)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: string(rune('0' + user.Id))})

	if err := h.PutUser(c); err == nil {
		t.Fatal("expected error disabling yourself")
	}
}

func TestPutUserRejectsInvalidRequest(t *testing.T) {
	db := newTestDB(t)
	h := &ApiAminHandler{DB: db}

	c, rec := newContext(t, `{"id":`)
	if err := h.PutUser(c); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("malformed json: status=%d err=%v", rec.Code, err)
	}

	c, rec = newContext(t, `{"id":"1","account":"","isActive":true,"isAdmin":false}`)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "1"})
	if err := h.PutUser(c); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("validation: status=%d err=%v", rec.Code, err)
	}

	c, rec = newContext(t, `{"id":"1","account":"tester","isActive":true,"isAdmin":false}`)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "1"})
	if err := h.PutUser(c); err == nil {
		t.Fatal("expected error for missing target user")
	}
}
