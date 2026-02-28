package routes

import (
	"github.com/albertoadami/nestled/internal/handlers"
	"github.com/gin-gonic/gin"
)

const ApiPrefix = "/api/v1"

func SetupRoutes(r *gin.Engine, userHandler *handlers.UserHandler, healthHandler *handlers.HealthHandler) {
	r.GET("/health", healthHandler.Health)

	apiGroup := r.Group(ApiPrefix)
	apiGroup.POST("/register", userHandler.RegisterUser)

}
