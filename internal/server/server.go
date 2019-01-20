package server

import (
	"fmt"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/gin-gonic/gin"
)

func Start(config k4ever.Config) {
	app := gin.Default()

	// Register all routes to the api as well as the frontend
	registerRoutes(app, config)

	// Run the webserver
	app.Run(fmt.Sprintf(":%d", config.HttpServerPort()))
}
