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
		getProduct(products, config)
		createProduct(products, config)
		getProductImage(products, config)
	}
}

func getProducts(router *gin.RouterGroup, config k4ever.Config) {
	router.GET("", func(c *gin.Context) {
		var products []models.Product
		if err := config.DB().Find(&products).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})
}

func getProduct(router *gin.RouterGroup, config k4ever.Config) {
	router.GET(":id/", func(c *gin.Context) {
		var product models.Product
		if err := config.DB().Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, product)
	})
}

func createProduct(router *gin.RouterGroup, config k4ever.Config) {
	router.POST("", func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := config.DB().Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, product)
	})
}

func getProductImage(router *gin.RouterGroup, config k4ever.Config) {
	router.GET(":id/image/", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"Hello": "World"})
	})
}
