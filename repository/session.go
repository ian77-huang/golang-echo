package repository

import (
	"time"

	"github.com/ian77-huang/golang-echo/pkg/cast"

	"github.com/ian77-huang/golang-echo/model"
)

func (r *sessionRepository) CreateSession(id, userId string, expiresAt time.Time) (*model.Session, error) {
	newSession := &model.Session{
		ID: id, UserID: userId, ExpiresAt: expiresAt, UpdatedAt: time.Now(), Status: 0, CountUpdate: 0,
	}
	tx := r.db.Create(newSession)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return newSession, nil
}

func (r *sessionRepository) UpdateSession(id string, expiresAt time.Time, sess *model.Session) (*model.Session, error) {
	updateSession := &model.Session{
		ExpiresAt:   expiresAt,
		UpdatedAt:   time.Now(),
		CountUpdate: sess.CountUpdate + 1,
	}

	tx := r.db.Model(&model.Session{}).Where("id = ?", id).Where("status = ?", 0).Updates(&updateSession)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return updateSession, nil
}

func (r *sessionRepository) DeleteSession(id string) (*model.Session, error) {
	deleteSession := &model.Session{UpdatedAt: time.Now(), Status: 1}

	ID, err := cast.StringToInt(id, 0)
	if err != nil {
		return nil, err
	}

	tx := r.db.Model(&model.Session{}).Where("userId = ?", ID).Updates(deleteSession)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return deleteSession, nil
}

func (r *sessionRepository) GetSession(id string) (*model.Session, error) {
	getSession := &model.Session{}
	tx := r.db.Model(&model.Session{}).Where("id = ?", id).Where("status = 0").First(getSession)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return getSession, nil
}
