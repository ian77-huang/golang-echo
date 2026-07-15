package model

import "time"

type User struct {
	Id        int        `json:"id"`
	Account   string     `json:"cccount"`
	Password  string     `json:"password"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
type UserProfile struct {
	UserID    int       `gorm:"primaryKey" json:"user_id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	Phone     string    `gorm:"type:varchar(50)" json:"phone"`
	Bio       string    `gorm:"type:text" json:"bio"`
	AvatarURL string    `gorm:"type:varchar(512)" json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	ID          string    `gorm:"column:id;type:text;primaryKey"`
	UserID      string    `gorm:"column:userId;type:text;not null"`
	ExpiresAt   time.Time `gorm:"column:expiresAt;type:datetime;not null"`
	CreatedAt   time.Time `gorm:"column:createdAt;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updatedAt;type:datetime;not null"`
	Status      int       `gorm:"column:status;type:integer;not null;default:0"` // 0: normal, 1: logout, 99: delete
	CountUpdate int       `gorm:"column:countUpdate;type:integer;not null;default:0"`
}
