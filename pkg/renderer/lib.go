package renderer

import (
	"fmt"
	"strings"
)

func parseTemplateName(name string) ([]string, string, error) {
	parts := strings.Split(name, ":")
	if len(parts) < 2 {
		return nil, "", fmt.Errorf("template name must be layout:tpl, got %q", name)
	}

	tpl := parts[len(parts)-1]
	if tpl == "" {
		return nil, "", fmt.Errorf("template path is empty, got %q", name)
	}

	layoutNames := parts[:len(parts)-1]

	return layoutNames, tpl, nil
}
