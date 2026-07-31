package frontend

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/ian77-huang/golang-echo/model"
	appMath "github.com/ian77-huang/golang-echo/pkg/math"
	"github.com/ian77-huang/golang-echo/pkg/store"
	storeRedis "github.com/ian77-huang/golang-echo/pkg/store/redis"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type renderer struct {
	name string
	data map[string]any
}

func (r *renderer) Render(_ *echo.Context, _ io.Writer, name string, data any) error {
	r.name = name
	r.data, _ = data.(map[string]any)
	return nil
}

func newTestStore(t *testing.T) *store.StoreServer {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(mr.Close)
	return store.New(&store.StoreOption{
		Redis: &store.StoreOptionRedis{
			IsUse:  true,
			Option: &storeRedis.RedisOption{Addr: mr.Addr()},
		},
	})
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.Bible{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestGetIndexRendersFrontendTemplate(t *testing.T) {
	db := newTestDB(t)
	id := appMath.CryptoNumberFromStringSHA(time.Now().Format("2006-01-02"), 100)
	if err := db.Create(&model.Bible{ID: id, Locale: "zh-TW", Verses: "test"}).Error; err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	r := &renderer{}
	e.Renderer = r
	h := &FrontendHandler{DB: db}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)
	if err := store.Middleware(newTestStore(t))(h.GetIndex)(c); err != nil {
		t.Fatal(err)
	}
	if r.name != "frontend:index:/index.html" || r.data["bible"] == nil {
		t.Fatalf("name=%q data=%#v", r.name, r.data)
	}
}

func TestGetIndexServesBibleFromCache(t *testing.T) {
	db := newTestDB(t)
	id := appMath.CryptoNumberFromStringSHA(time.Now().Format("2006-01-02"), 100)
	bible := &model.Bible{ID: id, Locale: "zh-TW", Verses: "cached"}
	if err := db.Create(bible).Error; err != nil {
		t.Fatal(err)
	}
	ss := newTestStore(t)
	if err := ss.Set(CACHE_KEY_BIBLE, bible, time.Hour); err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	r := &renderer{}
	e.Renderer = r
	h := &FrontendHandler{DB: db}
	rec := httptest.NewRecorder()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec)
	if err := store.Middleware(ss)(h.GetIndex)(c); err != nil {
		t.Fatal(err)
	}
	if r.name != "frontend:index:/index.html" || r.data["bible"] == nil {
		t.Fatalf("name=%q data=%#v", r.name, r.data)
	}
}

func TestNewRegistersFrontendRoutes(t *testing.T) {
	e := echo.New()
	New(&FrontendParameter{DB: newTestDB(t), Echo: e})
	routes := e.Router().Routes()
	got := make(map[string]bool)
	for _, route := range routes {
		if route.Method == http.MethodGet {
			got[route.Path] = true
		}
	}
	for _, expected := range []string{"/", "/ws", "/user", "/user/login", "/user/register", "/user/logout", "/user/profile", "/user/reset-password"} {
		if !got[expected] {
			t.Fatalf("missing route %s; got %v", expected, got)
		}
	}
}
