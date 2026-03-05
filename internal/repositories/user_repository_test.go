package repositories

import (
	"errors"
	"testing"

	apperrors "github.com/albertoadami/nestled/internal/errors"
	"github.com/albertoadami/nestled/internal/model"
	"github.com/albertoadami/nestled/internal/testhelpers"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func createTestUser() *model.User {
	return &model.User{
		Id:           uuid.New(),
		FirstName:    "John",
		LastName:     "Doe",
		Username:     "johndoe",
		Email:        "test@test.it",
		PasswordHash: "hashedpassword",
		Status:       model.UserStatusPending,
	}
}

func truncateUsers(t *testing.T, testDB *sqlx.DB) {
	_, err := testDB.Exec("TRUNCATE TABLE users CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate users: %v", err)
	}
}

func TestCreateUserSucessfully(t *testing.T) {

	db, terminate := testhelpers.SetupPostgres(t)
	defer terminate()
	truncateUsers(t, db)

	userRepo := NewUserRepository(db)

	user := createTestUser()

	id, err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if id != user.Id {
		t.Fatalf("expected id %v, got %v", user.Id, id)
	}

}

func TestCreateUserFailedDueToDuplicateUsername(t *testing.T) {

	db, terminate := testhelpers.SetupPostgres(t)
	defer terminate()

	truncateUsers(t, db)

	userRepo := NewUserRepository(db)

	user := createTestUser()
	_, err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	user.Email = "test-duplicated@test.it"
	user.Id = uuid.New()
	_, err = userRepo.Create(user)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, apperrors.ErrUsernameAlreadyExists) {
		t.Fatalf("expected error %v, got %v", apperrors.ErrUsernameAlreadyExists, err)
	}
}

func TestCreateUserFailedDueToDuplicateEmail(t *testing.T) {

	db, terminate := testhelpers.SetupPostgres(t)
	defer terminate()

	truncateUsers(t, db)

	userRepo := NewUserRepository(db)

	user := createTestUser()
	_, err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	user.Username = "johndoe-duplicated"
	user.Id = uuid.New()
	_, err = userRepo.Create(user)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, apperrors.ErrEmailAlreadyExists) {
		t.Fatalf("expected error %v, got %v", apperrors.ErrEmailAlreadyExists, err)
	}
}

func TestGetUserByUsernameSucessfully(t *testing.T) {

	db, terminate := testhelpers.SetupPostgres(t)
	defer terminate()
	truncateUsers(t, db)

	userRepo := NewUserRepository(db)
	user := createTestUser()
	_, err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedUser, err := userRepo.GetUserByUsername(user.Username)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrievedUser.Username != user.Username {
		t.Fatalf("expected username %v, got %v", user.Username, retrievedUser.Username)
	}

}

func TestGetUserByUsernameFailedDueToNonExistingUser(t *testing.T) {

	db, terminate := testhelpers.SetupPostgres(t)
	defer terminate()
	truncateUsers(t, db)

	userRepo := NewUserRepository(db)

	result, err := userRepo.GetUserByUsername("non-existing-username")
	assert.Nil(t, err, "expected err to be nil")
	assert.Nil(t, result, "expected result to be nil")
}

func TestGetUserByIdSucessfully(t *testing.T) {

	db, terminate := testhelpers.SetupPostgres(t)
	defer terminate()
	truncateUsers(t, db)

	userRepo := NewUserRepository(db)
	user := createTestUser()
	_, err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedUser, err := userRepo.GetUserById(user.Id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrievedUser.Username != user.Username {
		t.Fatalf("expected username %v, got %v", user.Username, retrievedUser.Username)
	}

}

func TestGetUserByIdFailedDueToNonExistingUser(t *testing.T) {

	db, terminate := testhelpers.SetupPostgres(t)
	defer terminate()
	truncateUsers(t, db)

	userRepo := NewUserRepository(db)

	result, err := userRepo.GetUserById(uuid.New())
	assert.Nil(t, err, "expected err to be nil")
	assert.Nil(t, result, "expected result to be nil")
}

func TestUpdateUserSuccessfully(t *testing.T) {

	db, terminate := testhelpers.SetupPostgres(t)
	defer terminate()
	truncateUsers(t, db)

	userRepo := NewUserRepository(db)
	user := createTestUser()
	_, err := userRepo.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	user.FirstName = "UpdatedFirstName"
	err = userRepo.Update(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedUser, err := userRepo.GetUserById(user.Id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if retrievedUser.FirstName != "UpdatedFirstName" {
		t.Fatalf("expected first name %v, got %v", "UpdatedFirstName", retrievedUser.FirstName)
	}

}
