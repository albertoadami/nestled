package routes

import (
	"github.com/albertoadami/nestled/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, healthHandler *handlers.HealthHandler) {
	r.GET("/health", healthHandler.Health)
}
