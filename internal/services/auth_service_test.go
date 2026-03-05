package services

import (
	"testing"

	"github.com/albertoadami/nestled/internal/auth"
	"github.com/albertoadami/nestled/internal/config"
	"github.com/albertoadami/nestled/internal/crypto"
	"github.com/albertoadami/nestled/internal/errors"
	"github.com/albertoadami/nestled/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// mockUserRepo allows us to simulate the repository behavior.
type mockUserRepo struct {
	getFn     func(username string) (*model.User, error)
	getFnById func(id uuid.UUID) (*model.User, error)
}

func (m *mockUserRepo) Create(user *model.User) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (m *mockUserRepo) GetUserById(id uuid.UUID) (*model.User, error) {
	return m.getFnById(id)
}

func (m *mockUserRepo) GetUserByUsername(username string) (*model.User, error) {
	return m.getFn(username)
}

func (m *mockUserRepo) Update(user *model.User) error {
	return nil
}

var tokenManager = auth.NewTokenManager(config.JWTConfig{Secret: "secret", Expiration: 1})

func TestGenerateToken_UserNotFound(t *testing.T) {
	mockRepo := &mockUserRepo{
		getFn: func(username string) (*model.User, error) {
			return nil, nil
		},
		getFnById: func(userId uuid.UUID) (*model.User, error) {
			return nil, nil
		},
	}
	service := NewAuthService(mockRepo, tokenManager)

	token, err := service.GenerateToken("noexist", "pwd")
	assert.Nil(t, token)
	assert.ErrorIs(t, err, errors.CredentialsInvalid)
}

func TestGenerateToken_InvalidPassword(t *testing.T) {
	// create a user with a known hashed password
	hash, _ := crypto.HashPassword("correct")
	user := &model.User{Id: uuid.New(), PasswordHash: hash}

	mockRepo := &mockUserRepo{
		getFn: func(username string) (*model.User, error) {
			return user, nil
		},
	}
	service := NewAuthService(mockRepo, tokenManager)

	token, err := service.GenerateToken("someuser", "wrong")
	assert.Nil(t, token)
	assert.ErrorIs(t, err, errors.CredentialsInvalid)
}

func TestGenerateToken_Success(t *testing.T) {
	hash, _ := crypto.HashPassword("correct")
	user := &model.User{Id: uuid.New(), PasswordHash: hash}

	mockRepo := &mockUserRepo{
		getFn: func(username string) (*model.User, error) {
			return user, nil
		},
	}
	service := NewAuthService(mockRepo, tokenManager)

	token, err := service.GenerateToken("someuser", "correct")
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Value)
}

func TestUpdatePasswordInvalidCurrentPassword(t *testing.T) {
	hash, _ := crypto.HashPassword("current")
	user := &model.User{Id: uuid.New(), PasswordHash: hash}

	mockRepo := &mockUserRepo{
		getFnById: func(userId uuid.UUID) (*model.User, error) {
			return user, nil
		},
	}
	service := NewUserService(mockRepo)

	err := service.ChangePassword(user, "wrong", "newpass")
	assert.ErrorIs(t, err, errors.ErrInvalidPassword)
}

func TestUpdatePasswordSuccess(t *testing.T) {
	hash, _ := crypto.HashPassword("current")
	user := &model.User{Id: uuid.New(), PasswordHash: hash}

	mockRepo := &mockUserRepo{
		getFnById: func(userId uuid.UUID) (*model.User, error) {
			return user, nil
		},
	}
	service := NewUserService(mockRepo)

	err := service.ChangePassword(user, "current", "newpass")
	assert.NoError(t, err)

	// Verify that the password hash has been updated
	assert.NotEqual(t, hash, user.PasswordHash)
}
