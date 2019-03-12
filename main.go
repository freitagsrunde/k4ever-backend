// k4ever api.
//
// defines all k4ever backend endpoints
//
//		Schemes: http, https
//		Host: localhost
//		BasePath: /api/v1
//		Version: 0.0.1
//		License: MIT (tbd)
//		Contact: Phillip Stagnet <phillip@freitagsrunde.org>
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//      SecurityDefinitions:
//        jwt:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package main

import (
	"github.com/freitagsrunde/k4ever-backend/cmd"
)

//go:generate swagger generate spec -o ./swagger.yml
func main() {
	cmd.Execute()
}
