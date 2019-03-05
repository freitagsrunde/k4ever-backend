package server

import (
	"github.com/freitagsrunde/k4ever-backend/internal/api"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/gin-gonic/gin"
)

func registerRoutes(app *gin.Engine, config k4ever.Config) {
	v1 := app.Group("/api/v1/")
	{
		api.ProductRoutesPublic(v1, config)
		v1.POST("/login/", authMiddleware.LoginHandler)
		api.UserRoutesPrivate(v1, config)
		api.VersionRoutesPublic(v1, config)
	}
	v1Private := app.Group("/api/v1/")
	{
		v1Private.Use(authMiddleware.MiddlewareFunc())

		api.PermissionRoutesPrivate(v1Private, config)
		api.ProductRoutesPrivate(v1Private, config)
	}
}
