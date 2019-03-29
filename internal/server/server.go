package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/appleboy/gin-jwt"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	ldap "gopkg.in/ldap.v2"
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
		fmt.Println(err.Error())
		fmt.Println("Could not connect to ldap")
	} else {
		err = ldapAuth(loginVars.Username, loginVars.Password, conn, configForAuth)
		if err == nil {
			user.UserName = loginVars.Username
			if err = configForAuth.DB().Where("user_name = ?", loginVars.Username).FirstOrCreate(&user).Error; err == nil {
				return &user, nil
			} else {
				return nil, err
			}
		}
	}

	// Check local db
	if err := configForAuth.DB().Where("user_name = ?", loginVars.Username).First(&user).Error; err != nil {
		time.Sleep(200 * time.Millisecond)
		return nil, jwt.ErrFailedAuthentication
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginVars.Password)); err != nil {
		time.Sleep(200 * time.Millisecond)
		return nil, err
	}

	return &user, nil
}

// Connect to ldap and return connection object
func connect(config k4ever.Config) (*ldap.Conn, error) {
	conn, err := ldap.Dial("tcp", config.LdapHost())

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to ldap server: %s", config.LdapHost())
	}

	if err := conn.Bind(config.LdapBind(), config.LdapPassword()); err != nil {
		return nil, fmt.Errorf("Failed to bind to ldap server: %s", config.LdapBind())
	}
	return conn, nil
}

// try to authenticate user against ldap
func ldapAuth(user string, password string, conn *ldap.Conn, config k4ever.Config) error {
	result, err := conn.Search(ldap.NewSearchRequest(
		config.LdapBaseDN(),
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter(user, config),
		[]string{"dn"},
		nil,
	))

	if err != nil {
		return fmt.Errorf("Failed to find user: %s", user)
	}

	if len(result.Entries) < 1 {
		return fmt.Errorf("User does not exist: %s", user)
	}

	if len(result.Entries) > 1 {
		return fmt.Errorf("Too many entries returned")
	}

	if err := conn.Bind(result.Entries[0].DN, password); err != nil {
		fmt.Errorf("Failed to auth. %s", err)
	}
	return nil
}

func filter(needle string, config k4ever.Config) string {
	res := strings.Replace(
		config.LdapFilterDN(),
		"{username}",
		needle,
		-1,
	)
	return res
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
