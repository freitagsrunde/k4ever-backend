package api

import (
	"net/http"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.RouterGroup, config k4ever.Config) {
	products := router.Group("/products/")
	{
		getProducts(products, config)
	}
}

func getProducts(router *gin.RouterGroup, config k4ever.Config) {
	router.GET("", func(c *gin.Context) {
		var products []models.Product
		if err := config.DB().Find(&products).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})
}
