package user

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	appConfig "github.com/ian77-huang/golang-echo/internal/config"
	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/internal/shared"
	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/service"

	"github.com/labstack/echo/v5"
)

func (h *ApiUserHandler) PostProfileUploadAvatar(c *echo.Context) error {

	file, err := c.FormFile("avatar")
	if err != nil {
		return response.ErrorBadRequest(c, shared.T(c, "file.upload.read.failed"))
	}

	config := appConfig.Load()
	maxSize := int64(config.MaxSizeUserProfileAvatar)
	if file.Size > maxSize {
		//
		return response.ErrorBadRequest(c, shared.T(c, "file.upload.too_large"))
	}

	src, err := file.Open()
	if err != nil {
		return response.ErrorInternalServerError(c, "file.upload.open.failed")
	}
	defer src.Close()

	buf := make([]byte, 512)
	_, err = src.Read(buf)
	if err != nil && err != io.EOF {
		return response.ErrorInternalServerError(c, "file.upload.read.failed")
	}

	// 判斷 MIME Type
	contentType := http.DetectContentType(buf)

	allowed := map[string]struct{}{
		"image/jpeg": {},
		"image/png":  {},
		"image/webp": {},
	}

	if _, ok := allowed[contentType]; !ok {
		return response.ErrorBadRequest(c, "file.upload.invalid_type")
	}

	if seeker, ok := src.(io.Seeker); ok {
		_, err = seeker.Seek(0, io.SeekStart)
		if err != nil {
			return response.ErrorInternalServerError(c, "seek failed")
		}
	}

	user := appAuth.GetUser[model.User](c)
	userService := service.NewUserService(h.DB)

	// 2. 確保儲存目錄存在
	uploadDir := "/uploads/user/profile/avatar"
	if err := os.MkdirAll(config.AssetsPath+uploadDir, os.ModePerm); err != nil {
		return response.ErrorInternalServerError(c, "file.upload.create.directory.failed")
	}

	// 3. 重新命名檔案防止重複 (以 User ID 1 為例，副檔名保留)
	ext := filepath.Ext(file.Filename)
	filename := user.ID + ext
	dstPath := filepath.Join(config.AssetsPath+uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return response.ErrorInternalServerError(c, "file.upload.create.file.failed")
	}
	defer dst.Close()

	// 4. 複製檔案內容
	if _, err = io.Copy(dst, src); err != nil {
		return response.ErrorInternalServerError(c, "file.upload.write.failed")
	}

	avatarURL := strings.TrimPrefix(config.AssetsPath, ".") + uploadDir + "/" + filename + "?t=" + strconv.FormatInt(time.Now().Unix(), 10)
	profile, err := userService.UpdateUserProfile(user.ID, &model.UserProfile{
		AvatarURL: avatarURL,
	})
	if err != nil {
		return err
	}
	log.Printf("\n======= avatarURL = %+v ========\n", avatarURL)
	return response.JsonOk(c, map[string]string{
		"message":    shared.T(c, "file.upload.success"),
		"avatar_url": profile.AvatarURL,
	})
}
