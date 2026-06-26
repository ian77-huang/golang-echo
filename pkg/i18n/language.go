package i18n

import "golang.org/x/text/language"

func languageTags(codes []string) ([]language.Tag, error) {
	tags := make([]language.Tag, 0, len(codes))
	for _, code := range codes {
		tag, err := language.Parse(code)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (i *I18n) supportedLang(candidates ...string) string {
	for _, candidate := range candidates {
		tags, _, err := language.ParseAcceptLanguage(candidate)
		if err != nil || len(tags) == 0 {
			continue
		}

		_, index, confidence := i.langMatcher.Match(tags...)
		if confidence == language.No {
			continue
		}

		return i.supportedLanguageCodes[index]
	}

	return i.defaultLang
}
