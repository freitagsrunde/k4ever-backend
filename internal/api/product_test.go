package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	app, router, conf := NewApiTest()

	getProducts(router, conf)

	result := PerformRequest(app, "GET", "/api/v1/products/")

	assert.Equal(t, http.StatusOK, result.Code)
}
