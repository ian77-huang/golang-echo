package service

import (
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
