package model

func (User) TableName() string {
	return "user"
}
func (UserProfile) TableName() string {
	return "user_profile"
}

func (Session) TableName() string {
	return "session"
}

func (Bible) TableName() string {
	return "bible"
}
