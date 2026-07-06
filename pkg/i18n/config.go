package i18n

import (
	"errors"

	"golang.org/x/text/language"
)

func normalizeConfig(config Config) (Config, []language.Tag, error) {
	if config.DefaultLang == "" {
		return config, nil, errors.New("default language is required")
	}

	if len(config.SupportedLanguageCodes) == 0 {
		return config, nil, errors.New("at least one supported language code is required")
	}

	if len(config.MessageFiles) > 0 && config.MessageFS == nil {
		return config, nil, errors.New("message fs is required when message files are configured")
	}

	tags, err := languageTags(config.SupportedLanguageCodes)
	if err != nil {
		return config, nil, err
	}

	return config, tags, nil
}
