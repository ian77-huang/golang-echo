package service

// import (
// 	"time"

// 	"github.com/ian77-huang/golang-echo/model"
// )

// func (s *SessionService) CreateSession(id, userId string, expiresAt time.Time) (*model.Session, error) {
// 	sess, err := s.repo.CreateSession(id, userId, expiresAt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return sess, nil
// }
// func (s *SessionService) UpdateSession(id string, expiresAt time.Time, sess *model.Session) (*model.Session, error) {
// 	updateSess, err := s.repo.UpdateSession(id, expiresAt, sess)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return updateSess, nil
// }
// func (s *SessionService) DeleteSession(id string) (*model.Session, error) {
// 	sess, err := s.repo.DeleteSession(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return sess, nil
// }
// func (s *SessionService) GetSession(id string) (*model.Session, error) {
// 	sess, err := s.repo.GetSession(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return sess, nil
// }
