package service

// "github.com/ian77-huang/golang-echo/pkg/cast"

// type RegisterParameter struct {
// 	Account  string
// 	Password string
// }

// func (s *UserService) IsAccountExist(account string) (bool, error) {
// 	return s.repo.IsAccountExist(account)
// }

// func (s *UserService) CreateUser(account, password string) (*model.User, error) {
// 	if ok, err := s.IsAccountExist(account); ok {
// 		return nil, err
// 	}

// 	user, err := s.repo.CreateUser(account, password)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// func (s *UserService) GetUser(id string) (*model.User, error) {
// 	userId, err := cast.StringToInt(id, 0)
// 	if err != nil {
// 		return nil, err
// 	}
// 	user, err := s.repo.GetUser(userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// func (s *UserService) GetUserByAccount(account string) (*model.User, error) {
// 	user, err := s.repo.GetUserByAccount(account)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }
