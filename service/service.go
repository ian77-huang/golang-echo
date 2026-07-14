package service

import (
	"github.com/ian77-huang/golang-echo/repository"
	"gorm.io/gorm"
)

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{repo: repository.NewUserRepository(db)}
}

func NewSessionService(db *gorm.DB) *SessionService {
	return &SessionService{repo: repository.NewSessionRepository(db)}
}
