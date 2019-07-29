package server

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-contrib/cors"

	"github.com/freitagsrunde/k4ever-backend/internal/api"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/gin-gonic/gin"
)

func Start(config k4ever.Config) {
	app := gin.Default()
	//app.Use(cors.Default())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	app.Use(cors.New(corsConfig))
	app.HandleMethodNotAllowed = true

	api.CreateAuthMiddleware(config)

	// Register all routes to the api as well as the frontend
	registerRoutes(app, config)

	banner, err := ioutil.ReadFile("./assets/banner")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(string(banner))

	// Run the webserver
	app.Run(fmt.Sprintf(":%d", config.HttpServerPort()))
}
