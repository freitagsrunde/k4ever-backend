package api

import (
	"net/http"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/gin-gonic/gin"
)

func VersionRoutesPublic(router *gin.RouterGroup, config k4ever.Config) {
	versionGroup := router.Group("/version/")
	{
		GetVersion(versionGroup, config)
	}
}

type VersionInformation struct {
	Version   string
	GitBranch string
	GitCommit string
	BuildTime string
}

func GetVersion(router *gin.RouterGroup, config k4ever.Config) {
	version := VersionInformation{Version: config.Version(), GitBranch: config.GitBranch(), GitCommit: config.GitCommit(), BuildTime: config.BuildTime()}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, version)
	})
}
