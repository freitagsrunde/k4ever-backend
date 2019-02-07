package params

import "github.com/freitagsrunde/k4ever-backend/internal/models"

// swagger:parameters createProduct
type ProductParam struct {
	// in: body
	// required: true
	Product models.Product
}
