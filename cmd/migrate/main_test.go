package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(filepath.Join(t.TempDir(), "app.db")), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestRunMigrationsAppliesProjectMigrations(t *testing.T) {
	migrationsDir, err := filepath.Abs(filepath.Join("..", "..", "migrations"))
	if err != nil {
		t.Fatal(err)
	}
	dbPath := filepath.Join(t.TempDir(), "nested", "app.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := runMigrations(db, "file://"+migrationsDir); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"user", "messages", "session", "user_profile", "bible"} {
		var name string
		if err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name); err != nil || name != table {
			t.Fatalf("%s table: %q %v", table, name, err)
		}
	}
	if err := runMigrations(db, "file://"+migrationsDir); err != nil {
		t.Fatalf("second migration should be idempotent: %v", err)
	}
}

func TestRunMigrationsReturnsInvalidSourceError(t *testing.T) {
	db, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "app.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := runMigrations(db, "file:///not-a-real-migrations-directory"); err == nil {
		t.Fatal("expected migration source error")
	}
}

func TestCreateDatabaseFolderCreatesParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "app.db")
	if err := createDatabaseFolder(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Fatal(err)
	}
}

func TestIsAccountExistAndCreateUser(t *testing.T) {
	db := openTestDB(t)
	exists, err := isAccountExist(db, "test")
	if err != nil || exists {
		t.Fatalf("fresh account: exists=%v err=%v", exists, err)
	}
	user := &model.User{Account: "test", Password: "12345678"}
	if err := createUser(db, user); err != nil {
		t.Fatal(err)
	}
	exists, err = isAccountExist(db, "test")
	if err != nil || !exists {
		t.Fatalf("after create: exists=%v err=%v", exists, err)
	}
	var saved model.User
	if err := db.Where("account = ?", "test").First(&saved).Error; err != nil {
		t.Fatal(err)
	}
	if saved.Password == "12345678" {
		t.Fatal("expected password to be hashed")
	}
	if ok, _ := appAuth.VerifyPassword("12345678", saved.Password); !ok {
		t.Fatal("expected password hash to verify")
	}
	if err := createUser(db, user); err != nil {
		t.Fatal(err)
	}
	var count int64
	if err := db.Model(&model.User{}).Where("account = ?", "test").Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected single user, got %d", count)
	}
}

func TestCreateUserPropagatesDatabaseError(t *testing.T) {
	db := openTestDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.Close()
	if err := createUser(db, &model.User{Account: "test"}); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSeedInitCreatesDefaultUsers(t *testing.T) {
	db := openTestDB(t)
	if err := seedInit(db); err != nil {
		t.Fatal(err)
	}
	var count int64
	if err := db.Model(&model.User{}).Where("account IN ?", []string{"test", "admin"}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("expected seeded users, got %d", count)
	}
	if err := seedInit(db); err != nil {
		t.Fatal(err)
	}
	if err := db.Model(&model.User{}).Where("account IN ?", []string{"test", "admin"}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("expected seed to be idempotent, got %d", count)
	}
}

func TestMainRunsMigrationsAndSeed(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	t.Setenv("DATABASE_PATH", filepath.Join(t.TempDir(), "app.db"))
	t.Setenv("SECRET_KEY", "test-secret")

	main()

	db, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_PATH")), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	var count int64
	if err := db.Model(&model.User{}).Where("account = ?", "admin").Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected seeded admin user, got %d", count)
	}
}
