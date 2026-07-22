package repository

import (
	"time"

	"github.com/ian77-huang/golang-echo/model"
	"gorm.io/gorm"
)

func (r *userRepository) CountUser(query interface{}, args ...interface{}) (int, error) {
	var count int64

	tx := r.db.Model(&model.User{})

	if query != nil && query != "" {
		tx = tx.Where(query, args...)
	}

	err := tx.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
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
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetPaginate(page, pageSize int, dest any) error {
	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.User{}).Offset(offset).Limit(pageSize).Find(dest).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUser(id int, updateData *model.User) (*model.User, error) {
	if err := r.db.Model(&model.User{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		return nil, err
	}
	return updateData, nil
}
func (r *userRepository) UpdateUserMap(id int, updateData map[string]interface{}) (*model.User, error) {
	result := r.db.Model(&model.User{}).Where("id = ?", id).Updates(updateData)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var updatedUser model.User
	if err := r.db.First(&updatedUser, id).Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (r *userRepository) GetUserByAccount(account string) (*model.User, error) {
	user := &model.User{}
	if err := r.db.Where("account = ?", account).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetUserProfile(id int) (*model.UserProfile, error) {
	userProfile := &model.UserProfile{}
	if err := r.db.Where("user_id = ?", id).First(userProfile).Error; err != nil {
		return nil, err
	}
	return userProfile, nil
}

func (r *userRepository) CreateUserProfile(insertData *model.UserProfile) (*model.UserProfile, error) {
	tx := r.db.Model(&model.UserProfile{}).Create(insertData)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return insertData, nil
}
func (r *userRepository) UpdateUserProfile(id int, updateData *model.UserProfile) (*model.UserProfile, error) {
	tx := r.db.Model(&model.UserProfile{}).Where("user_id = ?", id).Updates(updateData)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return updateData, nil
}
