package server

import (
	"github.com/freitagsrunde/k4ever-backend/internal/api"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/gin-gonic/gin"
)

func registerRoutes(app *gin.Engine, config k4ever.Config) {
	v1 := app.Group("/api/v1/")
	{
		api.ProductRoutes(v1, config)
	}
}
