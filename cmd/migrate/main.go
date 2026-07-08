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
)

func main() {
	dbPath := "./databases/main.db"

	// 💡 核心步驟：自動取得資料夾路徑（./databases），並直接建立它
	dir := filepath.Dir(dbPath)
	err := os.MkdirAll(dir, 0755) // 0755 是標準的資料夾權限
	if err != nil {
		log.Fatalf("無法建立資料夾: %v", err)
	}

	// 1. 連線至 SQLite 資料庫（這裡會在本機建立一個名為 mydb.db 的檔案）
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. 建立 migrate 的 SQLite 驅動個體
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 3. 指定 SQL 遷移檔案的路徑（預設讀取專案內的 ./migrations 目錄）
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 執行升級（Up）
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("SQLite 資料庫遷移成功！")

	// 5. 接下來啟動你的 Web 伺服器...
}
