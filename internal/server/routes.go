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
		api.AuthRoutes(v1, config)
	}
	v1Private := app.Group("/api/v1/")
	{
		v1Private.Use(AuthRequired())

		api.PermissionRoutes(v1Private, config)
		api.UserRoutes(v1Private, config)
		api.ProductRoutesPrivate(v1Private, config)
	}
}
