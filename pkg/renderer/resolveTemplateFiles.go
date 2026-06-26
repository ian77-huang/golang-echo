package renderer

import (
	"fmt"
	"path/filepath"
)

func (t *TemplateRenderer) resolveTemplateFiles(name string) ([]string, error) {
	layoutNames, tpl, err := parseTemplateName(name)
	if err != nil {
		return nil, err
	}

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
