package config

import (
	"github.com/ian77-huang/golang-echo/internal/locales"
	appi18n "github.com/ian77-huang/golang-echo/pkg/i18n"
)

func I18n() (*appi18n.I18n, error) {
	return appi18n.New(appi18n.Config{
		DefaultLang:            "zh-TW",
		SupportedLanguageCodes: []string{"zh-TW", "en"},
		MessageFS:              locales.FS,
		MessageFiles: []string{
			"active.zh-TW.toml",
			"active.en.toml",
			"errors.en.toml",
			"errors.zh-TW.toml",
			"file.en.toml",
			"file.zh-TW.toml",
			"index.en.toml",
			"index.zh-TW.toml",
			"placeholders.en.toml",
			"placeholders.zh-TW.toml",
			"users.en.toml",
			"users.zh-TW.toml",
			"validations.en.toml",
			"validations.zh-TW.toml",
		},
	})
}
