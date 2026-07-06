package i18n

import "testing"

func TestTemplateDataParsesKeyValuePairs(t *testing.T) {
	t.Parallel()

	got := templateData([]any{"Name=Yien", "City=Taipei"})

	if got["Name"] != "Yien" {
		t.Fatalf("expected Name to be Yien, got %#v", got["Name"])
	}

	if got["City"] != "Taipei" {
		t.Fatalf("expected City to be Taipei, got %#v", got["City"])
	}
}

func TestTemplateDataParsesTemplateDataPairs(t *testing.T) {
	t.Parallel()

	got := templateData([]any{
		KV("Name", "Yien"),
		KV("Count", 3),
	})

	if got["Name"] != "Yien" {
		t.Fatalf("expected Name to be Yien, got %#v", got["Name"])
	}

	if got["Count"] != 3 {
		t.Fatalf("expected Count to be 3, got %#v", got["Count"])
	}
}

func TestTemplateDataIgnoresInvalidPairs(t *testing.T) {
	t.Parallel()

	got := templateData([]any{"Name=Yien", "invalid"})

	if len(got) != 1 {
		t.Fatalf("expected one valid template data item, got %#v", got)
	}

	if got["Name"] != "Yien" {
		t.Fatalf("expected Name to be Yien, got %#v", got["Name"])
	}
}
