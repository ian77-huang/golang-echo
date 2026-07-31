package service

import (
	"errors"
	"testing"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type mockUserRepository struct {
	countUserFn   func(query interface{}, args ...interface{}) (int, error)
	createUserFn  func(account string, password string) (*model.User, error)
	getUserFn     func(id int) (*model.User, error)
	getPaginateFn func(page, pageSize int, dest any) error
	updateUserFn  func(id int, updateData *model.User) (*model.User, error)
	updateUserMapFn func(id int, updateData map[string]interface{}) (*model.User, error)
	getUserByAccountFn func(account string) (*model.User, error)
	getUserProfileFn   func(id int) (*model.UserProfile, error)
	createUserProfileFn func(insertData *model.UserProfile) (*model.UserProfile, error)
	updateUserProfileFn func(id int, updateData *model.UserProfile) (*model.UserProfile, error)
}

func (m *mockUserRepository) CountUser(query interface{}, args ...interface{}) (int, error) {
	return m.countUserFn(query, args...)
}
func (m *mockUserRepository) CreateUser(account string, password string) (*model.User, error) {
	return m.createUserFn(account, password)
}
func (m *mockUserRepository) GetUser(id int) (*model.User, error) {
	return m.getUserFn(id)
}
func (m *mockUserRepository) GetPaginate(page, pageSize int, dest any) error {
	return m.getPaginateFn(page, pageSize, dest)
}
func (m *mockUserRepository) UpdateUser(id int, updateData *model.User) (*model.User, error) {
	return m.updateUserFn(id, updateData)
}
func (m *mockUserRepository) UpdateUserMap(id int, updateData map[string]interface{}) (*model.User, error) {
	return m.updateUserMapFn(id, updateData)
}
func (m *mockUserRepository) GetUserByAccount(account string) (*model.User, error) {
	return m.getUserByAccountFn(account)
}
func (m *mockUserRepository) GetUserProfile(id int) (*model.UserProfile, error) {
	return m.getUserProfileFn(id)
}
func (m *mockUserRepository) CreateUserProfile(insertData *model.UserProfile) (*model.UserProfile, error) {
	return m.createUserProfileFn(insertData)
}
func (m *mockUserRepository) UpdateUserProfile(id int, updateData *model.UserProfile) (*model.UserProfile, error) {
	return m.updateUserProfileFn(id, updateData)
}

type mockSessionRepository struct {
	createSessionFn       func(id, userId string, expiresAt time.Time) (*model.Session, error)
	updateSessionFn       func(id string, expiresAt time.Time, sess *model.Session) (*model.Session, error)
	deleteSessionFn       func(id string) (*model.Session, error)
	deleteSessionUserIdFn func(userId string) error
	getSessionFn          func(id string) (*model.Session, error)
}

func (m *mockSessionRepository) CreateSession(id, userId string, expiresAt time.Time) (*model.Session, error) {
	return m.createSessionFn(id, userId, expiresAt)
}
func (m *mockSessionRepository) UpdateSession(id string, expiresAt time.Time, sess *model.Session) (*model.Session, error) {
	return m.updateSessionFn(id, expiresAt, sess)
}
func (m *mockSessionRepository) DeleteSession(id string) (*model.Session, error) {
	return m.deleteSessionFn(id)
}
func (m *mockSessionRepository) DeleteSessionUserId(userId string) error {
	return m.deleteSessionUserIdFn(userId)
}
func (m *mockSessionRepository) GetSession(id string) (*model.Session, error) {
	return m.getSessionFn(id)
}

type mockBibleRepository struct {
	getBibleFn func(id int) (*model.Bible, error)
}

func (m *mockBibleRepository) GetBible(id int) (*model.Bible, error) {
	return m.getBibleFn(id)
}

var errMock = errors.New("mock error")

func newTestUser(t *testing.T, id int, account string) *model.User {
	t.Helper()
	now := time.Now()
	return &model.User{Id: id, Account: account, CreatedAt: &now, UpdatedAt: &now}
}

func TestUserServiceAccountQueries(t *testing.T) {
	mock := &mockUserRepository{
		countUserFn: func(query interface{}, _ ...interface{}) (int, error) {
			if query == nil {
				return 5, nil
			}
			if query == "account = ?" {
				return 1, nil
			}
			return 0, nil
		},
	}
	s := &userService{repo: mock}

	exists, err := s.IsAccountExist("tester")
	if err != nil || !exists {
		t.Fatalf("IsAccountExist() = %v, %v", exists, err)
	}
	total, err := s.CountUserAll()
	if err != nil || total != 5 {
		t.Fatalf("CountUserAll() = %d, %v", total, err)
	}

	mock.countUserFn = func(_ interface{}, _ ...interface{}) (int, error) { return 0, errMock }
	if _, err := s.IsAccountExist("tester"); err == nil {
		t.Fatal("expected count error")
	}
}

func TestUserServiceGetUser(t *testing.T) {
	mock := &mockUserRepository{
		getUserFn: func(id int) (*model.User, error) { return newTestUser(t, id, "tester"), nil },
	}
	s := &userService{repo: mock}

	user, err := s.GetUser("7")
	if err != nil || user.ID != "7" || user.Data.Account != "tester" {
		t.Fatalf("GetUser() = %#v, %v", user, err)
	}

	mock.getUserFn = func(int) (*model.User, error) { return nil, gorm.ErrRecordNotFound }
	user, err = s.GetUser("7")
	if err != nil || user != nil {
		t.Fatalf("expected nil user for not found, got %#v, %v", user, err)
	}

	mock.getUserFn = func(int) (*model.User, error) { return nil, errMock }
	if _, err := s.GetUser("7"); err == nil {
		t.Fatal("expected lookup error")
	}

	if _, err := s.GetUser("not-a-number"); err == nil {
		t.Fatal("expected cast error")
	}
}

func TestUserServiceGetUserByAccount(t *testing.T) {
	mock := &mockUserRepository{
		getUserByAccountFn: func(account string) (*model.User, error) { return newTestUser(t, 1, account), nil },
	}
	s := &userService{repo: mock}

	user, err := s.GetUserByAccount("tester")
	if err != nil || user == nil || user.Data.Account != "tester" {
		t.Fatalf("GetUserByAccount() = %#v, %v", user, err)
	}

	mock.getUserByAccountFn = func(string) (*model.User, error) { return nil, gorm.ErrRecordNotFound }
	if user, err := s.GetUserByAccount("tester"); err != nil || user != nil {
		t.Fatalf("expected nil user, got %#v, %v", user, err)
	}

	mock.getUserByAccountFn = func(string) (*model.User, error) { return nil, errMock }
	if _, err := s.GetUserByAccount("tester"); err == nil {
		t.Fatal("expected lookup error")
	}
}

func TestUserServiceCreateAndUpdate(t *testing.T) {
	mock := &mockUserRepository{
		createUserFn: func(account, password string) (*model.User, error) {
			return &model.User{Id: 1, Account: account, Password: password}, nil
		},
		updateUserFn: func(id int, updateData *model.User) (*model.User, error) {
			return newTestUser(t, id, "updated"), nil
		},
		updateUserMapFn: func(id int, updateData map[string]interface{}) (*model.User, error) {
			return newTestUser(t, id, "updated"), nil
		},
	}
	s := &userService{repo: mock}

	user, err := s.CreateUser("tester", "password")
	if err != nil || user.ID != "1" || user.Password != "password" {
		t.Fatalf("CreateUser() = %#v, %v", user, err)
	}

	updated, err := s.UpdateUser("1", &model.User{})
	if err != nil || updated == nil {
		t.Fatalf("UpdateUser() = %#v, %v", updated, err)
	}

	updatedMap, err := s.UpdateUserMap("1", map[string]interface{}{"is_active": true})
	if err != nil || updatedMap == nil {
		t.Fatalf("UpdateUserMap() = %#v, %v", updatedMap, err)
	}

	if _, err := s.UpdateUser("abc", &model.User{}); err == nil {
		t.Fatal("expected cast error")
	}
	if _, err := s.UpdateUserMap("abc", map[string]interface{}{}); err == nil {
		t.Fatal("expected cast error")
	}

	passwordUpdated, err := s.UpdateUserPassword("1", "new-hash")
	if err != nil || passwordUpdated == nil {
		t.Fatalf("UpdateUserPassword() = %#v, %v", passwordUpdated, err)
	}

	mock.createUserFn = func(string, string) (*model.User, error) { return nil, errMock }
	if _, err := s.CreateUser("tester", "password"); err == nil {
		t.Fatal("expected create error")
	}
}

func TestUserServiceProfiles(t *testing.T) {
	mock := &mockUserRepository{
		getUserProfileFn: func(id int) (*model.UserProfile, error) { return &model.UserProfile{UserID: id, Name: "Yien"}, nil },
		updateUserProfileFn: func(id int, updateData *model.UserProfile) (*model.UserProfile, error) {
			return updateData, nil
		},
		createUserProfileFn: func(insertData *model.UserProfile) (*model.UserProfile, error) {
			return insertData, nil
		},
	}
	s := &userService{repo: mock}

	profile, err := s.GetUserProfile("1")
	if err != nil || profile.Name != "Yien" {
		t.Fatalf("GetUserProfile() = %#v, %v", profile, err)
	}

	mock.getUserProfileFn = func(int) (*model.UserProfile, error) { return nil, gorm.ErrRecordNotFound }
	if profile, err := s.GetUserProfile("1"); err != nil || profile != nil {
		t.Fatalf("expected nil profile, got %#v, %v", profile, err)
	}

	mock.getUserProfileFn = func(int) (*model.UserProfile, error) { return nil, gorm.ErrRecordNotFound }
	created, err := s.UpdateUserProfile("1", &model.UserProfile{Name: "New", Email: "a@b.c"})
	if err != nil || created.Name != "New" {
		t.Fatalf("UpdateUserProfile(create) = %#v, %v", created, err)
	}

	mock.getUserProfileFn = func(int) (*model.UserProfile, error) { return &model.UserProfile{UserID: 1}, nil }
	updated, err := s.UpdateUserProfile("1", &model.UserProfile{Name: "Updated", Email: "a@b.c"})
	if err != nil || updated.Name != "Updated" {
		t.Fatalf("UpdateUserProfile(update) = %#v, %v", updated, err)
	}

	if _, err := s.GetUserProfile("abc"); err == nil {
		t.Fatal("expected cast error")
	}
	if _, err := s.UpdateUserProfile("abc", &model.UserProfile{}); err == nil {
		t.Fatal("expected cast error")
	}
}

func TestUserServiceGetPaginate(t *testing.T) {
	var gotPage, gotSize int
	mock := &mockUserRepository{
		getPaginateFn: func(page, pageSize int, dest any) error {
			gotPage, gotSize = page, pageSize
			return nil
		},
	}
	s := &userService{repo: mock}

	if _, err := s.GetPaginate(0, 0); err != nil {
		t.Fatal(err)
	}
	if gotPage != 1 || gotSize != 10 {
		t.Fatalf("expected clamped (1, 10), got (%d, %d)", gotPage, gotSize)
	}

	if _, err := s.GetPaginate(2, 500); err != nil {
		t.Fatal(err)
	}
	if gotSize != 100 {
		t.Fatalf("expected pageSize clamped to 100, got %d", gotSize)
	}

	mock.getPaginateFn = func(_, _ int, _ any) error { return gorm.ErrRecordNotFound }
	users, err := s.GetPaginate(1, 10)
	if err != nil || users != nil {
		t.Fatalf("expected nil users, got %#v, %v", users, err)
	}

	mock.getPaginateFn = func(_, _ int, _ any) error { return errMock }
	if _, err := s.GetPaginate(1, 10); err == nil {
		t.Fatal("expected paginate error")
	}
}

func TestSessionService(t *testing.T) {
	expires := time.Now().Add(time.Hour)
	mock := &mockSessionRepository{
		createSessionFn: func(id, userId string, e time.Time) (*model.Session, error) {
			return &model.Session{ID: id, UserID: userId, ExpiresAt: e}, nil
		},
		updateSessionFn: func(id string, e time.Time, sess *model.Session) (*model.Session, error) {
			return &model.Session{ID: id, UserID: sess.UserID, ExpiresAt: e, CountUpdate: sess.CountUpdate + 1}, nil
		},
		deleteSessionFn: func(id string) (*model.Session, error) {
			return &model.Session{ID: id, UserID: "user-1"}, nil
		},
		deleteSessionUserIdFn: func(userId string) error { return nil },
		getSessionFn: func(id string) (*model.Session, error) {
			return &model.Session{ID: id, UserID: "user-1", ExpiresAt: expires}, nil
		},
	}
	s := &sessionService{repo: mock}

	sess, err := s.CreateSession("s-1", "user-1", expires)
	if err != nil || sess.ID != "s-1" || sess.UserID != "user-1" {
		t.Fatalf("CreateSession() = %#v, %v", sess, err)
	}

	sess, err = s.UpdateSession("s-1", expires, &model.Session{ID: "s-1", UserID: "user-1", CountUpdate: 2})
	if err != nil || sess.Data.CountUpdate != 2 {
		t.Fatalf("UpdateSession() = %#v, %v", sess, err)
	}

	sess, err = s.GetSession("s-1")
	if err != nil || sess.Data.UserID != "user-1" {
		t.Fatalf("GetSession() = %#v, %v", sess, err)
	}

	if err := s.DeleteSessionUserId("user-1"); err != nil {
		t.Fatal(err)
	}

	sess, err = s.DeleteSession("s-1")
	if err != nil || sess.Data.UserID != "user-1" {
		t.Fatalf("DeleteSession() = %#v, %v", sess, err)
	}

	mock.deleteSessionFn = func(string) (*model.Session, error) { return nil, nil }
	if sess, err := s.DeleteSession("s-1"); err != nil || sess != nil {
		t.Fatalf("expected nil session, got %#v, %v", sess, err)
	}

	mock.createSessionFn = func(string, string, time.Time) (*model.Session, error) { return nil, errMock }
	if _, err := s.CreateSession("s-1", "user-1", expires); err == nil {
		t.Fatal("expected create error")
	}
	mock.getSessionFn = func(string) (*model.Session, error) { return nil, errMock }
	if _, err := s.GetSession("s-1"); err == nil {
		t.Fatal("expected get error")
	}
}

func TestBibleService(t *testing.T) {
	mock := &mockBibleRepository{
		getBibleFn: func(id int) (*model.Bible, error) {
			return &model.Bible{ID: id, Locale: "zh-TW", Verses: "text"}, nil
		},
	}
	s := &bibleService{repo: mock}

	bible, err := s.GetBible(3)
	if err != nil || bible.ID != 3 || bible.Verses != "text" {
		t.Fatalf("GetBible() = %#v, %v", bible, err)
	}

	mock.getBibleFn = func(int) (*model.Bible, error) { return nil, gorm.ErrRecordNotFound }
	if bible, err := s.GetBible(3); err != nil || bible != nil {
		t.Fatalf("expected nil bible, got %#v, %v", bible, err)
	}

	mock.getBibleFn = func(id int) (*model.Bible, error) {
		if id < 1 || id > 100 {
			t.Fatalf("GetBibleByDate id %d out of range", id)
		}
		return &model.Bible{ID: id}, nil
	}
	if _, err := s.GetBibleByDate(); err != nil {
		t.Fatalf("GetBibleByDate() error: %v", err)
	}
}

func TestServiceConstructors(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if NewUserService(db) == nil || NewSessionService(db) == nil || NewBibleService(db) == nil {
		t.Fatal("expected non-nil services")
	}
}
