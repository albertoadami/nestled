package auth

import (
	"testing"

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
