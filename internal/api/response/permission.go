package response

import "github.com/freitagsrunde/k4ever-backend/internal/models"

// A PermissionsResponse returns a list of products
//
// swagger:response
type PermissionsResponse struct {
	// An array of permissions
	//
	// in: body
	Permissions []models.Permission
}
