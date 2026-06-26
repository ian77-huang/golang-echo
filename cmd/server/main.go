package main

import (
	"github.com/ian77-huang/golang-echo/internal/config"
	"github.com/ian77-huang/golang-echo/internal/router"
	"github.com/ian77-huang/golang-echo/pkg/renderer"
)

func main() {
	e := router.New()
	t := renderer.New(config.RendererTemplate())

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
