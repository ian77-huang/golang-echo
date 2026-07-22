package repository

import (
	"time"

	"github.com/ian77-huang/golang-echo/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}
type sessionRepository struct {
	db *gorm.DB
}
type bibleRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	CountUser(query interface{}, args ...interface{}) (int, error)
	CreateUser(account string, password string) (*model.User, error)
	GetUser(id int) (*model.User, error)
	GetPaginate(page, pageSize int, dest any) error
	UpdateUser(id int, updateData *model.User) (*model.User, error)
	UpdateUserMap(id int, updateData map[string]interface{}) (*model.User, error)
	GetUserByAccount(account string) (*model.User, error)
	GetUserProfile(id int) (*model.UserProfile, error)
	CreateUserProfile(insertData *model.UserProfile) (*model.UserProfile, error)
	UpdateUserProfile(id int, updateData *model.UserProfile) (*model.UserProfile, error)
}
type SessionRepository interface {
	CreateSession(id, userId string, expiresAt time.Time) (*model.Session, error)
	UpdateSession(id string, expiresAt time.Time, sess *model.Session) (*model.Session, error)
	DeleteSession(id string) (*model.Session, error)
	GetSession(id string) (*model.Session, error)
}
type BibleRepository interface {
	GetBible(id int) (*model.Bible, error)
}
