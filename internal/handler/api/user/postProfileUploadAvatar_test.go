package user

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var pngBytes = []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 0, 0, 0, 13, 'I', 'H', 'D', 'R'}

func avatarRequest(filename string, content []byte) *http.Request {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	fw, _ := w.CreateFormFile("avatar", filename)
	fw.Write(content)
	w.Close()
	req := httptest.NewRequest(http.MethodPost, "/api/user/profile/avatar", body)
	req.Header.Set(echo.HeaderContentType, w.FormDataContentType())
	return req
}

func newProfileDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.UserProfile{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestPostProfileUploadAvatarSavesFile(t *testing.T) {
	assets := t.TempDir()
	t.Setenv("ASSETS_PATH", assets)
	t.Setenv("MAX_SIZE_USER_PROFILE_AVATAR", "1048576")
	e := echo.New()
	h := &ApiUserHandler{DB: newProfileDB(t)}
	rec := httptest.NewRecorder()
	c := e.NewContext(avatarRequest("a.png", pngBytes), rec)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "7"})
	if err := h.PostProfileUploadAvatar(c); err != nil || rec.Code != http.StatusOK {
		t.Fatalf("status=%d err=%v body=%s", rec.Code, err, rec.Body.String())
	}
	if _, err := os.Stat(filepath.Join(assets, "uploads", "user", "profile", "avatar", "7.png")); err != nil {
		t.Fatalf("avatar file not saved: %v", err)
	}
}

func TestPostProfileUploadAvatarRejectsMissingFile(t *testing.T) {
	t.Setenv("ASSETS_PATH", t.TempDir())
	t.Setenv("MAX_SIZE_USER_PROFILE_AVATAR", "1048576")
	e := echo.New()
	h := &ApiUserHandler{}
	req := httptest.NewRequest(http.MethodPost, "/api/user/profile/avatar", bytes.NewBufferString("x"))
	req.Header.Set(echo.HeaderContentType, "text/plain")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := h.PostProfileUploadAvatar(c); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}

func TestPostProfileUploadAvatarRejectsTooLarge(t *testing.T) {
	t.Setenv("ASSETS_PATH", t.TempDir())
	t.Setenv("MAX_SIZE_USER_PROFILE_AVATAR", "-1")
	e := echo.New()
	h := &ApiUserHandler{}
	rec := httptest.NewRecorder()
	c := e.NewContext(avatarRequest("a.png", pngBytes), rec)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "7"})
	if err := h.PostProfileUploadAvatar(c); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}

func TestPostProfileUploadAvatarRejectsInvalidType(t *testing.T) {
	t.Setenv("ASSETS_PATH", t.TempDir())
	t.Setenv("MAX_SIZE_USER_PROFILE_AVATAR", "1048576")
	e := echo.New()
	h := &ApiUserHandler{DB: newProfileDB(t)}
	rec := httptest.NewRecorder()
	c := e.NewContext(avatarRequest("a.txt", []byte("plain text here")), rec)
	c.Set(appAuth.CONTEXT_KEY_USER, &appAuth.User[model.User]{ID: "7"})
	if err := h.PostProfileUploadAvatar(c); err != nil || rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d err=%v", rec.Code, err)
	}
}
