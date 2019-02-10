package api

import (
	"net/http"
	"strings"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func PermissionRoutesPrivate(router *gin.RouterGroup, config k4ever.Config) {
	permissions := router.Group("/permissions/")
	{
		getPermissions(permissions, config)
		getPermission(permissions, config)
		createPermission(permissions, config)
	}
}

// swagger:route GET /permissions/ permissions getPermissions
//
// Lists all permsissions
//
// This will show all permissions by default
//
// 		Produces:
//		- application/json
//
//		Responses:
//		  default: GenericError
//		  200: PermissionsResponse
//        404: GenericError
func getPermissions(router *gin.RouterGroup, config k4ever.Config) {
	router.GET("", func(c *gin.Context) {
		var permissions []models.Permission
		if err := config.DB().Find(&permissions).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, permissions)
	})
}

// swagger:route GET /permissions/{id} permissions getPermission
//
// Get detailed information of a permission
//
// 		Produces:
//   	- application/json
//
//		Responses:
// 		  default: GenericError
//		  200: Permission
//		  404: GenericError
func getPermission(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters getPermission
	type getPermissionParams struct {
		// in: path
		// required: true
		Id int `json:"id"`
	}
	router.GET(":id", func(c *gin.Context) {
		var permission models.Permission
		if err := config.DB().Find(&permission).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, permission)
	})
}

// swagger:route POST /permission/ permissions createPermission
//
// Create a new permission
//
// Creating a permission has no real function yet, since they arent
// being checked anywhere
// This will probably be a role in the future with fixef permissions
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		Responses:
//		  default: GenericError
//		  201: Permission
//		  400: GenericError
//        500: GenericError
func createPermission(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters createPermission
	type CreatePermissionsParams struct {
		// in: body
		// required: true
		Permission models.Permission
	}
	router.POST("", func(c *gin.Context) {
		var permission models.Permission
		if err := c.ShouldBindJSON(&permission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := config.DB().Create(&permission).Error; err != nil {
			if strings.HasPrefix(err.Error(), "UNIQUE constraint failed:") {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Permission already exists"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, permission)
	})
}
