package i18n

import (
	"io/fs"

	"github.com/BurntSushi/toml"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const defaultLocalizerKey = "localizer"

type Config struct {
	DefaultLang            string
	SupportedLanguageCodes []string
	MessageFS              fs.FS
	MessageFiles           []string
	LocalizerKey           string
}

type I18n struct {
	bundle                 *goi18n.Bundle
	defaultLang            string
	supportedLanguageCodes []string
	langMatcher            language.Matcher
	localizerKey           string
}

func New(config Config) (*I18n, error) {
	config, tags, err := normalizeConfig(config)
	if err != nil {
		return nil, err
	}

	defaultTag, err := language.Parse(config.DefaultLang)
	if err != nil {
		return nil, err
	}

	bundle := goi18n.NewBundle(defaultTag)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, messageFile := range config.MessageFiles {
		if _, err := bundle.LoadMessageFileFS(config.MessageFS, messageFile); err != nil {
			return nil, err
		}
	}

	return &I18n{
		bundle:                 bundle,
		defaultLang:            config.DefaultLang,
		supportedLanguageCodes: config.SupportedLanguageCodes,
		langMatcher:            language.NewMatcher(tags),
		localizerKey:           config.LocalizerKey,
	}, nil
}
