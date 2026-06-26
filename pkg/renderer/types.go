package renderer

import (
	"html/template"
	"sync"

	"github.com/labstack/echo/v5"
)

type TemplateFuncs func(c *echo.Context, data map[string]any) template.FuncMap
type Option func(config *RuntimeConfig)

type RuntimeConfig struct {
	Funcs []TemplateFuncs
}

type TemplateNode struct {
	FilePath string
	Layouts  map[string]TemplateNode
}
type TemplateConfig struct {
	BasePath        string
	Layouts         map[string]TemplateNode
	SharedTmplPaths []string
	Runtime         RuntimeConfig
	SharedData      func(c *echo.Context, layoutNames []string) map[string]any
}
type TemplateRenderer struct {
	config *TemplateConfig
	shared *template.Template
	// templates *template.Template
	cache map[string]*template.Template
	mu    sync.RWMutex
}
