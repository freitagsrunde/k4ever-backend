package server

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/appleboy/gin-jwt"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var authMiddleware *jwt.GinJWTMiddleware
var configForAuth k4ever.Config

func Start(config k4ever.Config) {
	configForAuth = config
	app := gin.Default()
	//app.Use(cors.Default())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	app.Use(cors.New(corsConfig))

	authMiddleware = &jwt.GinJWTMiddleware{
		Realm:      "emtpy",          // TODO
		Key:        []byte("secret"), // TODO
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
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

// swagger:model
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

// swagger:route POST /login/ auth authenticateP
//
// Return a jwt token on login
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		Responses:
//		  default: GenericError
//		  200: TokenResponse
//		  401: GenericError
func authenticate(c *gin.Context) (interface{}, error) {
	var loginVars login
	var user models.User
	if err := c.ShouldBindJSON(&loginVars); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	// Check ldap
	conn, err := connect(configForAuth)
	defer conn.Close()

	if err != nil {
		log.Debug("Could not connect to ldap, querying database")
	} else {
		err = ldapAuth(loginVars.Username, loginVars.Password, conn, configForAuth)
		if err == nil {
			user.UserName = loginVars.Username
			if err = configForAuth.DB().Where("user_name = ?", loginVars.Username).FirstOrCreate(&user).Error; err == nil {
				log.WithFields(log.Fields{"user": loginVars.Username}).Debug("Created user from ldap")
				return &user, nil
			} else {
				log.Error("Could not insert user into database after authenticating against ldap")
				return nil, err
			}
		}
	}

	// Check local db
	if err := configForAuth.DB().Where("user_name = ?", loginVars.Username).First(&user).Error; err != nil {
		log.Debug("Login failed: user not found")
		time.Sleep(200 * time.Millisecond)
		return nil, jwt.ErrFailedAuthentication
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginVars.Password)); err != nil {
		log.Debug("Login failed: password was wrong")
		time.Sleep(200 * time.Millisecond)
		return nil, err
	}

	return &user, nil
}

// A token with an expiry date
//
// swagger:response
type TokenResponse struct {
	// in: body
	Token Token
}

// This is just for swagger

// The returned token
//
// swagger:model Token
type Token struct {
	Code   string `json:"code"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

// swagger:parameters authenticateP
type authenticateParams struct {
	// in: body
	Login login
}
