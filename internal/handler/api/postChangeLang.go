package api

import (
	"net/http"
	"time"

	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/labstack/echo/v5"
)

func (h *ApiHandler) PostChangeLang(c *echo.Context) error {
	var req ChangeLangRequest

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "invalid request")
	}
	if err := c.Validate(req); err != nil {
		return response.ValidationError(c, err)
	}

	c.SetCookie(&http.Cookie{
		Name:     "lang",
		Value:    req.Code,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return response.JSON(c, map[string]any{
		"lang": req.Code,
	})
}
