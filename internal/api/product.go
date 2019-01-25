package api

import (
	"net/http"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ProductRoutesPublic(router *gin.RouterGroup, config k4ever.Config) {
	products := router.Group("/products/")
	{
		getProducts(products, config)
		getProduct(products, config)
		getProductImage(products, config)
	}
}

func ProductRoutesPrivate(router *gin.RouterGroup, config k4ever.Config) {
	products := router.Group("/products/")
	{
		createProduct(products, config)
		buyProduct(products, config)
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
			return
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

func buyProduct(router *gin.RouterGroup, config k4ever.Config) {
	router.POST(":id/buy", func(c *gin.Context) {
		var product models.Product
		tx := config.DB().Begin()
		// Get Product
		if err := tx.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		purchase := models.Purchase{Amount: product.Price}
		item := models.Item{Amount: 1, Product: product, ProductID: product.ID}
		// Create PurchaseItem
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		purchase.Items = append(purchase.Items, item)
		// Create Purchase
		if err := tx.Create(&purchase).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Update Balance
		var user models.User
		session := sessions.Default(c)
		userID := session.Get("user")
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}
		user.Balance = user.Balance - product.Price
		user.Purchases = append(user.Purchases, purchase)
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tx.Commit()
		c.JSON(http.StatusOK, purchase)
	})
}
