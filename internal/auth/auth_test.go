package auth

import (
	"testing"
	"time"

	"github.com/albertoadami/nestled/internal/config"
	"github.com/albertoadami/nestled/internal/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var tokenManager = NewTokenManager(
	config.JWTConfig{
		Secret:     "test_secret_key",
		Expiration: 6,
	},
)

var expiredTokenManager = NewTokenManager(
	config.JWTConfig{
		Secret:     "test_secret_key",
		Expiration: -1, // token already expired
	},
)

func TestGeneratingToken(t *testing.T) {
	userId := uuid.New()

	token, err := tokenManager.GenerateToken(userId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assert.NotNil(t, token.Value, "expected token value to be set, got nil")

}

func TestParsingInvalidToken(t *testing.T) {
	token := "invalid.token.value"

	_, err := tokenManager.ParseToken(token)

	assert.ErrorIs(t, err, errors.ErrInvalidToken)

}

func TestParsingValidTokenAndExtractInfo(t *testing.T) {
	userId := uuid.New()

	token, err := tokenManager.GenerateToken(userId)
	assert.NoError(t, err)

	tokenInfo, err := tokenManager.ParseToken(token.Value)
	assert.NoError(t, err)

	assert.Equal(t, userId.String(), tokenInfo.UserId)
	assert.Equal(t, token.ExpirationTime.Truncate(time.Second), tokenInfo.ExpirationTime.Truncate(time.Second))
}
