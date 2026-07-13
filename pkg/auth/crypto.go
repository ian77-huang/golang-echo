package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

var randomRead = rand.Read

func (a *Auth[TUser, TSession]) GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := randomRead(bytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func (a *Auth[TUser, TSession]) GenerateID(token string) string {
	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}
