package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albertoadami/nestled/internal/testhelpers"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestHealthHandlerSuccess(t *testing.T) {
	// Start PostgreSQL container and get connection
	db, terminate := testhelpers.SetupPostgres(t)
	defer terminate()

	// set up Gin router and handler
	gin.SetMode(gin.TestMode)
	router := gin.New()
	healthHandler := NewHealthHandler(db)

	router.GET("/health", healthHandler.Health)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestHealthHandlerDBDown(t *testing.T) {
	db, terminate := testhelpers.SetupPostgres(t)
	defer terminate()

	// Simulate DB down
	db.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	healthHandler := NewHealthHandler(db)
	router.GET("/health", healthHandler.Health)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
