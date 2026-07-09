package users

import (
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
