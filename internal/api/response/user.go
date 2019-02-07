package response

import "github.com/freitagsrunde/k4ever-backend/internal/models"

// A UsersResponse returns a list of users
//
// swagger:response
type UsersResponse struct {
	// An array of products
	//
	// in: body
	Users []models.User
}
