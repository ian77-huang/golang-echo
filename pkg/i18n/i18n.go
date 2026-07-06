package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/labstack/echo/v5"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const localizerKey = "localizer"

func New(config Config) (*I18n, error) {
	config, tags, err := normalizeConfig(config)
	if err != nil {
		return nil, err
	}

	defaultTag, err := language.Parse(config.DefaultLang)
	if err != nil {
		return nil, err
	}

	bundle := goi18n.NewBundle(defaultTag)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, messageFile := range config.MessageFiles {
		if _, err := bundle.LoadMessageFileFS(config.MessageFS, messageFile); err != nil {
			return nil, err
		}
	}

	return &I18n{
		bundle:                 bundle,
		defaultLang:            config.DefaultLang,
		supportedLanguageCodes: config.SupportedLanguageCodes,
		langMatcher:            language.NewMatcher(tags),
	}, nil
}

func T(c *echo.Context, messageID string, pairs ...any) string {

	localizer, ok := c.Get(localizerKey).(*goi18n.Localizer)
	if !ok {
		return "localizer does not exist."
	}

	return t(localizer, messageID, pairs...)
}

func KV(key string, value any) templateDataPair {
	return templateDataPair{
		key:   key,
		value: value,
	}
}

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

			langCode := i.supportedLang(langCandidates...)

			c.Set(localizerKey, goi18n.NewLocalizer(i.bundle, langCode))
			c.Set("lang", langCode)

			return next(c)
		}
	}
}
