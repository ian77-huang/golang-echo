package service

import (
	"time"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
)

func (s *sessionService) CreateSession(id, userId string, expiresAt time.Time) (*appAuth.Session[model.Session], error) {
	sess, err := s.repo.CreateSession(id, userId, expiresAt)
	if err != nil {
		return nil, err
	}

	return (&appAuth.Session[model.Session]{ID: sess.ID, UserID: userId, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
}
func (s *sessionService) UpdateSession(id string, expiresAt time.Time, sess *model.Session) (*appAuth.Session[model.Session], error) {
	updateSess, err := s.repo.UpdateSession(id, expiresAt, sess)
	if err != nil {
		return nil, err
	}
	return (&appAuth.Session[model.Session]{ID: sess.ID, UserID: updateSess.UserID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
}
func (s *sessionService) DeleteSession(id string) (*appAuth.Session[model.Session], error) {
	sess, err := s.repo.DeleteSession(id)
	if sess == nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return (&appAuth.Session[model.Session]{ID: sess.ID, UserID: sess.UserID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
}
func (s *sessionService) DeleteSessionUserId(userId string) error {
	err := s.repo.DeleteSessionUserId(userId)
	if err != nil {
		return err
	}
	return nil
}
func (s *sessionService) GetSession(id string) (*appAuth.Session[model.Session], error) {
	sess, err := s.repo.GetSession(id)
	if err != nil {
		return nil, err
	}
	return (&appAuth.Session[model.Session]{ID: sess.ID, UserID: sess.UserID, ExpiresAt: sess.ExpiresAt, Data: sess}), nil
}
