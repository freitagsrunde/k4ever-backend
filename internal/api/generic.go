package api

// A generic error message
// swagger:response
type GenericError struct {
	// The error message
	// in: body
	Body struct {
		// The error message
		//
		// Required: true
		// Example: Something went wrong
		Message string
	}
}
