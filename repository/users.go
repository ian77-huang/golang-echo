package repository

import (
	"errors"
	"time"

	"github.com/ian77-huang/golang-echo/model"

	"gorm.io/gorm"
)

func (r *userRepository) IsAccountExist(account string) (bool, error) {
	var count int64

	err := r.db.Model(&model.User{}).Where("account = ?", account).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) CreateUser(account string, password string) (*model.User, error) {
	now := time.Now()
	newUser := &model.User{Account: account, Password: password, CreatedAt: &now, UpdatedAt: &now}

	if err := r.db.Create(newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *userRepository) GetUser(id int) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Where("id = ?", id).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func (r *userRepository) GetUserByAccount(account string) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Where("account = ?", account).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
