package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3" // 引入 SQLite 驅動

	appConfig "github.com/ian77-huang/golang-echo/internal/config"
)

func main() {

	log.Printf("\n==== SQLite - 資料庫遷移 - 開始\n")

	config := appConfig.Load()
	if err := runMigrations(config.Databases.Path, "file://migrations"); err != nil {
		log.Fatal(err)
	}

	log.Printf("==== SQLite - 資料庫遷移 - 成功\n")
}

func runMigrations(dbPath, migrationsURL string) error {

	log.Printf("\n==== SQLite - 資料庫遷移 - 資料庫位置：%+v\n", dbPath)

	dir := filepath.Dir(dbPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

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
