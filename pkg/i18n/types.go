package i18n

import (
	"io/fs"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Config struct {
	DefaultLang            string
	SupportedLanguageCodes []string
	MessageFS              fs.FS
	MessageFiles           []string
}

type I18n struct {
	bundle                 *goi18n.Bundle
	defaultLang            string
	supportedLanguageCodes []string
	langMatcher            language.Matcher
}

type templateDataPair struct {
	key   string
	value any
}
