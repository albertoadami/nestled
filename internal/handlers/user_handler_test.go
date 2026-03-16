package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"encoding/json"

	"github.com/albertoadami/nestled/internal/crypto"

	"github.com/albertoadami/nestled/internal/dto"
	"github.com/albertoadami/nestled/internal/errors"
	"github.com/albertoadami/nestled/internal/model"
	"github.com/albertoadami/nestled/internal/testhelpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockUserService struct {
	createUserFn      func(req *dto.CreateUserRequest) (uuid.UUID, error)
	getByIdFn         func(id uuid.UUID) (*model.User, error)
	returnUpdateError bool
	isTokenValid      bool
}

func (m *mockUserService) CreateUser(req *dto.CreateUserRequest) (uuid.UUID, error) {
	return m.createUserFn(req)
}

func (m *mockUserService) GetUserById(id uuid.UUID) (*model.User, error) {
	return m.getByIdFn(id)
}

func (m *mockUserService) ChangePassword(user *model.User, currentPassword string, newPassword string) error {
	if m.returnUpdateError {
		return errors.CredentialsInvalid
	} else {
		return nil
	}
}

func (m *mockUserService) ActivateUser(token string) error {
	if m.isTokenValid {
		return nil
	} else {
		return errors.ErrInvalidToken
	}
}

func setupUserRouter(mockService *mockUserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewUserHandler(mockService, zap.NewNop())
	router.POST("/api/v1/users/register", handler.RegisterUserHandler)
	return router
}

func setUpUserProfileRouter(mockService *mockUserService, userId uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewUserHandler(mockService, zap.NewNop())
	// apply mock authentication as middleware before the handler
	router.GET("/api/v1/users/me", testhelpers.MockAuthentication(userId), handler.GetCurrentUserHandler)
	router.PATCH("/api/v1/users/me/password", testhelpers.MockAuthentication(userId), handler.ChangePasswordHandler)
	router.POST("api/v1/users/activate", handler.ActivateUserHandler)
	return router
}

func createUserRequest() *http.Request {
	body := `{"username":"test","email":"test@github.com","password":"secret123", "first_name":"Test","last_name":"User"}`
	req, _ := http.NewRequest("POST", "/api/v1/users/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestRegisterUserSuccessfully(t *testing.T) {

	userId := uuid.New()

	mockService := &mockUserService{
		createUserFn: func(req *dto.CreateUserRequest) (uuid.UUID, error) {
			return userId, nil
		},
	}

	router := setupUserRouter(mockService)

	req := createUserRequest()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status code 201, got %d", w.Code)
	}

	location := w.Header().Get("Location")
	expectedLocation := "/api/v1/users/" + userId.String()
	assert.Equal(t, expectedLocation, location)

}

func TestRegisterUsernameAlreadyExists(t *testing.T) {

	mockService := &mockUserService{
		createUserFn: func(req *dto.CreateUserRequest) (uuid.UUID, error) {
			return uuid.Nil, errors.ErrUsernameAlreadyExists
		},
	}

	router := setupUserRouter(mockService)

	req := createUserRequest()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegisterEmailAlreadyExists(t *testing.T) {

	mockService := &mockUserService{
		createUserFn: func(req *dto.CreateUserRequest) (uuid.UUID, error) {
			return uuid.Nil, errors.ErrEmailAlreadyExists
		},
	}

	router := setupUserRouter(mockService)

	req := createUserRequest()
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUserProfileSuccessfully(t *testing.T) {

	userId := uuid.New()

	mockService := &mockUserService{
		getByIdFn: func(id uuid.UUID) (*model.User, error) {
			if id == userId {
				return &model.User{
					Id:        userId,
					Username:  "test",
					Email:     "test@test.it",
					FirstName: "Test",
					LastName:  "User",
				}, nil
			}
			return nil, nil
		},
	}

	router := setUpUserProfileRouter(mockService, userId)

	req, _ := http.NewRequest("GET", "/api/v1/users/me", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := dto.UserResponse{
		Id:        userId.String(),
		Username:  "test",
		Email:     "test@test.it",
		FirstName: "Test",
		LastName:  "User",
	}
	expectedJSON, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expectedJSON), w.Body.String())
}

func TestChangePasswordCorrectly(t *testing.T) {

	userId := uuid.New()
	passwordHash, _ := crypto.HashPassword("oldpassword")

	mockService := &mockUserService{
		getByIdFn: func(id uuid.UUID) (*model.User, error) {
			if id == userId {
				return &model.User{
					Id:           userId,
					Username:     "test",
					Email:        "test@test.it",
					FirstName:    "Test",
					LastName:     "User",
					PasswordHash: passwordHash,
				}, nil
			}
			return nil, nil
		},
	}
	router := setUpUserProfileRouter(mockService, userId)

	passwordRequest := &dto.ChangePasswordRequest{
		CurrentPassword: "oldpassword",
		NewPassword:     "newpassword123",
	}
	body, _ := json.Marshal(passwordRequest)

	req, _ := http.NewRequest("PATCH", "/api/v1/users/me/password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

}

func TestChangePasswordInvalidCurrentPassword(t *testing.T) {

	userId := uuid.New()
	mockService := &mockUserService{
		getByIdFn: func(id uuid.UUID) (*model.User, error) {
			if id == userId {
				return &model.User{
					Id:           userId,
					Username:     "test",
					Email:        "test@test.it",
					FirstName:    "Test",
					LastName:     "User",
					PasswordHash: "blablah",
				}, nil
			}
			return nil, nil
		},
		returnUpdateError: true,
	}
	router := setUpUserProfileRouter(mockService, userId)

	passwordRequest := &dto.ChangePasswordRequest{
		CurrentPassword: "wrong_password",
		NewPassword:     "newpassword123",
	}
	body, _ := json.Marshal(passwordRequest)

	req, _ := http.NewRequest("PATCH", "/api/v1/users/me/password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func TestActivateUserWithValidToken(t *testing.T) {
	userId := uuid.New()
	mockService := &mockUserService{isTokenValid: true}

	router := setUpUserProfileRouter(mockService, userId)

	activateUserRequest := &dto.ActivateUserToken{Token: "test-token"}
	body, _ := json.Marshal(activateUserRequest)

	req, _ := http.NewRequest("POST", "/api/v1/users/activate", bytes.NewReader(body))
	req.Header.Set("Content'Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

}

func TestActivateUserWithInvalidToken(t *testing.T) {
	userId := uuid.New()
	mockService := &mockUserService{isTokenValid: false}

	router := setUpUserProfileRouter(mockService, userId)

	activateUserRequest := &dto.ActivateUserToken{Token: "test-token"}
	body, _ := json.Marshal(activateUserRequest)

	req, _ := http.NewRequest("POST", "/api/v1/users/activate", bytes.NewReader(body))
	req.Header.Set("Content'Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

}
