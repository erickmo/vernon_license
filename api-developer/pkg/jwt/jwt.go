// Package jwt menyediakan fungsi sign dan verify JWT HS256 untuk Vernon App.
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims adalah payload JWT untuk Vernon App.
type Claims struct {
	Sub  string `json:"sub"`  // user UUID
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// Sign menghasilkan JWT string dari claims.
// Expiry: 24 jam dari sekarang.
func Sign(userID, name, role, secret string) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		Sub:  userID,
		Name: name,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("jwt.Sign: failed to sign token")
	}
	return signed, nil
}

// Verify mem-parse dan memvalidasi JWT string.
// Returns Claims jika valid, error jika expired atau invalid.
func Verify(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("jwt.Verify: unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("jwt.Verify: invalid token")
	}
	return claims, nil
}
