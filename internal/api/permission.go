package api

import (
	"net/http"
	"strings"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func PermissionRoutesPrivate(router *gin.RouterGroup, config k4ever.Config) {
	permissions := router.Group("/permissions/")
	{
		getPermissions(permissions, config)
		getPermission(permissions, config)
		createPermission(permissions, config)
	}
}

func getPermissions(router *gin.RouterGroup, config k4ever.Config) {
	router.GET("", func(c *gin.Context) {
		var permissions []models.Permission
		if err := config.DB().Find(&permissions).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, permissions)
	})
}

func getPermission(router *gin.RouterGroup, config k4ever.Config) {
	router.GET(":id", func(c *gin.Context) {
		var permission models.Permission
		if err := config.DB().Find(&permission).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, permission)
	})
}

func createPermission(router *gin.RouterGroup, config k4ever.Config) {
	router.POST("", func(c *gin.Context) {
		var permission models.Permission
		if err := c.ShouldBindJSON(&permission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := config.DB().Create(&permission).Error; err != nil {
			if strings.HasPrefix(err.Error(), "UNIQUE constraint failed:") {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Permission already exists"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, permission)
	})
}
