package service

import (
	"errors"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/pkg/cast"
	"gorm.io/gorm"
)

func (s *userService) IsAccountExist(account string) (bool, error) {
	count, err := s.repo.CountUser("account = ?", account)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *userService) CountUserAll() (int, error) {
	return s.repo.CountUser(nil, nil)
}

func (s *userService) GetUser(id string) (*appAuth.User[model.User], error) {
	userId, err := cast.StringToInt(id, 0)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.GetUser(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	authUser := &appAuth.User[model.User]{ID: cast.IntToString(user.Id), Data: user, Password: user.Password}
	return authUser, nil
}

func (s *userService) GetPaginate(page, pageSize int) ([]UserOmitPassword, error) {
	if page <= 0 {
		page = 1
	}

	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	var users []UserOmitPassword
	err := s.repo.GetPaginate(page, pageSize, &users)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userService) GetUserByAccount(account string) (*appAuth.User[model.User], error) {
	user, err := s.repo.GetUserByAccount(account)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	authUser := &appAuth.User[model.User]{ID: cast.IntToString(user.Id), Data: user, Password: user.Password}
	return authUser, nil
}

func (s *userService) CreateUser(account, password string) (*appAuth.User[model.User], error) {
	user, err := s.repo.CreateUser(account, password)
	if err != nil {
		return nil, err
	}
	authUser := &appAuth.User[model.User]{ID: cast.IntToString(user.Id), Password: user.Password, Data: user}
	return authUser, nil
}
func (s *userService) UpdateUser(id string, updateData *model.User) (*appAuth.User[model.User], error) {
	userId, err := cast.StringToInt(id, 0)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.UpdateUser(userId, updateData)
	if err != nil {
		return nil, err
	}
	authUser := &appAuth.User[model.User]{ID: cast.IntToString(user.Id), Password: user.Password, Data: user}
	return authUser, nil
}
func (s *userService) UpdateUserMap(id string, updateData map[string]interface{}) (*appAuth.User[model.User], error) {
	userId, err := cast.StringToInt(id, 0)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.UpdateUserMap(userId, updateData)
	if err != nil {
		return nil, err
	}
	authUser := &appAuth.User[model.User]{ID: cast.IntToString(user.Id), Password: user.Password, Data: user}
	return authUser, nil
}
func (s *userService) UpdateUserPassword(id string, passwordHash string) (*appAuth.User[model.User], error) {
	updateData := &model.User{Password: passwordHash}
	user, err := s.UpdateUser(id, updateData)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserProfile(id string) (*model.UserProfile, error) {
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

func (s *userService) UpdateUserProfile(id string, updateData *model.UserProfile) (*model.UserProfile, error) {
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
			AvatarURL: updateData.AvatarURL,
			UpdatedAt: time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}

	return profile, nil
}
