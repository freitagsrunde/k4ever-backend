package api

import (
	"net/http"
	"strconv"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func PurchaseRoutes(router *gin.RouterGroup, config k4ever.Config) {
	purchases := router.Group("/:id/purchases/")
	{
		getPurchaseHistory(purchases, config)
	}
}

func getPurchaseHistory(router *gin.RouterGroup, config k4ever.Config) {
	router.GET("", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		user := models.User{}
		user.ID = uint(id)
		var purchases []models.Purchase
		if err = config.DB().Preload("Items").Model(&user).Related(&purchases).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		c.JSON(http.StatusOK, purchases)
	})
}