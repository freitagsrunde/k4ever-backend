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

// swagger:route GET /swagger/ swagger getSwagger
//
// Get the swagger yml
//
//		Produces:
//		- application/yml
//
//		Responses:
//		  default: swaggerResponse
//		  200: swaggerResponse
func getSwagger(router *gin.RouterGroup, config k4ever.Config) {
	// Just a swagger file
	type swaggerResponse struct {
		// A swagger file
		//
		// in: body
		Swagger string
	}
	box := packr.NewBox("../../")
	s, err := box.FindString("swagger.yml")
	if err != nil {
		fmt.Println(err.Error())
		s = "undefined"
	}

	router.GET("", func(c *gin.Context) {
		c.Header("Content-Type", "application/yaml; charset=utf-8")
		c.String(http.StatusOK, s)
	})
}
