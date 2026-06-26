package i18n

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/labstack/echo/v5"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type templateDataPair struct {
	key   string
	value any
}

func (i *I18n) TemplateFuncs(c *echo.Context, _ map[string]any) template.FuncMap {
	if c == nil {
		return fallbackTemplateFuncs()
	}

	localizer, ok := i.Localizer(c)
	if !ok {
		return fallbackTemplateFuncs()
	}

	return template.FuncMap{
		"kv": templateDataValue,
		"t": func(messageID string, pairs ...any) string {
			message, err := localizer.Localize(&goi18n.LocalizeConfig{
				MessageID:    messageID,
				TemplateData: templateData(pairs),
			})
			if err != nil {
				return messageID
			}

			return message
		},
	}
}

func fallbackTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"kv": templateDataValue,
		"t": func(messageID string, _ ...any) string {
			return messageID
		},
	}
}

func templateDataValue(key string, value any) templateDataPair {
	return templateDataPair{
		key:   key,
		value: value,
	}
}

func templateData(pairs []any) map[string]any {
	if len(pairs) == 0 {
		return nil
	}

	data := make(map[string]any, len(pairs))
	for _, pair := range pairs {
		if kv, ok := pair.(templateDataPair); ok {
			data[kv.key] = kv.value
			continue
		}

		key, value, ok := strings.Cut(fmt.Sprint(pair), "=")
		if !ok {
			continue
		}

		data[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return data
}
