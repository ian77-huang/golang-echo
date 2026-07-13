package users

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func testDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestUserCRUD(t *testing.T) {
	db := testDB(t)
	created, err := CreateUser(db, "yien", "hash")
	if err != nil || created.Id == 0 {
		t.Fatalf("CreateUser() = %#v, %v", created, err)
	}
	exists, err := IsAccountExist(db, "yien")
	if err != nil || !exists {
		t.Fatalf("IsAccountExist() = %v, %v", exists, err)
	}
	byID, err := GetUser(db, created.Id)
	if err != nil || byID.Account != "yien" {
		t.Fatalf("GetUser() = %#v, %v", byID, err)
	}
	byAccount, err := GetUserByAccount(db, "yien")
	if err != nil || byAccount.Id != created.Id {
		t.Fatalf("GetUserByAccount() = %#v, %v", byAccount, err)
	}
	missing, err := GetUser(db, 99)
	if err != nil || missing != nil {
		t.Fatalf("missing user = %#v, %v", missing, err)
	}
	missingByAccount, err := GetUserByAccount(db, "missing")
	if err != nil || missingByAccount != nil {
		t.Fatalf("missing account = %#v, %v", missingByAccount, err)
	}
}

func TestUserOperationsReturnDatabaseErrors(t *testing.T) {
	db := testDB(t)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := IsAccountExist(db, "user"); err == nil {
		t.Fatal("expected account lookup database error")
	}
	if _, err := CreateUser(db, "user", "hash"); err == nil {
		t.Fatal("expected create database error")
	}
	if _, err := GetUser(db, 1); err == nil {
		t.Fatal("expected get database error")
	}
	if _, err := GetUserByAccount(db, "user"); err == nil {
		t.Fatal("expected account get database error")
	}
}
