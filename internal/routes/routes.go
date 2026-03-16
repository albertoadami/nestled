package routes

import (
	"github.com/albertoadami/nestled/internal/auth"
	"github.com/albertoadami/nestled/internal/handlers"
	"github.com/albertoadami/nestled/internal/middleware"
	"github.com/gin-gonic/gin"
)

const ApiPrefix = "/api/v1"

func SetupRoutes(r *gin.Engine, userHandler *handlers.UserHandler, healthHandler *handlers.HealthHandler, authHandler *handlers.AuthHandler, tokenManager *auth.TokenManager) {

	r.GET("/health", healthHandler.HealthHandler)

	apiGroup := r.Group(ApiPrefix)
	protected := r.Group(ApiPrefix).Use(middleware.BearerAuthentication(tokenManager))

	apiGroup.POST("/register", userHandler.RegisterUserHandler)
	apiGroup.POST("/auth/token", authHandler.GenerateTokenHandler)
	apiGroup.POST("/users/activate", userHandler.ActivateUserHandler)
	protected.GET("/users/me", userHandler.GetCurrentUserHandler)
	protected.PATCH("/users/me/password", userHandler.ChangePasswordHandler)
}
