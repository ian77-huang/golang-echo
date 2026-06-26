package renderer

import (
	// "fmt"
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v5"
)

func (t *TemplateRenderer) Render(c *echo.Context, w io.Writer, name string, data any) error {
	data, err := t.mergeSharedData(c, name, data)
	if err != nil {
		return err
	}

	t.mu.RLock()
	tmpl, ok := t.cache[name]
	t.mu.RUnlock()

	if ok {
		return tmpl.ExecuteTemplate(w, "base", data)
	}

	tmpl, err = t.buildTemplate(name)
	if err != nil {
		return err
	}

	t.mu.Lock()
	t.cache[name] = tmpl
	t.mu.Unlock()

	return tmpl.ExecuteTemplate(w, "base", data)
}

func New(config *TemplateConfig) *TemplateRenderer {
	shared := template.New("shared")

	filePaths := make([]string, 0)

	for _, path := range config.SharedTmplPaths {
		filePaths = append(filePaths, filepath.Join(config.BasePath, path))
	}

	if len(filePaths) > 0 {
		template.Must(shared.ParseFiles(filePaths...))
	}

	return &TemplateRenderer{
		config: config,
		shared: shared,
		cache:  make(map[string]*template.Template),
	}
}
