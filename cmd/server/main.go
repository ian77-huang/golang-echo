package main

import (
	"os"

	appConfig "github.com/ian77-huang/golang-echo/internal/config"

	"github.com/ian77-huang/golang-echo/internal/router"

	"github.com/ian77-huang/golang-echo/pkg/database"
	"github.com/ian77-huang/golang-echo/pkg/renderer"
	"github.com/ian77-huang/golang-echo/pkg/store"
	"github.com/ian77-huang/golang-echo/pkg/validator"
	"github.com/ian77-huang/golang-echo/service"

	"github.com/labstack/echo/v5"
	echomiddleware "github.com/labstack/echo/v5/middleware"
)

type User struct {
}
type Session struct {
}

func main() {
	name := os.Args[1]
	if name == "" {
		name = "echo"
	}
	port := os.Args[2]
	if port == "" {
		port = ":1323"
	}

	e, err := newServer(name)
	if err != nil {
		panic(err)
	}

	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "name", name, "error", err)
	}
}

func newServer(serverName string) (*echo.Echo, error) {
	config := appConfig.Load()

	translator, err := appConfig.I18n()
	if err != nil {
		return nil, err
	}

	db, err := appConfig.DB(config.Databases.Path)
	if err != nil {
		return nil, err
	}

	auth := appConfig.Auth(&appConfig.AuthParameter{
		UserService:    service.NewUserService(db),
		SessionService: service.NewSessionService(db),
	})

	storeServer := appConfig.Store(config.RedisURL)

	e := router.New(&router.RouterParameter{DB: db})
	t := renderer.New(appConfig.RendererTemplate(
		renderer.WithFuncs(translator.TemplateFuncs),
	))

	e.Use(echomiddleware.GzipWithConfig(echomiddleware.GzipConfig{
		Level: 5,
	}))

	e.Static("/assets", "assets")
	e.Validator = validator.New()
	e.Renderer = t

	e.Use(appConfig.Middleware(&appConfig.ConfigMiddleware{ServerName: serverName}))
	e.Use(echomiddleware.RequestLogger())
	e.Use(echomiddleware.Recover())
	e.Use(translator.Middleware())
	e.Use(database.Middleware(db))
	e.Use(store.Middleware(storeServer))

	e.Use(auth.Middleware())

	return e, nil
}
