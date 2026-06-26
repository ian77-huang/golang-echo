package i18n

import (
	"testing"

	"golang.org/x/text/language"
)

func TestSupportedLangMatchesSupportedLanguageCodes(t *testing.T) {
	t.Parallel()

	translator := &I18n{
		defaultLang:            "zh-TW",
		supportedLanguageCodes: []string{"zh-TW", "en"},
		langMatcher:            newTestMatcher([]string{"zh-TW", "en"}),
	}

	tests := []struct {
		name       string
		candidates []string
		want       string
	}{
		{
			name:       "exact traditional chinese",
			candidates: []string{"zh-TW"},
			want:       "zh-TW",
		},
		{
			name:       "regional english",
			candidates: []string{"en-US"},
			want:       "en",
		},
		{
			name:       "browser accept language",
			candidates: []string{"zh-Hant-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7"},
			want:       "zh-TW",
		},
		{
			name:       "uses next candidate when first is unsupported",
			candidates: []string{"fr-FR", "en-US"},
			want:       "en",
		},
		{
			name:       "defaults when no candidate is supported",
			candidates: []string{"fr-FR"},
			want:       "zh-TW",
		},
		{
			name:       "defaults when candidates are empty",
			candidates: []string{"", "  "},
			want:       "zh-TW",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := translator.supportedLang(tt.candidates...); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func newTestMatcher(codes []string) language.Matcher {
	tags, err := languageTags(codes)
	if err != nil {
		panic(err)
	}

	return language.NewMatcher(tags)
}
