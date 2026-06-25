package renderer

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

func (t *TemplateRenderer) buildTemplate(name string) (*template.Template, error) {
	filePaths, err := t.resolveTemplateFiles(name)
	if err != nil {
		return nil, err
	}

	tmpl, err := t.shared.Clone()
	if err != nil {
		return nil, err
	}

	_, err = tmpl.ParseFiles(filePaths...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (t *TemplateRenderer) resolveTemplateFiles(name string) ([]string, error) {
	parts := strings.Split(name, ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("template name must be layout:tpl, got %q", name)
	}

	tpl := parts[len(parts)-1]
	if tpl == "" {
		return nil, fmt.Errorf("template path is empty, got %q", name)
	}

	layoutNames := parts[:len(parts)-1]

	filePaths := make([]string, 0)

	tempPath := t.config.BasePath
	tempNodes := t.config.Layouts

	for _, layoutName := range layoutNames {
		node, ok := tempNodes[layoutName]
		if !ok {
			return nil, fmt.Errorf("layout %q not found", layoutName)
		}

		if node.FilePath == "" {
			return nil, fmt.Errorf("layout %q FilePath is empty", layoutName)
		}

		tempPath = filepath.Join(tempPath, layoutName)
		filePaths = append(filePaths, filepath.Join(tempPath, node.FilePath))

		tempNodes = node.Layouts
	}

	filePaths = append(filePaths, filepath.Join(tempPath, tpl))

	return filePaths, nil
}
