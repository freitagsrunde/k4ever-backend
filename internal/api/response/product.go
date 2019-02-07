package response

import "github.com/freitagsrunde/k4ever-backend/internal/models"

// A ProductsResponse returns a list of products
//
// swagger:response productsResponse
type ProductsResponse struct {
	// An array of products
	//
	// in: body
	Products []models.Product
}
