package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqlite(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("=== Error：Unable to connect to the database. ====")
	}
	return db
}
