package session

import "time"

// Session 代表資料庫中的 Session 資料表結構
type Session struct {
	ID          string    `gorm:"column:id;type:text;primaryKey"`
	UserID      string    `gorm:"column:userId;type:text;not null"`
	ExpiresAt   time.Time `gorm:"column:expiresAt;type:datetime;not null"`
	CreatedAt   time.Time `gorm:"column:createdAt;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updatedAt;type:datetime;not null"`
	Status      int       `gorm:"column:status;type:integer;not null;default:0"` // 0: normal, 1: logout, 99: delete
	CountUpdate int       `gorm:"column:countUpdate;type:integer;not null;default:0"`
}
