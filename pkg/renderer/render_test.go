package renderer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildTemplateReturnsErrorWhenTemplatePathIsEmpty(t *testing.T) {
	t.Parallel()

	basePath := t.TempDir()
	frontendPath := filepath.Join(basePath, "frontend")
	if err := os.MkdirAll(frontendPath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(basePath, "base.html"), []byte(`{{define "base"}}{{end}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(frontendPath, "layout.html"), []byte(`{{define "content-layout"}}{{end}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	renderer := New(&TemplateConfig{
		BasePath: basePath,
		Layouts: map[string]TemplateNode{
			"frontend": {
				FilePath: "layout.html",
			},
		},
		SharedTmplPaths: []string{"base.html"},
	})

	_, err := renderer.buildTemplate("frontend:")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "template path is empty") {
		t.Fatalf("expected template path error, got %q", err.Error())
	}
}
