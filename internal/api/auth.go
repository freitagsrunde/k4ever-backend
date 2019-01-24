package api

import (
	"net/http"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AuthRoutesPublic(router *gin.RouterGroup, config k4ever.Config) {
	auth := router.Group("")
	{
		login(auth, config)
		logout(auth, config)
	}
}

func login(router *gin.RouterGroup, config k4ever.Config) {
	router.POST("/login2/", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		password := user.Password
		if err := config.DB().Where("user_name = ?", user.UserName).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated"})
	})
}

func logout(router *gin.RouterGroup, config k4ever.Config) {
	router.POST("/logout/", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}
		session.Delete("user")
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	})
}
