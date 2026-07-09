package auth

import "time"

func (a *Auth[TUser, TSession]) CalculateExpiry(days int) time.Time {
	return time.Now().AddDate(0, 0, days)
}
