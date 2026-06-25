package renderer

import (
	"html/template"
	"sync"
)

type TemplateNode struct {
	FilePath string
	Layouts  map[string]TemplateNode
}
type TemplateConfig struct {
	BasePath        string
	Layouts         map[string]TemplateNode
	SharedTmplPaths []string
}
type TemplateRenderer struct {
	config *TemplateConfig
	shared *template.Template
	// templates *template.Template
	cache map[string]*template.Template
	mu    sync.RWMutex
}
