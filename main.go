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
// swagger:meta
package main

import (
	"github.com/freitagsrunde/k4ever-backend/cmd"
)

func main() {
	cmd.Execute()
}
