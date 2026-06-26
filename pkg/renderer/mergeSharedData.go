package renderer

import (
	"maps"

	"github.com/labstack/echo/v5"
)

func (t *TemplateRenderer) mergeSharedData(c *echo.Context, name string, data any) (any, error) {
	layoutNames, _, err := parseTemplateName(name)
	if err != nil {
		return nil, err
	}

	shared := map[string]any{}

	if t.config.SharedData != nil {
		maps.Copy(shared, t.config.SharedData(c, layoutNames))
	}

	if pageData, ok := data.(map[string]any); ok {
		maps.Copy(shared, pageData)
	}

	return shared, nil
}
