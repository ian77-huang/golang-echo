package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresSql(cfg *ConfigPostgreSQL) (*gorm.DB, error) {
	dsn := cfg.DSN
	if dsn == "" {
		if cfg.SSLMode == "" {
			cfg.SSLMode = "disable"
		}
		if cfg.timeZone == "" {
			cfg.timeZone = "Asia/Taipei"
		}
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			cfg.Host,
			cfg.User,
			cfg.Password,
			cfg.DBName,
			cfg.Port,
			cfg.SSLMode,
			cfg.timeZone,
		)
	}
	if cfg.maxIdleConns == 0 {
		cfg.maxIdleConns = 10
	}
	if cfg.maxOpenConns == 0 {
		cfg.maxOpenConns = 100
	}
	if cfg.connMaxLifetime == 0 {
		cfg.connMaxLifetime = time.Hour
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("=== Error：Unable to connect to the mysql database: %w ====", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("=== Error：Unable to get sql.DB: %w ====", err)
	}

	sqlDB.SetMaxIdleConns(cfg.maxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.maxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.connMaxLifetime)

	return db, nil
}
