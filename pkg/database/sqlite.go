// package database

// import (
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func NewSqlite(path string) *gorm.DB {
// 	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
// 	if err != nil {
// 		panic("=== Error：Unable to connect to the sqlite database. ====")
// 	}
// 	return db
// }

package database

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite(cfg *ConfigSqlite) (*gorm.DB, error) {
	if cfg.maxOpenConns == 0 {
		cfg.maxOpenConns = 10
	}
	if cfg.maxIdleConns == 0 {
		cfg.maxIdleConns = 1
	}
	if cfg.connMaxLifetime == 0 {
		cfg.connMaxLifetime = time.Hour
	}
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("=== Error：Unable to connect to the sqlite database: %w ====", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("=== Error：Unable to get sql.DB: %w ====", err)
	}

	_, _ = sqlDB.Exec("PRAGMA journal_mode=WAL;")

	sqlDB.SetMaxOpenConns(cfg.maxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
