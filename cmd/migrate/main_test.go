package main

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestRunMigrationsAppliesProjectMigrations(t *testing.T) {
	migrationsDir, err := filepath.Abs(filepath.Join("..", "..", "migrations"))
	if err != nil {
		t.Fatal(err)
	}
	dbPath := filepath.Join(t.TempDir(), "nested", "app.db")
	if err := runMigrations(dbPath, "file://"+migrationsDir); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var name string
	if err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users'").Scan(&name); err != nil || name != "users" {
		t.Fatalf("users table: %q %v", name, err)
	}
	if err := runMigrations(dbPath, "file://"+migrationsDir); err != nil {
		t.Fatalf("second migration should be idempotent: %v", err)
	}
}

func TestRunMigrationsReturnsInvalidSourceError(t *testing.T) {
	if err := runMigrations(filepath.Join(t.TempDir(), "app.db"), "file:///not-a-real-migrations-directory"); err == nil {
		t.Fatal("expected migration source error")
	}
}
