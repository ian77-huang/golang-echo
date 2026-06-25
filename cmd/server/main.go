package main

import (
	"github.com/ian77-huang/golang-echo/internal/router"
	"github.com/ian77-huang/golang-echo/pkg/renderer"
)

func main() {
	config := &renderer.TemplateConfig{
		BasePath: "views",
		Layouts: map[string]renderer.TemplateNode{
			"frontend": {
				FilePath: "layout.html",
			},
			"admin": {
				FilePath: "layout.html",
			},
		},
		SharedTmplPaths: []string{"base.html"},
	}

	e := router.New()
	t := renderer.New(config)

	e.Renderer = t

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

// .
// ├── cmd/
// │   └── api/
// │       └── main.go
// ├── internal/
// │   ├── config/
// │   ├── handler/
// │   ├── service/
// │   ├── repository/
// │   └── model/
// ├── pkg/
// ├── migrations/
// ├── go.mod
// └── .env
