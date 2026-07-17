package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysql(cfg *ConfigMySQL) (*gorm.DB, error) {
	dsn := cfg.DSN
	if dsn == "" {
		if cfg.Charset == "" {
			cfg.Charset = "utf8mb4"
		}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DBName,
			cfg.Charset,
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
