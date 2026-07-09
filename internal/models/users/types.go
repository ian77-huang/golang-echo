package users

import "time"

type User struct {
	Id        int        `json:"id"`
	Account   string     `json:"cccount"`
	Password  string     `json:"password"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
