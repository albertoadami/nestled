package routes

import (
	"github.com/albertoadami/nestled/internal/auth"
	"github.com/albertoadami/nestled/internal/handlers"
	"github.com/albertoadami/nestled/internal/middleware"
	"github.com/gin-gonic/gin"
)

const ApiPrefix = "/api/v1"

func SetupRoutes(r *gin.Engine, userHandler *handlers.UserHandler, healthHandler *handlers.HealthHandler, authHandler *handlers.AuthHandler, tokenManager *auth.TokenManager) {

	r.GET("/health", healthHandler.Health)

	apiGroup := r.Group(ApiPrefix)
	protected := r.Group(ApiPrefix).Use(middleware.BearerAuthentication(tokenManager))

	apiGroup.POST("/register", userHandler.RegisterUser)
	apiGroup.POST("/auth/token", authHandler.GenerateToken)
	protected.GET("/users/me", userHandler.GetCurrentUser)
	protected.PATCH("/users/me/password", userHandler.ChangePassword)
}
