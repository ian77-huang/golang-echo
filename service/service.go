package service

import (
	"github.com/ian77-huang/golang-echo/repository"
	"gorm.io/gorm"
)

func NewUserService(db *gorm.DB) UserService {
	return &userService{repo: repository.NewUserRepository(db)}
}

func NewSessionService(db *gorm.DB) SessionService {
	return &sessionService{repo: repository.NewSessionRepository(db)}
}

func NewBibleService(db *gorm.DB) BibleService {
	return &bibleService{repo: repository.NewBibleRepository(db)}
}
