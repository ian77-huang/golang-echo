package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var signJWT = func(token *jwt.Token, key []byte) (string, error) { return token.SignedString(key) }
var parseJWT = func(tokenStr string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, claims, keyFunc)
}

func (a *Auth[TUser, TSession]) getJwtSecret() []byte {
	return []byte(a.config.SecretKey)
}
func (a *Auth[TUser, TSession]) createToken(id string, expiresAt time.Time) (string, error) {
	claims := &JwtCustomClaims{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return signJWT(token, a.getJwtSecret())
}

func (a *Auth[TUser, TSession]) parseToken(tokenStr string) (*JwtCustomClaims, error) {
	token, err := parseJWT(tokenStr, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.getJwtSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}
