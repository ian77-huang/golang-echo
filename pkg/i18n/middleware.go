package i18n

import (
	"github.com/labstack/echo/v5"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

func (i *I18n) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			langCandidates := []string{}

			if cookie, err := c.Cookie("lang"); err == nil && cookie.Value != "" {
				langCandidates = append(langCandidates, cookie.Value)
			}
			langCandidates = append(
				langCandidates,
				c.Request().Header.Get("X-Language"),
				c.Request().Header.Get("Accept-Language"),
			)

			c.Set(i.localizerKey, goi18n.NewLocalizer(i.bundle, i.supportedLang(langCandidates...)))

			return next(c)
		}
	}
}

func (i *I18n) Localizer(c *echo.Context) (*goi18n.Localizer, bool) {
	localizer, ok := c.Get(i.localizerKey).(*goi18n.Localizer)
	return localizer, ok
}
