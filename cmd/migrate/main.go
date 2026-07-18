package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"

	appConfig "github.com/ian77-huang/golang-echo/internal/config"
	"github.com/ian77-huang/golang-echo/model"
	appAuth "github.com/ian77-huang/golang-echo/pkg/auth"
)

func main() {

	log.Printf("\n==== SQLite - 資料庫遷移 - 開始\n")

	config := appConfig.Load()

	if err := createDatabaseFolder(config.Databases.Path); err != nil {
		log.Fatal(err)
	}

	db, err := appConfig.DB(config.Databases.Path)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := runMigrations(sqlDB, "file://migrations"); err != nil {
		log.Fatal(err)
	}

	log.Printf("==== SQLite - 資料庫遷移 - 成功\n")

	if err := seedInit(db); err != nil {
		log.Fatal(err)
	}
}

func seedInit(db *gorm.DB) error {
	log.Printf("==== SQLite - Seed Init - 開始\n")

	if err := createUser(db, &model.User{Account: "test", Password: "12345678"}); err != nil {
		log.Printf("==== SQLite - Seed - User test error = %+v\n", err)
	}

	if err := createUser(db, &model.User{Account: "admin", Password: "admin12345678", IsAdmin: true}); err != nil {
		log.Printf("==== SQLite - Seed - User admin error = %+v\n", err)
	}

	log.Printf("==== SQLite - Seed Init - 完成\n")

	return nil
}

func createUser(db *gorm.DB, data *model.User) error {
	count, err := isAccountExist(db, data.Account)
	if err != nil {
		return err
	}
	if count == false {
		hashPassword, err := appAuth.HashPassword(data.Password)
		if err != nil {
			return err
		}

		now := time.Now()
		data.Password = hashPassword
		data.CreatedAt = &now

		if err := db.Create(data).Error; err != nil {
			return err
		}
	}
	return nil
}

func isAccountExist(db *gorm.DB, account string) (bool, error) {
	var count int64

	err := db.Model(&model.User{}).Where("account = ?", account).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func createDatabaseFolder(path string) error {
	log.Printf("\n==== SQLite - 資料庫遷移 - 資料庫位置：%+v\n", path)

	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}

func runMigrations(db *sql.DB, migrationsURL string) error {

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsURL,
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
