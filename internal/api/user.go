package api

import (
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
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
		transferToUser(users, config)
	}
}

// swagger:route GET /users/ users getUsers
//
// Lists all users
//
// This will show all available users by default
//
// 		Produces:
//      - applications/json
//
//		Security:
//		  jwt:
//
//		Responses:
//		  default: GenericError
// 	 	  200: UsersResponse
//		  404: GenericError
func getUsers(router *gin.RouterGroup, config k4ever.Config) {
	// A UsersResponse returns a list of users
	//
	// swagger:response
	type UsersResponse struct {
		// An array of users
		//
		// in: body
		Users []models.User
	}
	router.GET("", func(c *gin.Context) {
		var err error
		params := models.DefaultParams{}
		params.SortBy = c.DefaultQuery("sort_by", "user_name")
		params.Order = c.DefaultQuery("order", "asc")
		offset := c.Query("offset")
		if offset != "" {
			params.Offset, err = strconv.Atoi(offset)
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "offset is not a number"})
			return
		}
		limit := c.Query("limit")
		if limit != "" {
			params.Limit, err = strconv.Atoi(limit)
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "limit is not a number"})
			return
		}

		users, err := k4ever.GetUsers(params, config)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}
		c.JSON(http.StatusOK, users)
	})
}

// swagger:route GET /users/{name}/ user getUser
//
// Get detailed information of a user
//
// This will show detailed information for a specific user
//
//		Produces:
//		- application/json
//
//		Security:
//        jwt:
//
//		Responses:
//		  default: GenericError
//	  	  200: User
//		  404: GenericError
func getUser(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters getUser
	type getUserParams struct {
		// in:path
		// required: true
		Name string `json:"name"`
	}
	router.GET(":name/", func(c *gin.Context) {
		var user models.User
		var err error
		if user, err = k4ever.GetUser(c.Param("name"), config); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	})
}

// Input params for creating a user
//
// swagger:model
type newUser struct {
	UserName    string `json:"name"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

// swagger:route POST /users/ users createUser
//
// Create a new user
//
// 		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		Security:
//        jwt:
//
//		Responses:
//		  default: GenericError
//        201: User
//		  400: GenericError
//	      500: GenericError
func createUser(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters createUser
	type CreateUserParams struct {
		// in: body
		// required: true
		NewUser newUser
	}
	router.POST("", func(c *gin.Context) {
		var bind newUser
		var user models.User
		if err := c.ShouldBindJSON(&bind); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.UserName = bind.UserName
		user.Password = bind.Password
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

// swagger:route PUT /users/{name}/permissions/ user permission addPermissionToUser
//
// Add permission to user
//
// Links an existing permission to a user
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		Security:
//		  jwt:
//
//		Responses:
//		  default: GenericError
//        203: User
//		  400: GenericError
//		  404: GenericError
func addPermissionToUser(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters addPermissionToUser
	type AddPermissionParam struct {
		// in: path
		// required: true
		Name string `json:"name"`
		// in: body
		// required: true
		Permission models.Permission
	}
	router.PUT(":name/permissions/", func(c *gin.Context) {
		var user models.User
		var err error
		var permission models.Permission
		if err = c.ShouldBindJSON(&permission); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if user, err = k4ever.GetUser(c.Param("name"), config); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err = config.DB().Where("name = ?", permission.Name).First(&permission).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		user.Permissions = append(user.Permissions, permission)
		config.DB().Save(&user)

		c.JSON(http.StatusAccepted, user)
	})
}

// swagger:model
type Balance struct {
	Amount float64
}

// swagger:route PUT /users/{name}/balance/ user balance addBalance
//
// Add balance
//
// Add the given balance to the logged in user
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		Security:
//        jwt:
//
//		Responses:
//		  default: GenericError
//		  200: User
//		  400: GenericError
//        404: GenericError
//        500: GenericError
func addBalance(router *gin.RouterGroup, config k4ever.Config) {

	// swagger:parameters addBalance
	type AddBalanceParams struct {
		// in: path
		// required: true
		Name string `json:"name"`

		// in: body
		// required: true
		Balance Balance
	}
	router.PUT(":name/balance/", func(c *gin.Context) {
		var err error
		var balance Balance
		if err := c.ShouldBindJSON(&balance); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims := jwt.ExtractClaims(c)
		username := claims["name"]

		if username == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		if username != c.Param("name") {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "You can only update you own account"})
			return
		}

		// Add balance to specified account
		var balanceHistory models.History
		balanceHistory, err = k4ever.AddBalance(c.Param("name"), balance.Amount, config)
		if err != nil {
			if err.Error() == "record not found" {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, balanceHistory)
	})
}

// swagger:route PUT /users/{name}/transfer/ user balance transferToUser
//
// Transfer money from the current user to the user in the path
//
// Transfers the exact given amount from the body from the current user to the user in the path. The transfer fails if the amount is 0 or lower.
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		Security:
//        jwt:
//
//		Responses:
//		  default: GenericError
//		  200: History
//		  400: GenericError
//        404: GenericError
//        500: GenericError
func transferToUser(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters transferToUser
	type TransferToUserParams struct {
		// in: path
		// required: true
		Name string `json:"name"`

		// in: body
		// required: true
		balance Balance
	}
	router.PUT(":name/transfer/", func(c *gin.Context) {
		// Get current user
		claims := jwt.ExtractClaims(c)
		username := claims["name"]
		if username == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		// Parse amount from body
		var balance Balance
		if err := c.ShouldBindJSON(&balance); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update balances
		transfer, err := k4ever.TransferToUser(username.(string), c.Param("name"), balance.Amount, config)
		if err != nil {
			if err.Error() == "record not found" {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, transfer)
	})
}
