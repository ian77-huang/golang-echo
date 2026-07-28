package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/pkg/math"
	"gorm.io/gorm"
)

func (s *bibleService) GetBible(id int) (*model.Bible, error) {
	val, err, _ := s.g.Do("GetBible-"+strconv.Itoa(id), func() (interface{}, error) {
		bible, err := s.repo.GetBible(id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return bible, nil
	})

	if val == nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return val.(*model.Bible), nil
}

func (s *bibleService) GetBibleByDate() (*model.Bible, error) {
	now := time.Now()
	dateStr := now.Format("2006-01-02")

	val, err, _ := s.g.Do("GetBibleByDate-"+dateStr, func() (interface{}, error) {
		rand := math.CryptoNumberFromStringSHA(dateStr, 100)

		bible, err := s.GetBible(rand)
		if err != nil {
			return nil, err
		}
		return bible, nil
	})

	if err != nil {
		return nil, err
	}

	return val.(*model.Bible), nil
}
