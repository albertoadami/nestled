package main

import (
	"github.com/albertoadami/nestled/internal/config"
	"github.com/albertoadami/nestled/internal/database"
	"github.com/albertoadami/nestled/internal/handlers"
	"github.com/albertoadami/nestled/internal/repositories"
	"github.com/albertoadami/nestled/internal/routes"
	"github.com/albertoadami/nestled/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Starting Nestled API server...")
	r := gin.Default()

	configuration, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	database, err := database.Connect(&configuration.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()

	// repositories
	userRepository := repositories.NewUserRepository(database)

	// services
	userService := services.NewUserService(userRepository)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, logger)
	healthHandler := handlers.NewHealthHandler(database)

	routes.SetupRoutes(r, userHandler, healthHandler)

	// Start the server
	if err := r.Run(); err != nil {
		panic(err)
	}

}
