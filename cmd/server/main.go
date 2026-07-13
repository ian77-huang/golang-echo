package main

import (
	appConfig "github.com/ian77-huang/golang-echo/internal/config"

	"github.com/ian77-huang/golang-echo/internal/router"

	"github.com/ian77-huang/golang-echo/pkg/database"
	"github.com/ian77-huang/golang-echo/pkg/renderer"
	"github.com/ian77-huang/golang-echo/pkg/validator"

	"github.com/labstack/echo/v5"
	echomiddleware "github.com/labstack/echo/v5/middleware"
)

type User struct {
}
type Session struct {
}

func main() {
	e, err := newServer()
	if err != nil {
		panic(err)
	}

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

func newServer() (*echo.Echo, error) {
	config := appConfig.Load()

	translator, err := appConfig.I18n()
	if err != nil {
		return nil, err
	}

	db := database.NewSqlite(config.Databases.Path)

	auth := appConfig.Auth(db)

	e := router.New()
	t := renderer.New(appConfig.RendererTemplate(
		renderer.WithFuncs(translator.TemplateFuncs),
	))

	e.Validator = validator.New()
	e.Renderer = t

	e.Use(echomiddleware.RequestLogger())
	e.Use(echomiddleware.Recover())
	e.Use(translator.Middleware())
	e.Use(database.Middleware(db))

	e.Use(auth.Middleware())

	return e, nil
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
