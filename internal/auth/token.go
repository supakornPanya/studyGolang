package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represent the data inseid JWT
type Claims struct {
	Username  string `json:"username"`
	UserID    string `json:"user_id"`
	CanRead   bool   `json:"can_read"`
	CanWrite  bool   `json:"can_write"`
	CanUpdate bool   `json:"can_update"`
	CanDelete bool   `json:"can_delete"`
	jwt.RegisteredClaims
}

// GenerateToken
func GenerateToken(username, userID string, canRead, canWrite, canUpdate, canDelete bool, secret []byte) (string, error) {
	claims := &Claims{
		Username: username,
		UserID:   userID,
		CanRead:   canRead,
		CanWrite:  canWrite,
		CanUpdate: canUpdate,
		CanDelete: canDelete,
		RegisteredClaims: jwt.RegisteredClaims{
			//Set expire token
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	//Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

//Validation token
func ValidateToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}