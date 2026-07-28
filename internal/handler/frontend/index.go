package frontend

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/pkg/store"
	"github.com/ian77-huang/golang-echo/service"
	"github.com/labstack/echo/v5"
)

const CACHE_KEY_BIBLE = "bible:index"

func (f *FrontendHandler) GetIndex(c *echo.Context) error {
	storeServer := store.LoadStore(c)

	var bible *model.Bible
	err := storeServer.GetByteKeys(&[]store.Keys{
		{Key: CACHE_KEY_BIBLE, Target: &bible},
	})
	if err != nil {
		return err
	}

	g, _ := errgroup.WithContext(context.Background())

	if bible == nil {
		g.Go(func() error {

			bibleService := service.NewBibleService(f.DB)
			bible, err = bibleService.GetBibleByDate()
			if err != nil {
				return err
			}

			now := time.Now()

			endOfDay := time.Date(
				now.Year(), now.Month(), now.Day(),
				23, 59, 59, 999999999,
				now.Location(),
			)

			duration := endOfDay.Sub(now)

			storeServer.Set(CACHE_KEY_BIBLE, bible, duration)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return response.Render(c, "frontend:index:/index.html", map[string]any{
		"bible": bible,
	})
}
