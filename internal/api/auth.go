package api

import (
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// swagger:model
type login struct {
	Username string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var AuthMiddleware *jwt.GinJWTMiddleware
var configForAuth k4ever.Config

func CreateAuthMiddleware(config k4ever.Config) {
	configForAuth = config
	AuthMiddleware = &jwt.GinJWTMiddleware{
		Realm:      "emtpy",          // TODO
		Key:        []byte("secret"), // TODO
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"id":   v.Uid,
					"name": v.UserName,
				}
			}
			return nil
		},
		//IdentityHandler: getIdentity,
		Authenticator: authenticate,
	}
}

func getIdentity(claims jwt.MapClaims) interface{} {
	user := &models.User{}
	//uid, err := strconv.ParseUint(claims["id"].(string), 10, 64)
	/*err := configForAuth.DB().Where("user_name = ?", claims["name"].(string)).First(&user).Error
	if err != nil {
		return nil
	}*/
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
	/*
		conn, err := utils.LdapConnect(configForAuth)
		log.Debug("Testing")

		if err != nil {
			log.WithFields(log.Fields{"error": err.Error()}).Debug("Could not connect to ldap, querying database")
		}
		if conn != nil {
			defer conn.Close()
		}
		if err == nil {
			err = utils.LdapAuth(loginVars.Username, loginVars.Password, conn, configForAuth)
			if err == nil {
				user.UserName = loginVars.Username
				if err = configForAuth.DB().Where("user_name = ?", loginVars.Username).FirstOrCreate(&user).Error; err == nil {
					log.WithFields(log.Fields{"user": loginVars.Username}).Debug("Logged in via ldap")
					return &user, nil
				} else {
					log.WithFields(log.Fields{"error": err.Error()}).Error("Could not query/update database after authenticating against ldap")
					return nil, err
				}
			} else {
				log.WithFields(log.Fields{"error": err.Error()}).Debug("Failed to authenticate against ldap. Trying database entries...")
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
			return nil, jwt.ErrFailedAuthentication
		}
	*/

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
