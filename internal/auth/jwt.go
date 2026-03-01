package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Token struct {
	Value          string
	ExpirationTime time.Time
}

func GenerateToken(userId uuid.UUID, secretKey string, expireHours int) (*Token, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(expireHours)).Local().UTC()

	claims := jwt.MapClaims{
		"userId": userId.String(),
		"exp":    expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &Token{
		Value:          signed,
		ExpirationTime: expirationTime,
	}, nil
}
