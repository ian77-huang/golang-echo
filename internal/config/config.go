package config

import (
	"github.com/ian77-huang/golang-echo/pkg/renderer"
	"github.com/labstack/echo/v5"
)

func RendererTemplate() *renderer.TemplateConfig {
	return &renderer.TemplateConfig{
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
		SharedData: func(c *echo.Context, layoutNames []string) map[string]any {

			return map[string]any{}
		},
	}
}
