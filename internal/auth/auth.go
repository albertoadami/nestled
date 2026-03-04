package auth

import (
	"fmt"
	"time"

	"github.com/albertoadami/nestled/internal/config"
	"github.com/albertoadami/nestled/internal/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Token struct {
	Value          string
	ExpirationTime time.Time
}

type TokenInfo struct {
	UserId         string
	ExpirationTime time.Time
}

type TokenManager struct {
	jwtConfig config.JWTConfig
}

func NewTokenManager(jwtConfig config.JWTConfig) *TokenManager {
	return &TokenManager{jwtConfig: jwtConfig}
}

func (tm *TokenManager) GenerateToken(userId uuid.UUID) (*Token, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(tm.jwtConfig.Expiration)).Local().UTC()

	claims := jwt.MapClaims{
		"userId": userId.String(),
		"exp":    expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(tm.jwtConfig.Secret))

	if err != nil {
		return nil, errors.ErrGeneratingToken
	}

	return &Token{
		Value:          signed,
		ExpirationTime: expirationTime,
	}, nil
}

func (tm *TokenManager) ParseToken(tokenValue string) (*TokenInfo, error) {
	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(tm.jwtConfig.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.ErrInvalidToken
	}

	return &TokenInfo{
		UserId:         claims["userId"].(string),
		ExpirationTime: time.Unix(int64(claims["exp"].(float64)), 0).UTC(),
	}, nil
}
