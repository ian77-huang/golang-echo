package session

import (
	"time"

	"github.com/ian77-huang/golang-echo/pkg/cast"

	"gorm.io/gorm"
)

func (Session) TableName() string {
	return "session"
}

func CreateSession(db *gorm.DB, id, userId string, expiresAt time.Time) (*Session, error) {
	newSession := &Session{
		ID: id, UserID: userId, ExpiresAt: expiresAt, UpdatedAt: time.Now(), Status: 0, CountUpdate: 0,
	}
	tx := db.Create(newSession)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return newSession, nil
}

func UpdateSession(db *gorm.DB, id string, expiresAt time.Time, sess *Session) (*Session, error) {
	updateSession := &Session{
		ExpiresAt:   expiresAt,
		UpdatedAt:   time.Now(),
		CountUpdate: sess.CountUpdate + 1,
	}

	tx := db.Model(&Session{}).Where("id = ?", id).Where("status = ?", 0).Updates(&updateSession)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return updateSession, nil
}

func DeleteSession(db *gorm.DB, id string) (*Session, error) {
	deleteSession := &Session{UpdatedAt: time.Now(), Status: 1}

	ID, err := cast.StringToInt(id, 0)
	if err != nil {
		return nil, err
	}

	tx := db.Model(&Session{}).Where("userId = ?", ID).Updates(deleteSession)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return deleteSession, nil
}

func GetSession(db *gorm.DB, id string) (*Session, error) {
	getSession := &Session{}
	tx := db.Model(&Session{}).Where("id = ?", id).Where("status = 0").First(getSession)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return getSession, nil
}
