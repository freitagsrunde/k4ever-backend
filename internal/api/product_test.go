package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	app, router, conf := NewApiTest()

	getProducts(router, conf)

	result := PerformRequest(app, "GET", "/api/v1")

	fmt.Println(router)

	assert.Equal(t, http.StatusOK, result.Code)
}
