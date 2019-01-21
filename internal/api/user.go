package api

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, config k4ever.Config) {
	users := router.Group("/users/")
	{
		getUsers(users, config)
		getUser(users, config)
		createUser(users, config)
		addPermissionToUser(users, config)
		PurchaseRoutes(users, config)
	}
}

func getUsers(router *gin.RouterGroup, config k4ever.Config) {
	router.GET("", func(c *gin.Context) {
		var users []models.User
		if err := config.DB().Find(&users).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})
}

func getUser(router *gin.RouterGroup, config k4ever.Config) {
	router.GET(":id", func(c *gin.Context) {
		var user models.User
		if err := config.DB().Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})
}

func createUser(router *gin.RouterGroup, config k4ever.Config) {
	router.POST("", func(c *gin.Context) {
		var user models.User
		var err error
		if err = c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error while hashing password"})
			return
		}
		user.Password = string(password)
		if err = config.DB().Create(&user).Error; err != nil {
			if strings.HasPrefix(err.Error(), "UNIQUE constraint failed:") {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
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
