package model

func (User) TableName() string {
	return "users"
}

func (Session) TableName() string {
	return "session"
}
