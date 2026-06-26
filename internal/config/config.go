package config

import (
	"github.com/ian77-huang/golang-echo/internal/locales"
	appi18n "github.com/ian77-huang/golang-echo/pkg/i18n"
	"github.com/ian77-huang/golang-echo/pkg/renderer"
	"github.com/labstack/echo/v5"
)

func I18n() (*appi18n.I18n, error) {
	return appi18n.New(appi18n.Config{
		DefaultLang:            "zh-TW",
		SupportedLanguageCodes: []string{"zh-TW", "en"},
		MessageFS:              locales.FS,
		MessageFiles: []string{
			"active.zh-TW.toml",
			"active.en.toml",
		},
	})
}

func RendererTemplate(options ...renderer.Option) *renderer.TemplateConfig {
	runtime := renderer.RuntimeConfig{}
	for _, option := range options {
		option(&runtime)
	}

	config := &renderer.TemplateConfig{
		BasePath: "views",
		Layouts: map[string]renderer.TemplateNode{
			"frontend": {
				FilePath: "layout.html",
			},
			"admin": {
				FilePath: "layout.html",
			},
		},
		SharedTmplPaths: []string{"base.html"},
		Runtime:         runtime,
		SharedData: func(c *echo.Context, layoutNames []string) map[string]any {

			return map[string]any{}
		},
	}

	return config
}
