package users

import (
	"log"
	"net/http"

	// "time"

	"github.com/ian77-huang/golang-echo/internal/config"
	// "github.com/ian77-huang/golang-echo/internal/models/users"
	"github.com/ian77-huang/golang-echo/internal/response"
	// "github.com/ian77-huang/golang-echo/pkg/argon2"
	// "github.com/ian77-huang/golang-echo/pkg/database"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"

	"github.com/labstack/echo/v5"
)

type Request struct {
	Account         string `json:"account" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}
type Users struct {
}
type Session struct {
}

func PostRegister(c *echo.Context) error {
	var req Request

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request")
	}

	config := config.Load()
	if len(req.Account) < config.Users.MinLengthAccount {
		err := response.NewFieldError("account", "min", config.Users.MinLengthAccount)
		return response.ValidationCustomError(c, err)
	}
	if len(req.Password) < config.Users.MinLengthPassword {
		err := response.NewFieldError("password", "min", config.Users.MinLengthPassword)
		return response.ValidationCustomError(c, err)
	}
	if err := c.Validate(req); err != nil {
		return response.ValidationError(c, err)
	}
	// bbb, okbbb := c.Get("rbbrAuth").(*appAuth.Auth[Users, Session])
	cc := c.Get("rbbrAuth")
	if cc == nil {
		return echo.NewHTTPError(500, "====== 自定義 Context 轉換失敗 ======")
	}
	data, ok := cc.(appAuth.Auth[Users, Session])
	if !ok {
		log.Printf("\n======%+v ======== %+v ++++++\n", data, ok)
		return echo.NewHTTPError(500, "====== 自定義 Context 轉換失敗 2 ======")
	}
	// auth := appAuth.GetAuth[Users, Session](c)

	// ID, err := auth.Register(c, req.Account, req.Password)
	// if err != nil {
	// 	return err
	// }

	// db, err := database.GetDBConnect(c)
	// if err != nil {
	// 	return err // 若失敗，直接回傳錯誤（Go 會自動處理 HTTP 狀態碼）
	// }

	// if ok, _ := users.IsAccountExist(db, req.Account); ok {
	// 	return response.Error(c, http.StatusConflict, "account already exists")
	// }

	// c.SetCookie(&http.Cookie{
	// 	Name:     "lang",
	// 	Value:    req.Account,
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(365 * 24 * time.Hour),
	// 	HttpOnly: true,
	// 	SameSite: http.SameSiteLaxMode,
	// })

	// now := time.Now()

	// passwordHash, err := argon2.HashPassword(req.Password)
	// if err != nil {
	// 	// 這裡通常代表系統底層的隨機數產生器（crypto/rand）發生問題
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "failed to secure password")
	// }

	// newUser := users.User{
	// 	Account:   req.Account,
	// 	Password:  passwordHash, // 💡 實務上建議使用 golang.org/x/crypto/bcrypt 進行密碼雜湊加密
	// 	CreatedAt: &now,
	// 	UpdatedAt: &now,
	// }

	// // 5. 使用 GORM 寫入資料庫 (指定寫入到 "users" 資料表)
	// if err := users.CreateUser(db, &newUser); err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	// }

	ID := "12345"
	// 6. 回傳成功訊息與建立好的使用者 ID
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "user created successfully",
		"id":      ID,
	})
}
