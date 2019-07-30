package api

import (
	"net/http"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func PurchaseRoutes(router *gin.RouterGroup, config k4ever.Config) {
	purchases := router.Group("/:name/purchases/")
	{
		getPurchaseHistory(purchases, config)
	}
}

// swagger:route GET /users/{name}/purchases/ users purchases getPurchaseHistory
//
// Get a list of all purchases
//
//		Produces:
//		- application/json
//
//		Security:
//        jwt:
//
//		Responses:
//		  default: GenericError
//		  200: HistoryArray
//		  400: GenericError
func getPurchaseHistory(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters getPurchaseHistory
	type getPurchaseHistoryParams struct {
		// in: path
		// required: true
		Name string `json:"name"`
	}
	router.GET("", func(c *gin.Context) {
		/*var user models.User
		var err error
		if user, err = k4ever.GetUser(c.Param("name"), config); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		*/
		var histories []models.History
		/*
			if err = config.DB().Preload("Items").Model(&user).Related(&histories).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
				return
			}
		*/
		c.JSON(http.StatusOK, histories)
	})
}
