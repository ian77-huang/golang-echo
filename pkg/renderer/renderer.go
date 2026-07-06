package renderer

import (
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
		return t.executeTemplate(c, w, tmpl, data)
	}

	tmpl, err = t.buildTemplate(name)
	if err != nil {
		return err
	}

	t.mu.Lock()
	t.cache[name] = tmpl
	t.mu.Unlock()

	return t.executeTemplate(c, w, tmpl, data)
}

func New(config *TemplateConfig) *TemplateRenderer {
	shared := template.New("shared")
	if funcs := templateFuncs(config.Runtime.Funcs, nil, nil); len(funcs) > 0 {
		shared.Funcs(funcs)
	}

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

func WithFuncs(funcs ...TemplateFuncs) Option {
	return func(config *RuntimeConfig) {
		config.Funcs = append(config.Funcs, funcs...)
	}
}

func (t *TemplateRenderer) executeTemplate(c *echo.Context, w io.Writer, tmpl *template.Template, data any) error {
	tmpl, err := tmpl.Clone()
	if err != nil {
		return err
	}

	pageData, _ := data.(map[string]any)
	if funcs := templateFuncs(t.config.Runtime.Funcs, c, pageData); len(funcs) > 0 {
		tmpl.Funcs(funcs)
	}

	return tmpl.ExecuteTemplate(w, "base", data)
}

func templateFuncs(providers []TemplateFuncs, c *echo.Context, data map[string]any) template.FuncMap {
	if len(providers) == 0 {
		return nil
	}

	funcs := template.FuncMap{}
	for _, provider := range providers {
		if provider == nil {
			continue
		}

		for name, fn := range provider(c, data) {
			funcs[name] = fn
		}
	}

	if len(funcs) == 0 {
		return nil
	}

	return funcs
}
