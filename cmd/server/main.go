package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appConfig "github.com/ian77-huang/golang-echo/internal/config"

	"github.com/ian77-huang/golang-echo/internal/router"

	"github.com/ian77-huang/golang-echo/pkg/database"
	"github.com/ian77-huang/golang-echo/pkg/renderer"
	"github.com/ian77-huang/golang-echo/pkg/store"
	"github.com/ian77-huang/golang-echo/pkg/validator"
	"github.com/ian77-huang/golang-echo/pkg/ws"
	"github.com/ian77-huang/golang-echo/service"

	"github.com/labstack/echo/v5"
	echomiddleware "github.com/labstack/echo/v5/middleware"
)

func main() {
	wsHub := ws.NewHub()

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

	e.Use(ws.Middleware(wsHub))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:    port,
		Handler: e,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Error("failed to start server", "error", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		e.Logger.Error("failed to stop server", "error", err)
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

	storeServer := appConfig.Store(config.RedisURL)

	auth := appConfig.Auth(&appConfig.AuthParameter{
		UserService:    service.NewUserService(db),
		SessionService: service.NewSessionService(db),
	}, storeServer)

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
