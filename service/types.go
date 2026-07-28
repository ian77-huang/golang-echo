package service

import (
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"github.com/ian77-huang/golang-echo/repository"
)

type UserService interface {
	IsAccountExist(account string) (bool, error)
	CountUserAll() (int, error)
	GetUser(id string) (*appAuth.User[model.User], error)
	GetPaginate(page, pageSize int) ([]UserOmitPassword, error)
	GetUserByAccount(account string) (*appAuth.User[model.User], error)
	CreateUser(account, password string) (*appAuth.User[model.User], error)
	UpdateUser(id string, updateData *model.User) (*appAuth.User[model.User], error)
	UpdateUserMap(id string, updateData map[string]interface{}) (*appAuth.User[model.User], error)
	UpdateUserPassword(id string, passwordHash string) (*appAuth.User[model.User], error)
	GetUserProfile(id string) (*model.UserProfile, error)
	UpdateUserProfile(id string, updateData *model.UserProfile) (*model.UserProfile, error)
}
type userService struct {
	repo repository.UserRepository
}

type SessionService interface {
	CreateSession(id, userId string, expiresAt time.Time) (*appAuth.Session[model.Session], error)
	UpdateSession(id string, expiresAt time.Time, sess *model.Session) (*appAuth.Session[model.Session], error)
	DeleteSession(id string) (*appAuth.Session[model.Session], error)
	GetSession(id string) (*appAuth.Session[model.Session], error)
}
type sessionService struct {
	repo repository.SessionRepository
}

type BibleService interface {
	GetBible(id int) (*model.Bible, error)
	GetBibleByDate() (*model.Bible, error)
}
type bibleService struct {
	repo repository.BibleRepository
	g    singleflight.Group
}

type UserOmitPassword struct {
	Id        int        `json:"id"`
	Account   string     `json:"account"`
	IsActive  bool       `json:"isActive" gorm:"default:true"`
	IsAdmin   bool       `json:"isAdmin" gorm:"default:false"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
