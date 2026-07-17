package database

import "time"

type DBConfigSet struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
}
type DBType int

const (
	Sqlite DBType = iota
	Mysql
	PostgreSQL
)

type DBConfig struct {
	Driver DBType            `json:"driver"`
	Mysql  *ConfigMySQL      `json:"mysql,omitempty"`
	Pgsql  *ConfigPostgreSQL `json:"pgsql,omitempty"`
	Sqlite *ConfigSqlite     `json:"sqlite,omitempty"`
}

type ConfigPostgreSQL struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	DSN      string
	timeZone string
	DBConfigSet
}

type ConfigMySQL struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Charset  string
	DSN      string
	DBConfigSet
}

type ConfigSqlite struct {
	Path string
	DBConfigSet
}
