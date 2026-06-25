package renderer

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestResolveTemplateFilesWithNestedLayouts(t *testing.T) {
	t.Parallel()

	renderer := &TemplateRenderer{
		config: &TemplateConfig{
			BasePath: "views",
			Layouts: map[string]TemplateNode{
				"frontend": {
					FilePath: "layout.html",
					Layouts: map[string]TemplateNode{
						"user": {
							FilePath: "layout.html",
						},
					},
				},
			},
		},
	}

	got, err := renderer.resolveTemplateFiles("frontend:user:index/index.html")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		filepath.Join("views", "frontend", "layout.html"),
		filepath.Join("views", "frontend", "user", "layout.html"),
		filepath.Join("views", "frontend", "user", "index", "index.html"),
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("resolved files mismatch:\nwant: %#v\n got: %#v", want, got)
	}
}

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
