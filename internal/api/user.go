package api

import (
	"net/http"
	"strings"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func UserRoutesPrivate(router *gin.RouterGroup, config k4ever.Config) {
	users := router.Group("/users/")
	{
		getUsers(users, config)
		getUser(users, config)
		createUser(users, config)
		addPermissionToUser(users, config)
		PurchaseRoutes(users, config)
		addBalance(users, config)
	}
}

func getUsers(router *gin.RouterGroup, config k4ever.Config) {
	router.GET("", func(c *gin.Context) {
		var users []models.User
		if err := config.DB().Find(&users); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "No user was"})
			return
		}
		c.JSON(http.StatusOK, users)
	})
}

func getUser(router *gin.RouterGroup, config k4ever.Config) {
	router.GET(":id", func(c *gin.Context) {
		var user models.User
		var err error
		if user, err = k4ever.GetUser(c.Param("name"), config); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	})
}

func createUser(router *gin.RouterGroup, config k4ever.Config) {
	type newUser struct {
		UserName    string `json:"name"`
		Password    string `json:"password""`
		DisplayName string `json:"display_name"`
	}
	router.POST("", func(c *gin.Context) {
		var bind newUser
		var user models.User
		if err := c.ShouldBindJSON(&bind); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user.UserName = bind.UserName
		user.DisplayName = bind.DisplayName
		if err := k4ever.CreateUser(&user, config); err != nil {
			if strings.HasPrefix(err.Error(), "Username") {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	})
}

func addPermissionToUser(router *gin.RouterGroup, config k4ever.Config) {
	router.PUT(":id/permissions/", func(c *gin.Context) {
		var user models.User
		var permission models.Permission
		if err := c.ShouldBindJSON(&permission); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := config.DB().Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err := config.DB().Where("name = ?", permission.Name).First(&permission).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		user.Permissions = append(user.Permissions, permission)
		config.DB().Save(&user)

		c.JSON(http.StatusAccepted, user)
	})
}

func addBalance(router *gin.RouterGroup, config k4ever.Config) {
	type Balance struct {
		Amount float64
	}
	router.PUT(":id/balance/", func(c *gin.Context) {
		var user models.User
		var balance Balance
		if err := c.ShouldBindJSON(&balance); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tx := config.DB().Begin()
		if err := tx.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "No such user id"})
			return
		}
		user.Balance = user.Balance + balance.Amount
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tx.Commit()
		c.JSON(http.StatusOK, user)
	})
}
