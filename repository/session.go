package repository

import (
	"errors"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	"gorm.io/gorm"
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
	var sess *model.Session
	txSess := r.db.Model(&model.Session{}).Where("id = ?", id).First(&sess)

	if txSess.Error != nil && !errors.Is(txSess.Error, gorm.ErrRecordNotFound) {
		return nil, txSess.Error
	} else {
		if sess == nil {
			return nil, nil
		}
		deleteSession := &model.Session{UpdatedAt: time.Now(), Status: 1}
		tx := r.db.Model(&model.Session{}).Where("userId = ?", sess.UserID).Updates(deleteSession)
		if tx.Error != nil {
			return nil, tx.Error
		}

		return sess, nil
	}
}

func (r *sessionRepository) GetSession(id string) (*model.Session, error) {
	getSession := &model.Session{}
	tx := r.db.Model(&model.Session{}).Where("id = ?", id).Where("status = 0").First(getSession)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return getSession, nil
}
