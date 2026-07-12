package users

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

func (User) TableName() string {
	return "users"
}

func IsAccountExist(db *gorm.DB, account string) (bool, error) {
	var count int64

	err := db.Model(&User{}).Where("account = ?", account).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CreateUser(db *gorm.DB, account string, password string) (*User, error) {
	now := time.Now()
	newUser := &User{Account: account, Password: password, CreatedAt: &now, UpdatedAt: &now}

	if err := db.Create(newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil
}

func GetUser(db *gorm.DB, id int) (*User, error) {
	user := &User{}
	if err := db.Where("id = ?", id).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func GetUserByAccount(db *gorm.DB, account string) (*User, error) {
	user := &User{}
	if err := db.Where("account = ?", account).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
