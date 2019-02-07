package server

import (
	"fmt"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var authMiddleware *jwt.GinJWTMiddleware
var configForAuth k4ever.Config

func Start(config k4ever.Config) {
	configForAuth = config
	app := gin.Default()
	app.Use(cors.Default())

	authMiddleware = &jwt.GinJWTMiddleware{
		Realm:      "emtpy",          // TODO
		Key:        []byte("secret"), // TODO
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				fmt.Print(v.UserName)
				return jwt.MapClaims{
					"id":   v.ID,
					"name": v.UserName,
				}
			}
			return nil
		},
		//IdentityHandler: getIdentity,
		Authenticator: authenticate,
	}

	// Register all routes to the api as well as the frontend
	registerRoutes(app, config)

	// Run the webserver
	app.Run(fmt.Sprintf(":%d", config.HttpServerPort()))
}

type login struct {
	Username string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func getIdentity(claims jwt.MapClaims) interface{} {
	user := &models.User{}
	//uid, err := strconv.ParseUint(claims["id"].(string), 10, 64)
	err := configForAuth.DB().Where("user_name = ?", claims["name"].(string)).First(&user).Error
	if err != nil {
		return nil
	}
	//user.ID = uint(uid)
	return user
}

func authenticate(c *gin.Context) (interface{}, error) {
	var loginVars login
	var user models.User
	if err := c.ShouldBindJSON(&loginVars); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}
	if err := configForAuth.DB().Where("user_name = ?", loginVars.Username).First(&user).Error; err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginVars.Password)); err != nil {
		return nil, err
	}

	return &user, nil
}
