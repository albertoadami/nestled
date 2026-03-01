package routes

import (
	"github.com/albertoadami/nestled/internal/handlers"
	"github.com/gin-gonic/gin"
)

const ApiPrefix = "/api/v1"

func SetupRoutes(r *gin.Engine, userHandler *handlers.UserHandler, healthHandler *handlers.HealthHandler, authHandler *handlers.AuthHandler) {
	r.GET("/health", healthHandler.Health)

	apiGroup := r.Group(ApiPrefix)
	apiGroup.POST("/register", userHandler.RegisterUser)
	apiGroup.POST("/auth/token", authHandler.GenerateToken)
}
