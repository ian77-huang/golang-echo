package auth

import "time"

func (a *Auth[TUser, TSession]) extendDateDays(days int) time.Time {
	return time.Now().Add(time.Duration(days) * time.Hour * 24)
}
