package server

import (
	"github.com/freitagsrunde/k4ever-backend/internal/api"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func registerRoutes(app *gin.Engine, config k4ever.Config) {
	app.GET("/metrics", gin.WrapH(promhttp.Handler()))
	v1 := app.Group("/api/v1/")
	{
		api.ProductRoutesPublic(v1, config)
		v1.POST("/login/", authMiddleware.LoginHandler)
		api.VersionRoutesPublic(v1, config)
		api.SwaggerRoutesPublic(v1, config)
	}
	v1Private := app.Group("/api/v1/")
	{
		v1Private.Use(authMiddleware.MiddlewareFunc())

		api.PermissionRoutesPrivate(v1Private, config)
		api.ProductRoutesPrivate(v1Private, config)
		api.UserRoutesPrivate(v1Private, config)
		v1Private.Static("files", "./files")
	}
}
