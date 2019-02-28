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

// swagger:model
type VersionInformation struct {
	Version   string `json:"version"`
	GitBranch string `json:"branch"`
	GitCommit string `json:"commit"`
	BuildTime string `json:"build_time"`
}

// swagger:route GET /version/ getVersion
//
// Fetch version information
//
//		Produces:
//		- application/json
//
//		Responses:
//		  default: VersionInformation
//		  200: VersionInformation
func GetVersion(router *gin.RouterGroup, config k4ever.Config) {
	version := VersionInformation{Version: config.Version(), GitBranch: config.GitBranch(), GitCommit: config.GitCommit(), BuildTime: config.BuildTime()}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, version)
	})
}
