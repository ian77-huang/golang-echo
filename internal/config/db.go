package config

import (
	"github.com/ian77-huang/golang-echo/pkg/database"
	"gorm.io/gorm"
)

func DB(path string) (*gorm.DB, error) {
	dbConfig := &database.DBConfig{
		Driver: database.Sqlite, Sqlite: &database.ConfigSqlite{
			Path: path,
		}}
	return database.New(dbConfig)
}
