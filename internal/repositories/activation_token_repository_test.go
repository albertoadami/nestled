package repositories

import (
	"os"
	"testing"
	"time"

	"github.com/albertoadami/nestled/internal/errors"
	"github.com/albertoadami/nestled/internal/model"
	"github.com/albertoadami/nestled/internal/testhelpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var _activationTokenRepository ActivationTokenRepository
var _userRepository UserRepository

func TestMain(m *testing.M) {
	db, terminate := testhelpers.SetupPostgres(nil)
	defer terminate()
	_activationTokenRepository = NewActivationTokenRepository(db)
	_userRepository = NewUserRepository(db)

	os.Exit(m.Run())
}

func createActivationToken(userId uuid.UUID) *model.ActivationToken {
	return &model.ActivationToken{Id: uuid.New(), UserId: userId, Token: "test-token", ExpiresAt: time.Now()}
}

func TestCreateSuccessfully(t *testing.T) {
	user := createTestUser()

	userId, err := _userRepository.Create(user)
	assert.NoError(t, err)

	activationToken := createActivationToken(userId)

	result, err := _activationTokenRepository.Create(activationToken)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, activationToken.Id, result)
}

func TestGetByTokenReturnsNil(t *testing.T) {
	result, err := _activationTokenRepository.GetByToken("not_existing")

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestGetByTokenReturnsResult(t *testing.T) {
	user := createTestUser()
	userId, err := _userRepository.Create(user)
	assert.NoError(t, err)

	activationToken := createActivationToken(userId)

	result, err := _activationTokenRepository.Create(activationToken)
	assert.NoError(t, err)

	getResult, err := _activationTokenRepository.GetByToken(activationToken.Token)

	assert.NoError(t, err)
	assert.Equal(t, getResult.Id, result)

}

func TestDeleteNoRows(t *testing.T) {
	err := _activationTokenRepository.DeleteById(uuid.New())

	assert.Equal(t, err, errors.ErrNoRowsAffected)

}

func TestDeleteSuccessfully(t *testing.T) {
	user := createTestUser()
	userId, err := _userRepository.Create(user)
	assert.NoError(t, err)

	activationToken := createActivationToken(userId)

	result, err := _activationTokenRepository.Create(activationToken)
	assert.NoError(t, err)

	err = _activationTokenRepository.DeleteById(result)
	assert.NoError(t, err)

}
