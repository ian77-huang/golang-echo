package service

import (
	"errors"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/pkg/math"
	"gorm.io/gorm"
)

func (s *bibleService) GetBible(id int) (*model.Bible, error) {
	bible, err := s.repo.GetBible(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return bible, nil
}

func (s *bibleService) GetBibleByDate() (*model.Bible, error) {
	now := time.Now()
	dateStr := now.Format("2006-01-02")

	rand := math.CryptoNumberFromStringSHA(dateStr, 100)

	bible, err := s.GetBible(rand)
	if err != nil {
		return nil, err
	}
	return bible, nil
}
