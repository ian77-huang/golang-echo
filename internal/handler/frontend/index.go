package frontend

import (
	"time"

	"github.com/ian77-huang/golang-echo/internal/response"
	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/pkg/store"
	storeRedis "github.com/ian77-huang/golang-echo/pkg/store/redis"
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

	needUpdateData := []store.RedisMSET{}

	if bible == nil {
		bibleService := service.NewBibleService(f.DB)
		bible, err = bibleService.GetBibleByDate()
		if err != nil {
			return response.ErrorInternalServerError(c, "invalid_request")
		}

		now := time.Now()

		endOfDay := time.Date(
			now.Year(), now.Month(), now.Day(),
			23, 59, 59, 999999999,
			now.Location(),
		)

		duration := endOfDay.Sub(now)

		needUpdateData = append(needUpdateData, storeRedis.RedisMSET{
			Key: CACHE_KEY_BIBLE, Value: bible, Expiration: duration,
		})
	}

	if len(needUpdateData) != 0 {
		if err := storeServer.MSet(needUpdateData); err != nil {
			return err
		}
	}

	return response.Render(c, "frontend:index:/index.html", map[string]any{
		"bible": bible,
	})
}
