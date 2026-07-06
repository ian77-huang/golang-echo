package i18n

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/labstack/echo/v5"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

func (i *I18n) TemplateFuncs(c *echo.Context, _ map[string]any) template.FuncMap {
	if c == nil {
		return fallbackTemplateFuncs()
	}

	localizer, ok := i.Localizer(c)
	if !ok {
		return fallbackTemplateFuncs()
	}

	return template.FuncMap{
		"kv": KV,
		"t": func(messageID string, pairs ...any) string {
			return t(localizer, messageID, pairs...)
		},
	}
}
func (i *I18n) Localizer(c *echo.Context) (*goi18n.Localizer, bool) {
	localizer, ok := c.Get(localizerKey).(*goi18n.Localizer)
	return localizer, ok
}

func t(localizer *goi18n.Localizer, messageID string, pairs ...any) string {
	message, err := localizer.Localize(&goi18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData(pairs),
	})
	if err != nil {
		return messageID
	}

	return message
}

func fallbackTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"kv": KV,
		"t": func(messageID string, _ ...any) string {
			return messageID
		},
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
