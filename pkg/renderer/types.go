package renderer

import (
	"html/template"
	"sync"

	"github.com/labstack/echo/v5"
)

type TemplateNode struct {
	FilePath string
	Layouts  map[string]TemplateNode
}
type TemplateConfig struct {
	BasePath        string
	Layouts         map[string]TemplateNode
	SharedTmplPaths []string
	SharedData      func(c *echo.Context, layoutNames []string) map[string]any
}
type TemplateRenderer struct {
	config *TemplateConfig
	shared *template.Template
	// templates *template.Template
	cache map[string]*template.Template
	mu    sync.RWMutex
}
