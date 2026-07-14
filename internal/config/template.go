package config

import (
	"github.com/ian77-huang/golang-echo/internal/shared"
	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/pkg/renderer"
	"github.com/labstack/echo/v5"
)

func RendererTemplate(options ...renderer.Option) *renderer.TemplateConfig {
	runtime := renderer.RuntimeConfig{}
	for _, option := range options {
		option(&runtime)
	}

	config := &renderer.TemplateConfig{
		BasePath: "internal/views",
		Layouts: map[string]renderer.TemplateNode{
			"frontend": {
				FilePath: "layout.html",
				Layouts: map[string]renderer.TemplateNode{
					"users": {
						FilePath: "layout.html",
					},
				},
			},
			"admin": {
				FilePath: "layout.html",
			},
		},
		SharedTmplPaths: []string{"base.html"},
		Runtime:         runtime,
		SharedData: func(c *echo.Context, layoutNames []string) map[string]any {
			realPath := c.Request().URL.Path
			lang, _ := c.Get("lang").(string)

			isSignedIn := appAuth.IsSignedIn[model.User](c)

			users := SetMenusUsers(MenuUsersRules{Path: realPath, IsSignedIn: isSignedIn, T: shared.TFactory(c)})

			menus := SetMenus(MenuRules{users: users, T: shared.TFactory(c)})

			return map[string]any{
				"Lang":       lang,
				"Users":      users,
				"Menus":      menus,
				"IsSignedIn": isSignedIn,
			}
		},
	}

	return config
}
