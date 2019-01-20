package api

import (
	"net/http"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AuthRoutes(router *gin.RouterGroup, config k4ever.Config) {
	auth := router.Group("")
	{
		login(auth, config)
	}
}

func login(router *gin.RouterGroup, config k4ever.Config) {
	router.POST("/login/", func(c *gin.Context) {
		session := sessions.Default(c)
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
		session.Set("user", user.ID)
		if err := session.Save(); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated"})
	})
}
