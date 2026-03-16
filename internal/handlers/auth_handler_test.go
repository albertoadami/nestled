package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/albertoadami/nestled/internal/auth"
	appErrorrs "github.com/albertoadami/nestled/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockAuthService struct {
	generateTokenFn func(username string, password string) (*auth.Token, error)
}

func (m *mockAuthService) GenerateToken(username string, password string) (*auth.Token, error) {
	return m.generateTokenFn(username, password)
}

func setupUserAuthRouter(mockService *mockAuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewAuthHandler(mockService)
	router.POST("/api/v1/auth/token", handler.GenerateTokenHandler)
	return router
}

func TestGenerateTokenSueccessfully(t *testing.T) {

	mockService := &mockAuthService{
		generateTokenFn: func(username string, password string) (*auth.Token, error) {
			return &auth.Token{Value: "mocked.jwt.token"}, nil
		},
	}
	router := setupUserAuthRouter(mockService)

	body := `{"username":"test","password":"secret123"}`
	req, _ := http.NewRequest("POST", "/api/v1/auth/token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mocked.jwt.token")
}

func TestGenerateTokenInvalidCredentials(t *testing.T) {
	mockService := &mockAuthService{
		generateTokenFn: func(username string, password string) (*auth.Token, error) {
			return nil, appErrorrs.CredentialsInvalid
		},
	}
	router := setupUserAuthRouter(mockService)

	body := `{"username":"test","password":"wrong"}`
	req, _ := http.NewRequest("POST", "/api/v1/auth/token", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid username or password")
}
