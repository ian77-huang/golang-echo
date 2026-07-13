package auth

import "github.com/ian77-huang/golang-echo/pkg/argon2"

func HashPassword(password string) (string, error) {
	passwordHash, err := argon2.HashPassword(password)
	if err != nil {
		return "", NewError("error.auth.FailedToSecurePassword", "failed to secure password")
	}
	return passwordHash, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	verify, err := argon2.VerifyPassword(password, encodedHash)
	if err != nil {
		return false, NewError("error.auth.FailedToSecurePassword", "failed to secure password")
	}
	return verify, nil
}
