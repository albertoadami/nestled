package auth

import (
	"testing"
	"time"

	"github.com/albertoadami/nestled/internal/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGeneratingToken(t *testing.T) {
	userId := uuid.New()
	secretKey := "test_secret_key"
	expireHours := 6

	token, err := GenerateToken(userId, secretKey, expireHours)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assert.NotNil(t, token.Value, "expected token value to be set, got nil")

}

func TestParsingInvalidToken(t *testing.T) {
	token := "invalid.token.value"

	_, err := ParseToken(token, "test_secret_key")

	assert.ErrorIs(t, err, errors.ErrInvalidToken)

}

func TestParsingValidTokenAndExtractInfo(t *testing.T) {
	userId := uuid.New()
	secretKey := "test_secret_key"

	token, err := GenerateToken(userId, secretKey, 6)
	assert.NoError(t, err)

	tokenInfo, err := ParseToken(token.Value, secretKey)
	assert.NoError(t, err)

	assert.Equal(t, userId.String(), tokenInfo.UserId)
	assert.Equal(t, token.ExpirationTime.Truncate(time.Second), tokenInfo.ExpirationTime.Truncate(time.Second))
}
