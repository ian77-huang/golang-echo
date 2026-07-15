package service

import (
	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/pkg/cast"
)

func (s *UserService) GetUserProfile(id string) (*model.UserProfile, error) {
	userId, err := cast.StringToInt(id, 0)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.GetUserProfile(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
