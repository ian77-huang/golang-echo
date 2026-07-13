package session

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestSessionLifecycle(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&Session{}); err != nil {
		t.Fatal(err)
	}
	expires := time.Now().Add(time.Hour)
	created, err := CreateSession(db, "session-1", "1", expires)
	if err != nil || created.ID != "session-1" {
		t.Fatalf("CreateSession() = %#v, %v", created, err)
	}
	loaded, err := GetSession(db, "session-1")
	if err != nil || loaded.UserID != "1" {
		t.Fatalf("GetSession() = %#v, %v", loaded, err)
	}
	updated, err := UpdateSession(db, "session-1", expires.Add(time.Hour), loaded)
	if err != nil || updated.CountUpdate != 1 {
		t.Fatalf("UpdateSession() = %#v, %v", updated, err)
	}
	if _, err := DeleteSession(db, "1"); err != nil {
		t.Fatal(err)
	}
	if _, err := GetSession(db, "session-1"); err == nil {
		t.Fatal("expected logged-out session to be excluded")
	}
}

func TestSessionOperationsReturnDatabaseErrors(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&Session{}); err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := CreateSession(db, "id", "1", time.Now()); err == nil {
		t.Fatal("expected create error")
	}
	if _, err := UpdateSession(db, "id", time.Now(), &Session{}); err == nil {
		t.Fatal("expected update error")
	}
	if _, err := DeleteSession(db, "bad-id"); err == nil {
		t.Fatal("expected invalid id error")
	}
	if _, err := DeleteSession(db, "1"); err == nil {
		t.Fatal("expected delete database error")
	}
	if _, err := GetSession(db, "id"); err == nil {
		t.Fatal("expected get error")
	}
}
