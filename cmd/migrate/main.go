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

	dbPath := config.Databases.Path

	log.Printf("\n==== SQLite - 資料庫遷移 - 資料庫位置：%+v\n", dbPath)

	dir := filepath.Dir(dbPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("==== SQLite - 無法建立資料夾: %v\n", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Printf("==== SQLite - 資料庫遷移 - 成功\n")
}
