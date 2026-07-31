package config

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDBOpensSqliteDatabase(t *testing.T) {
	db, err := DB(t.TempDir() + "/config.db")
	if err != nil || db == nil {
		t.Fatalf("db=%v err=%v", db, err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("ping: %v", err)
	}
	sqlDB.Close()
}

func TestDBReturnsErrorForMissingDirectory(t *testing.T) {
	if db, err := DB("/no_such_dir_xyz/app.db"); err == nil {
		t.Fatalf("expected error, got db %v", db)
	}
}

func TestStoreBuildsRedisBackedServer(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(mr.Close)
	ss := Store("redis://" + mr.Addr())
	if ss == nil {
		t.Fatal("store is nil")
	}
	if err := ss.Set("k", "v", time.Minute); err != nil {
		t.Fatalf("set: %v", err)
	}
	var got string
	if err := ss.GetByte("k", &got); err != nil || got != "v" {
		t.Fatalf("get: got=%q err=%v", got, err)
	}
}

func TestSetMenusAdmin(t *testing.T) {
	menus := SetMenusAdmin(MenuRules{T: func(id string, _ ...any) string { return id }})
	if len(menus) != 2 || menus[1].Url != "/admin/user" {
		t.Fatalf("unexpected menus: %#v", menus)
	}
}

func TestSetMenusUsersAllBranches(t *testing.T) {
	tf := func(id string, _ ...any) string { return id }
	for _, tt := range []struct {
		path     string
		signedIn bool
		admin    bool
		want     int
	}{
		{"/user/login", true, false, 0},
		{"/user/register", true, false, 0},
		{"/user", true, false, 2},
		{"/user", true, true, 3},
		{"/user/profile", true, false, 1},
		{"/user/profile", true, true, 2},
		{"/user/reset-password", true, false, 1},
		{"/user/reset-password", true, true, 2},
		{"/other", true, false, 2},
		{"/other", true, true, 3},
	} {
		got := SetMenusUsers(MenuUsersRules{Path: tt.path, IsSignedIn: tt.signedIn, IsAdmin: tt.admin, T: tf})
		if len(got) != tt.want {
			t.Fatalf("%s signedIn=%v admin=%v: got %d %#v", tt.path, tt.signedIn, tt.admin, len(got), got)
		}
	}
}

func TestAuthWithStoreCachesAndDeletesSessions(t *testing.T) {
	t.Setenv("SECRET_KEY", "test-secret")
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Session{}); err != nil {
		t.Fatal(err)
	}
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(mr.Close)
	ss := Store("redis://" + mr.Addr())

	auth := Auth(&AuthParameter{
		UserService:    service.NewUserService(db),
		SessionService: service.NewSessionService(db),
	}, ss)
	e := echo.New()
	rec := httptest.NewRecorder()
	id, err := auth.ActionRegister(e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), rec), "store-user", "password")
	if err != nil || id == "" || len(rec.Result().Cookies()) != 1 {
		t.Fatalf("id=%q cookies=%#v err=%v", id, rec.Result().Cookies(), err)
	}
	cookie := rec.Result().Cookies()[0]
	var saved model.Session
	if err := db.Where("userId = ?", id).First(&saved).Error; err != nil {
		t.Fatal(err)
	}

	if err := ss.Delete(CACHE_KEY_SESSION_ID + saved.ID); err != nil {
		t.Fatal(err)
	}
	middlewareRec := httptest.NewRecorder()
	if err := auth.Middleware()(func(c *echo.Context) error { return nil })(e.NewContext(requestWithCookie(cookie), middlewareRec)); err != nil {
		t.Fatal(err)
	}
	if _, err := auth.ActionLogout(contextWithUser(e, auth, id, saved.ID)); err != nil {
		t.Fatal(err)
	}
}
