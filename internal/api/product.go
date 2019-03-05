package api

import (
	"net/http"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/gin-gonic/gin"
)

func ProductRoutesPublic(router *gin.RouterGroup, config k4ever.Config) {
	products := router.Group("/products/")
	{
		getProducts(products, config)
		getProduct(products, config)
		getProductImage(products, config)
	}
}

func ProductRoutesPrivate(router *gin.RouterGroup, config k4ever.Config) {
	products := router.Group("/products/")
	{
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
		var products []models.Product
		if err := config.DB().Find(&products).Error; err != nil {
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
			if strings.Contains(err.Error(), "UNIQUE") {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
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
func getProductImage(router *gin.RouterGroup, config k4ever.Config) {
	// swagger:parameters getProductImage
	type getProductImageParams struct {
		// in: path
		// required: true
		Id int `json:"id"`
	}
	router.GET(":id/image/", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"Hello": "World"})
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
		var product models.Product
		tx := config.DB().Begin()
		// Get Product
		if err := tx.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if product.Disabled == true {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The product is disabled and cannot be bought"})
			return
		}
    
    purchase := models.Purchase{Total: product.Price}
		item := models.PurchaseItem{Amount: 1, Product: product, ProductID: product.ID}
		item := models.PurchaseItem{Amount: 1}
		item.ProductID = product.ID
		item.Name = product.Name
		item.Price = product.Price
		item.Description = product.Description
		item.Deposit = product.Deposit
		item.Barcode = product.Barcode
		item.Image = product.Image
		// Create PurchaseItem
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		purchase.Items = append(purchase.Items, item)
		// Create Purchase
		if err := tx.Create(&purchase).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Update Balance
		var user models.User
		claims := jwt.ExtractClaims(c)
		userID := claims["id"]
		if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}
		user.Balance = user.Balance - product.Price
		user.Purchases = append(user.Purchases, purchase)
		if err := tx.Save(&user).Error; err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tx.Commit()
		c.JSON(http.StatusOK, purchase)
	})
}
