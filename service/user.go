package service

import (
	"errors"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	"github.com/ian77-huang/golang-echo/pkg/cast"
	"gorm.io/gorm"
)

// func (s *UserService) GetUser(id string) (*model.User, error) {
// 	userId, err := cast.StringToInt(id, 0)
// 	if err != nil {
// 		return nil, err
// 	}
// 	user, err := s.repo.GetUser(userId)
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return user, nil
// }
// func (s *UserService) GetUserByAccount(account string) (*model.User, error) {
// 	user, err := s.repo.GetUserByAccount(account)
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

//		return user, nil
//	}
func (s *UserService) GetUserProfile(id string) (*model.UserProfile, error) {
	userId, err := cast.StringToInt(id, 0)
	if err != nil {
		return nil, err
	}
	profile, err := s.repo.GetUserProfile(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *UserService) UpdateUserProfile(id string, updateData *model.UserProfile) (*model.UserProfile, error) {
	userId, err := cast.StringToInt(id, 0)
	if err != nil {
		return nil, err
	}

	var profile *model.UserProfile

	profile, err = s.repo.GetUserProfile(userId)
	if err != nil {
		profile, err = s.repo.CreateUserProfile(&model.UserProfile{
			UserID:    userId,
			Name:      updateData.Name,
			Email:     updateData.Email,
			Phone:     updateData.Phone,
			Bio:       updateData.Bio,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return nil, err
		}
	} else {
		profile, err = s.repo.UpdateUserProfile(userId, &model.UserProfile{
			Name:      updateData.Name,
			Email:     updateData.Email,
			Phone:     updateData.Phone,
			Bio:       updateData.Bio,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}

	return profile, nil
}
