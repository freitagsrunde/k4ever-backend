package server

import (
	"fmt"
	"net/http"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Start(config k4ever.Config) {
	app := gin.Default()
	store := sessions.NewCookieStore([]byte("secret"))
	app.Use(sessions.Sessions("auth", store))

	// Register all routes to the api as well as the frontend
	registerRoutes(app, config)

	// Run the webserver
	app.Run(fmt.Sprintf(":%d", config.HttpServerPort()))
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Next()
	}
}
