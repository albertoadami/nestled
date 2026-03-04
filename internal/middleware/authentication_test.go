package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albertoadami/nestled/internal/auth"
	"github.com/albertoadami/nestled/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var tokenManager = auth.NewTokenManager(config.JWTConfig{Secret: "secret", Expiration: 1})
var expiredTokenManager = auth.NewTokenManager(config.JWTConfig{Secret: "secret", Expiration: -1})

func setupMiddlewareRouter(tokenManager *auth.TokenManager) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(BearerAuthentication(tokenManager))
	router.GET("/test", func(c *gin.Context) {
		userId, _ := c.Get("userId")
		c.JSON(http.StatusOK, gin.H{"userId": userId})
	})
	return router
}

func TestBearerAuthentication_MissingHeader(t *testing.T) {
	router := setupMiddlewareRouter(tokenManager)

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBearerAuthentication_InvalidToken(t *testing.T) {
	router := setupMiddlewareRouter(tokenManager)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestBearerAuthentication_ValidToken(t *testing.T) {
	userId := uuid.New()

	token, _ := tokenManager.GenerateToken(userId)

	router := setupMiddlewareRouter(tokenManager)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token.Value)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBearerAuthentication_ExpiredToken(t *testing.T) {
	userId := uuid.New()

	// generate token already expired (-1 hour)
	token, _ := expiredTokenManager.GenerateToken(userId)

	router := setupMiddlewareRouter(expiredTokenManager)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token.Value)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
