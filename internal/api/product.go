package api

import (
	"fmt"
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/freitagsrunde/k4ever-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func ProductRoutesPublic(router *gin.RouterGroup, config k4ever.Config) {
	products := router.Group("/products/")
	{
		getProduct(products, config)
		setProductImage(products, config)
	}
}

func ProductRoutesPrivate(router *gin.RouterGroup, config k4ever.Config) {
	products := router.Group("/products/")
	{
		getProducts(products, config)
		createProduct(products, config)
		buyProduct(products, config)
	}
}

// swagger:route GET /products/ products getProducts
//
// Lists all available prodcuts
//
// This will show all available products by default
//
// 		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		Responses:
//		  default: GenericError
//		  200: productsResponse
//		  404: GenericError
func getProducts(router *gin.RouterGroup, config k4ever.Config) {
	// A ProductsResponse returns a list of products
	//
	// swagger:response productsResponse
	type ProductsResponse struct {
		// An array of products
		//
		// in: body
		Products []models.Product
	}
	router.GET("", func(c *gin.Context) {
		var err error
		params := models.DefaultParams{}
		params.SortBy = c.DefaultQuery("sort_by", "name")
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
		claims := jwt.ExtractClaims(c)
		username := claims["name"]

		products, err := k4ever.GetProducts(username.(string), params, config)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, GenericError{Body: struct{ Message string }{Message: err.Error()}})
			return
		}
		c.JSON(http.StatusOK, products)
	})
}

// swagger:route GET /products/{id}/ products getProduct
//
// Get information for a product by id
//
// This will show detailed information for a specific product
//
// 		Consumes:
// 		- application/json
//
//		Produces:
//		- application/json
//
//		Responses:
//		  default: GenericError
//		  200: Product
//		  404: GenericError
func getProduct(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters getProduct
	type getProductParams struct {
		// in: path
		// required: true
		Id int `json:"id"`
	}
	router.GET(":id/", func(c *gin.Context) {
		var product models.Product
		if err := config.DB().Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, product)
	})
}

// swagger:route POST /products/ products createProduct
//
// Create a new product
//
// Create a new product (currently with all fields available)
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
//        default: GenericError
//		  201: Product
//		  400: GenericError
//        500: GenericError
func createProduct(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters createProduct
	type ProductParam struct {
		// in: body
		// required: true
		Product models.Product
	}
	router.POST("", func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := config.DB().Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, product)
	})
}

// swagger:route GET /products/{id}/image/ getProductImage
//
// Not yet implemented
//
// Returns a product image or path to it (tbd)
//
// 		Produces:
//		- application/json
//
//		Responses:
//		  default: GenericError
//		  502: GenericError
func setProductImage(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters getProductImage
	type getProductImageParams struct {
		// in: path
		// required: true
		Id int `json:"id"`
	}
	router.PUT(":id/image/", func(c *gin.Context) {
		var product models.Product
		if err := config.DB().First(&product).Where("id = ?", c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "no such product"})
			return
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not get file from key file"})
			return
		}

		f, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not open file"})
			return
		}

		bytes := utils.StreamToByte(f)

		path, err := utils.UploadFile(bytes, "products/"+product.Name, fileHeader.Filename)
		if err != nil {
			fmt.Println(err.Error())
			if err.Error() == "file already exists" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "file already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while saving file"})
			return
		}
		fmt.Println(path)

		product.Image = config.HttpServerHost() + ":" + strconv.Itoa(config.HttpServerPort()) + path
		k4ever.UpdateProduct(&product, config)
		c.JSON(http.StatusOK, product)
	})
}

// swagger:route POST /products/{id}/buy/ buyProduct
//
// Buy a product as the current user
//
// Buys a product according to the user read from the jwt header
//
//		Produces:
//		- application/json
//
//		Security:
//		  jwt:
//
//		Responses:
//		  default: GenericError
//		  200: Purchase
//		  400: GenericError
//		  404: GenericError
//        500: GenericError
func buyProduct(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters buyProduct
	type buyProductParams struct {
		// in: path
		// required: true
		Id int `json:"id"`
	}
	router.POST(":id/buy/", func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		username := claims["name"]
		if username == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}
		purchase, err := k4ever.BuyProduct(c.Param("id"), username.(string), config)
		if err != nil {
			if err.Error() == "record not found" {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Record not found"})
				return
			}
			if err.Error() == "The product is disabled and cannot be bought" {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, purchase)
	})
}
