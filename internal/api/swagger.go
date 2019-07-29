package api

import (
	"fmt"
	"net/http"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
)

func SwaggerRoutesPublic(router *gin.RouterGroup, config k4ever.Config) {
	swagger := router.Group("swagger/")
	{
		getSwagger(swagger, config)
	}
}

// Just a swagger file
//
// swagger:model
type SwaggerResponse struct {
	// A swagger file
	//
	// in: body
	Swagger string
}

// swagger:route GET /swagger/ swagger getSwagger
//
// Get the swagger yml
//
//		Produces:
//		- application/yml
//
//		Responses:
//		  default: SwaggerResponse
//		  200: SwaggerResponse
func getSwagger(router *gin.RouterGroup, config k4ever.Config) {
	box := packr.NewBox("../../")
	s, err := box.FindString("swagger.yml")
	if err != nil {
		fmt.Println(err.Error())
		s = "undefined"
	}

	router.GET("", func(c *gin.Context) {
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		swagger := SwaggerResponse{Swagger: s}
		c.String(http.StatusOK, swagger.Swagger)
	})
}
