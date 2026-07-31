package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/ian77-huang/golang-echo/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.UserProfile{}, &model.Session{}, &model.Bible{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestUserRepositoryCRUD(t *testing.T) {
	db := testDB(t)
	repo := NewUserRepository(db)

	if _, err := repo.GetUser(1); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}

	created, err := repo.CreateUser("tester", "hash")
	if err != nil {
		t.Fatal(err)
	}
	if created.Id == 0 || created.Account != "tester" || created.Password != "hash" {
		t.Fatalf("unexpected user: %#v", created)
	}

	got, err := repo.GetUser(created.Id)
	if err != nil || got.Account != "tester" {
		t.Fatalf("GetUser() = %#v, %v", got, err)
	}

	byAccount, err := repo.GetUserByAccount("tester")
	if err != nil || byAccount.Id != created.Id {
		t.Fatalf("GetUserByAccount() = %#v, %v", byAccount, err)
	}

	count, err := repo.CountUser("account = ?", "tester")
	if err != nil || count != 1 {
		t.Fatalf("CountUser() = %d, %v", count, err)
	}
	countAll, err := repo.CountUser(nil, nil)
	if err != nil || countAll != 1 {
		t.Fatalf("CountUser(all) = %d, %v", countAll, err)
	}

	var page []model.User
	if err := repo.GetPaginate(1, 10, &page); err != nil {
		t.Fatal(err)
	}
	if len(page) != 1 {
		t.Fatalf("pagination returned %d users", len(page))
	}

	updated, err := repo.UpdateUser(created.Id, &model.User{IsAdmin: true})
	if err != nil || !updated.IsAdmin {
		t.Fatalf("UpdateUser() = %#v, %v", updated, err)
	}

	updatedMap, err := repo.UpdateUserMap(created.Id, map[string]interface{}{"is_active": false})
	if err != nil || updatedMap.IsActive {
		t.Fatalf("UpdateUserMap() = %#v, %v", updatedMap, err)
	}

	if _, err := repo.UpdateUserMap(9999, map[string]interface{}{"is_active": false}); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found for missing user, got %v", err)
	}
}

func TestUserProfileRepository(t *testing.T) {
	db := testDB(t)
	repo := NewUserRepository(db)

	if _, err := repo.GetUserProfile(1); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}

	profile, err := repo.CreateUserProfile(&model.UserProfile{
		UserID: 1, Name: "Yien", Email: "yien@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if profile.UserID != 1 {
		t.Fatalf("unexpected profile: %#v", profile)
	}

	got, err := repo.GetUserProfile(1)
	if err != nil || got.Name != "Yien" {
		t.Fatalf("GetUserProfile() = %#v, %v", got, err)
	}

	updated, err := repo.UpdateUserProfile(1, &model.UserProfile{Name: "Updated", Email: "yien@example.com"})
	if err != nil || updated.Name != "Updated" {
		t.Fatalf("UpdateUserProfile() = %#v, %v", updated, err)
	}
}

func TestSessionRepository(t *testing.T) {
	db := testDB(t)
	repo := NewSessionRepository(db)

	expires := time.Now().Add(time.Hour)
	sess, err := repo.CreateSession("session-1", "user-1", expires)
	if err != nil {
		t.Fatal(err)
	}
	if sess.Status != 0 || sess.CountUpdate != 0 {
		t.Fatalf("unexpected session: %#v", sess)
	}

	got, err := repo.GetSession("session-1")
	if err != nil || got.UserID != "user-1" {
		t.Fatalf("GetSession() = %#v, %v", got, err)
	}

	newExpires := time.Now().Add(24 * time.Hour)
	updated, err := repo.UpdateSession("session-1", newExpires, sess)
	if err != nil || updated.CountUpdate != 1 {
		t.Fatalf("UpdateSession() = %#v, %v", updated, err)
	}
	if !updated.ExpiresAt.After(expires) {
		t.Fatal("expected extended expiration")
	}

	if err := repo.DeleteSessionUserId("user-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetSession("session-1"); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected session to be deactivated, got %v", err)
	}

	if _, err := repo.CreateSession("session-2", "user-2", expires); err != nil {
		t.Fatal(err)
	}
	deleted, err := repo.DeleteSession("session-2")
	if err != nil || deleted.UserID != "user-2" {
		t.Fatalf("DeleteSession() = %#v, %v", deleted, err)
	}

	// Deleting a missing session does not error.
	missing, err := repo.DeleteSession("missing")
	if err != nil {
		t.Fatalf("DeleteSession(missing) error = %v", err)
	}
	_ = missing
}

func TestBibleRepository(t *testing.T) {
	db := testDB(t)
	repo := NewBibleRepository(db)

	if _, err := repo.GetBible(1); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found, got %v", err)
	}

	bible := &model.Bible{Locale: "zh-TW", Verses: "text"}
	if err := db.Create(bible).Error; err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetBible(bible.ID)
	if err != nil || got.Verses != "text" {
		t.Fatalf("GetBible() = %#v, %v", got, err)
	}
}

func TestRepositoryConstructors(t *testing.T) {
	db := testDB(t)
	if NewUserRepository(db) == nil || NewSessionRepository(db) == nil || NewBibleRepository(db) == nil {
		t.Fatal("expected non-nil repositories")
	}
}
