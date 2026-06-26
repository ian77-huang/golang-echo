package main

import (
	"github.com/ian77-huang/golang-echo/internal/config"
	"github.com/ian77-huang/golang-echo/internal/router"
	"github.com/ian77-huang/golang-echo/pkg/renderer"
	echomiddleware "github.com/labstack/echo/v5/middleware"
)

func main() {
	translator, err := config.I18n()
	if err != nil {
		panic(err)
	}

	e := router.New()
	t := renderer.New(config.RendererTemplate(
		renderer.WithFuncs(translator.TemplateFuncs),
	))

	e.Use(echomiddleware.RequestLogger())
	e.Use(echomiddleware.Recover())
	e.Use(translator.Middleware())

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
