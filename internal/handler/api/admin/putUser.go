package admin

import (
	"errors"
	"log"
	"time"

	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/internal/shared"
	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
)

func (h *ApiAminHandler) PutUser(c *echo.Context) error {
	var req RequestPutUser

	if err := c.Bind(&req); err != nil {
		log.Printf("\n====== 1 %+v ======\n", err)
		return response.ErrorBadRequest(c, "invalid request")
	}

	if err := c.Validate(req); err != nil {
		log.Printf("\n====== 2 %+v ======\n", err)
		return response.ValidationError(c, err)
	}

	user := auth.GetUser[model.User](c)
	if user.ID == req.Id {
		if req.IsActive != nil && !*req.IsActive {
			log.Printf("\n====== 3 ======\n")
			return errors.New(shared.T(c, "user.cannot_disable_yourself"))
		}
	}

	userService := service.NewUserService(h.DB)
	sessionService := service.NewSessionService(h.DB)

	updatedAt := time.Now()
	user, err := userService.UpdateUserMap(req.Id, map[string]any{
		"is_active":  *req.IsActive,
		"is_admin":   *req.IsAdmin,
		"updated_at": updatedAt,
	})
	if err != nil {
		return response.ValidationErrorAuth(c, err)
	}

	if !user.Data.IsActive {
		_, err := sessionService.DeleteSession(user.ID)
		if err != nil {
			return errors.New(shared.T(c, "user.user_modify_failed"))
		}
	}

	return response.JsonOk(c, map[string]any{
		"message": shared.T(c, "user.profile.update.success"),
		"data":    user,
	})
}
