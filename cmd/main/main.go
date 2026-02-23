package main

import (
	"github.com/albertoadami/nestled/internal/handlers"
	"github.com/albertoadami/nestled/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()

	// Setup routes
	routes.SetupRoutes(r, healthHandler)

	// Start the server
	if err := r.Run(); err != nil {
		panic(err)
	}

}
