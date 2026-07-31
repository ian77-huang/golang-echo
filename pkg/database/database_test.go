package database

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
)

func TestMiddlewareStoresDatabase(t *testing.T) {
	db, err := NewSqlite(&ConfigSqlite{Path: ":memory:"})
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	err = Middleware(db)(func(c *echo.Context) error {
		got, err := GetDBConnect(c)
		if err != nil || got != db {
			t.Fatalf("GetDBConnect() = %v, %v", got, err)
		}
		return nil
	})(c)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetDBConnectRequiresValidValue(t *testing.T) {
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	if _, err := GetDBConnect(c); err == nil {
		t.Fatal("expected missing database error")
	}
	c.Set(contextDBKey, "not-a-db")
	if _, err := GetDBConnect(c); err == nil {
		t.Fatal("expected invalid database error")
	}
}

func TestNewSqliteAppliesDefaults(t *testing.T) {
	cfg := &ConfigSqlite{Path: ":memory:"}
	db, err := NewSqlite(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.maxOpenConns != 10 || cfg.maxIdleConns != 1 || cfg.connMaxLifetime != time.Hour {
		t.Fatalf("unexpected defaults: %#v", cfg)
	}
	var one int
	if err := db.Raw("SELECT 1").Scan(&one).Error; err != nil || one != 1 {
		t.Fatalf("connection unusable: %v", err)
	}
}

func TestNewSqliteReturnsUnusablePathError(t *testing.T) {
	if _, err := NewSqlite(&ConfigSqlite{Path: "/dev/null/app.db"}); err == nil {
		t.Fatal("expected sqlite path error")
	}
}

func TestNewBuildsRequestedDrivers(t *testing.T) {
	db, err := New(&DBConfig{Driver: Sqlite, Sqlite: &ConfigSqlite{Path: ":memory:"}})
	if err != nil || db == nil {
		t.Fatalf("sqlite driver: db=%v err=%v", db, err)
	}
}

func TestNewRejectsMissingDriverConfigs(t *testing.T) {
	for _, cfg := range []*DBConfig{
		{Driver: Sqlite},
		{Driver: PostgreSQL},
		{Driver: Mysql},
		{Driver: DBType(99)},
	} {
		if _, err := New(cfg); err == nil {
			t.Fatalf("expected error for %#v", cfg)
		}
	}
}

func TestNewMysqlDefaultsAndConnectionError(t *testing.T) {
	cfg := &ConfigMySQL{User: "u", Password: "p", Host: "localhost", Port: 1, DBName: "d"}
	if _, err := NewMysql(cfg); err == nil {
		t.Fatal("expected mysql connection error")
	}
	if cfg.Charset != "utf8mb4" || cfg.maxIdleConns != 10 || cfg.maxOpenConns != 100 || cfg.connMaxLifetime != time.Hour {
		t.Fatalf("unexpected defaults: %#v", cfg)
	}
}

func TestNewMysqlReturnsInvalidDSNError(t *testing.T) {
	if _, err := NewMysql(&ConfigMySQL{DSN: "bad dsn"}); err == nil {
		t.Fatal("expected invalid dsn error")
	}
}

func TestNewPostgresSqlDefaultsAndConnectionError(t *testing.T) {
	cfg := &ConfigPostgreSQL{User: "u", Password: "p", Host: "localhost", Port: 1, DBName: "d"}
	if _, err := NewPostgresSql(cfg); err == nil {
		t.Fatal("expected postgres connection error")
	}
	if cfg.SSLMode != "disable" || cfg.timeZone != "Asia/Taipei" || cfg.maxIdleConns != 10 || cfg.maxOpenConns != 100 || cfg.connMaxLifetime != time.Hour {
		t.Fatalf("unexpected defaults: %#v", cfg)
	}
}

func TestNewPostgresSqlReturnsInvalidDSNError(t *testing.T) {
	if _, err := NewPostgresSql(&ConfigPostgreSQL{DSN: "bad dsn"}); err == nil {
		t.Fatal("expected invalid dsn error")
	}
}
