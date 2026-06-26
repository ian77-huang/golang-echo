package renderer

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
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

func TestRenderUsesRequestTemplateFuncs(t *testing.T) {
	t.Parallel()

	basePath := t.TempDir()
	frontendPath := filepath.Join(basePath, "frontend", "index")
	if err := os.MkdirAll(frontendPath, 0o755); err != nil {
		t.Fatal(err)
	}

	writeTestFile(t, filepath.Join(basePath, "base.html"), `{{define "base"}}{{template "content" .}}{{end}}`)
	writeTestFile(t, filepath.Join(basePath, "frontend", "layout.html"), `{{define "layout"}}{{end}}`)
	writeTestFile(t, filepath.Join(frontendPath, "index.html"), `{{define "content"}}{{t "welcome"}}{{end}}`)

	renderer := New(&TemplateConfig{
		BasePath: basePath,
		Layouts: map[string]TemplateNode{
			"frontend": {
				FilePath: "layout.html",
			},
		},
		SharedTmplPaths: []string{"base.html"},
		Runtime: RuntimeConfig{
			Funcs: []TemplateFuncs{
				func(c *echo.Context, data map[string]any) template.FuncMap {
					if c == nil {
						return template.FuncMap{
							"t": func(key string) string {
								return key
							},
						}
					}

					return template.FuncMap{
						"t": func(key string) string {
							return "translated " + key
						},
					}
				},
			},
		},
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var out bytes.Buffer
	if err := renderer.Render(c, io.Writer(&out), "frontend:index/index.html", map[string]any{}); err != nil {
		t.Fatal(err)
	}

	if got, want := out.String(), "translated welcome"; got != want {
		t.Fatalf("rendered output mismatch: want %q, got %q", want, got)
	}
}

func TestRenderDoesNotRequireTemplateFuncs(t *testing.T) {
	t.Parallel()

	basePath := t.TempDir()
	frontendPath := filepath.Join(basePath, "frontend", "index")
	if err := os.MkdirAll(frontendPath, 0o755); err != nil {
		t.Fatal(err)
	}

	writeTestFile(t, filepath.Join(basePath, "base.html"), `{{define "base"}}{{template "content" .}}{{end}}`)
	writeTestFile(t, filepath.Join(basePath, "frontend", "layout.html"), `{{define "layout"}}{{end}}`)
	writeTestFile(t, filepath.Join(frontendPath, "index.html"), `{{define "content"}}hello{{end}}`)

	renderer := New(&TemplateConfig{
		BasePath: basePath,
		Layouts: map[string]TemplateNode{
			"frontend": {
				FilePath: "layout.html",
			},
		},
		SharedTmplPaths: []string{"base.html"},
	})

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var out bytes.Buffer
	if err := renderer.Render(c, io.Writer(&out), "frontend:index/index.html", map[string]any{}); err != nil {
		t.Fatal(err)
	}

	if got, want := out.String(), "hello"; got != want {
		t.Fatalf("rendered output mismatch: want %q, got %q", want, got)
	}
}

func writeTestFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
