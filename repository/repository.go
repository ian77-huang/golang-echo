package repository

import "gorm.io/gorm"

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}
