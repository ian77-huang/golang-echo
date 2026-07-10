package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (a *Auth[TUser, TSession]) getJwtSecret() []byte {
	return []byte(a.config.SecretKey)
}
func (a *Auth[TUser, TSession]) createToken(id string) (string, error) {
	claims := &JwtCustomClaims{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(a.getJwtSecret())
}
