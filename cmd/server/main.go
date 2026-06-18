package main

import (
	"github.com/ian77-huang/golang-echo/internal/router"
)

func main() {
	e := router.New()

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
