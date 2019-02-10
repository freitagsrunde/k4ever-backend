package api

import (
	"net/http"
	"strings"

	"github.com/freitagsrunde/k4ever-backend/internal/api/response"
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

// swagger:route GET /users/ users getUsers
//
// Lists all users
//
// This will show all available users by default
//
// 		Produces:
//      - applications/json
//
//		Responses:
//		  default: GenericError
// 	 	  200: UsersResponse
//		  404: GenericError
func getUsers(router *gin.RouterGroup, config k4ever.Config) {
	router.GET("", func(c *gin.Context) {
		var users []models.User
		if err := config.DB().Find(&users); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "No user was found"})
			return
		}
		c.JSON(http.StatusOK, response.UsersResponse{Users: users})
	})
}

// swagger:route GET /users/{id}/ user getUser
//
// Get detailed information of a user
//
// This will show detailed information for a specific user
//
//		Produces:
//		- application/json
//
//		Responses:
//		  default: GenericError
//	  	  200: User
//		  404: GenericError
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

// Input params for creating a user
//
// swagger:model
type newUser struct {
	UserName    string `json:"name"`
	Password    string `json:"password""`
	DisplayName string `json:"display_name"`
}

// swagger:route POST /users/ product createUser
//
// Create a new user
//
// 		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		Responses:
//		  default: GenericError
//        201: CreateUserParams
//		  400: GenericError
//	      500: GenericError
func createUser(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters createUser
	type CreateUserParams struct {
		// in: path
		// required: true
		Id string `json:"id"`
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

// swagger:route PUT /users/{id}/permissions/ user permission addPermissionToUser
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
		Id string `json:"id"`
		// in: body
		// required: true
		Permission models.Permission
	}
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

// swagger:model
type Balance struct {
	Amount float64
}

// swagger:route PUT /users/{id}/balance/ user balance addBalance
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
		Id string `json:"id"`

		// in: body
		// required: true
		Balance Balance
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
