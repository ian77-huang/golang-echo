package service

import (
	"time"

	"github.com/ian77-huang/golang-echo/repository"
)

type UserService struct {
	repo repository.UserRepository
}

type SessionService struct {
	repo repository.SessionRepository
}

type BibleService struct {
	repo repository.BibleRepository
}

type UserOmitPassword struct {
	Id        int        `json:"id"`
	Account   string     `json:"account"`
	IsActive  bool       `json:"isActive" gorm:"default:true"`
	IsAdmin   bool       `json:"isAdmin" gorm:"default:false"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
