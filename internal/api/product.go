package api

import (
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/freitagsrunde/k4ever-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func ProductRoutesPublic(router *gin.RouterGroup, config k4ever.Config) {
	products := router.Group("/products/")
	{
		getProduct(products, config)
	}
}

func ProductRoutesPrivate(router *gin.RouterGroup, config k4ever.Config) {
	products := router.Group("/products/")
	{
		getProducts(products, config)
		createProduct(products, config)
		updateProduct(products, config)
		deleteProduct(products, config)
		buyProduct(products, config)
		setProductImage(products, config)
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
//		Security:
//		  jwt:
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
		if !utils.CheckRole(0, c) {
			return
		}
		params, err := utils.ParseDefaultParams(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		claims := jwt.ExtractClaims(c)
		username := claims["name"]
		// Since at this point a user was already validated we can use a default user if non is found for testing
		if username == nil {
			username = ""
		}

		products, err := k4ever.GetProducts(username.(string), params, config)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, GenericError{Body: struct{ Message string }{Message: err.Error()}})
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
//		Security:
//		  jwt:
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
		if !utils.CheckRole(2, c) {
			return
		}
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

// swagger:route PUT /products/{id}/ updateProduct
//
// Update a product
//
// 		Produces:
//		- application/json
//
//		Consumes:
//		- application/json
//
//		Responses:
//		  default: GenericError
//		  200: Product
//		  400: GenericError
//		  401: GenericError
//		  500: GenericError
func updateProduct(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters updateProduct
	type ProductParam struct {
		//in: body
		// required: true
		Product models.Product
	}
	router.PUT(":id/", func(c *gin.Context) {
		if !utils.CheckRole(2, c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		uintID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		product.ID = uint(uintID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is not an int"})
			return
		}

		if err = k4ever.UpdateProduct(&product, config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	})
}

// swager:route DELETE /products/{id}/ deleteProduct
//
// Delete a product
//
//		Responses:
//		  default: GenericError
//		  200: string
//		  401: GenericError
//		  500: GenericError
func deleteProduct(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters deleteProduct
	type deleteProductParams struct {
		// in: path
		// required: true
		Id int
	}
	router.DELETE(":id/", func(c *gin.Context) {
		if !utils.CheckRole(2, c) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		uintID, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Id is not an int"})
			return
		}

		if err = k4ever.DeleteProduct(uint(uintID), config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "Deleted")
	})
}

// swagger:route PUT /products/{id}/image/ setProductImage
//
// set the product image for a single product
//
// Set the product image from the form value "file"
//
// 		Produces:
//		- application/json
//
//		Consumes:
//		- multipart/form-data
//
//		Security:
//		  jwt:
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

		// in: form
		// required: true
		File string `json:"file"`
	}
	router.PUT(":id/image/", func(c *gin.Context) {
		if !utils.CheckRole(2, c) {
			return
		}
		var product models.Product
		if err := config.DB().Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		imagePath := setImage(c, product, config)

		// Return if empty string is returned because have already send an error message
		if imagePath == "" {
			return
		}

		// Set prouct image
		product.Image = imagePath
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
//		  200: History
//		  400: GenericError
//		  404: GenericError
//        500: GenericError
func buyProduct(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters buyProduct
	type buyProductParams struct {
		// in: path
		// required: true
		Id int `json:"id"`

		// in: body
		// required: false
		Deposit bool `json:"deposit"`
	}
	router.POST(":id/buy/", func(c *gin.Context) {
		if !utils.CheckRole(1, c) {
			return
		}
		claims := jwt.ExtractClaims(c)
		username := claims["name"]
		if username == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}

		var productParams buyProductParams
		if c.Request.Body == http.NoBody {
			if err := c.ShouldBindJSON(&productParams); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		purchase, err := k4ever.BuyProduct(c.Param("id"), productParams.Deposit, username.(string), config)
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
