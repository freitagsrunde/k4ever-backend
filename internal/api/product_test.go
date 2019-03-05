package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	app, router, conf := NewApiTest()

	getProducts(router, conf)

	result := PerformRequest(app, "GET", "/api/v1")

	assert.Equal(t, http.StatusOK, result.Code)
}

func TestCreateProduct(t *testing.T) {
	app, router, conf := NewApiTest()

	createProduct(router, conf)

	var json = []byte(`{
		"name":"Product",
		"price":1.50
	}`)

	result := PerformRequestWithBody(app, "POST", "/api/v1", json)

	assert.Equal(t, http.StatusCreated, result.Code)

}

func TestCreateDuplicateProduct(t *testing.T) {
	app, router, conf := NewApiTest()

	createProduct(router, conf)

	var json = []byte(`{
		"name":"Product",
		"price":1.50
	}`)

	result := PerformRequestWithBody(app, "POST", "/api/v1", json)
	result2 := PerformRequestWithBody(app, "POST", "/api/v1", json)

	assert.Equal(t, http.StatusCreated, result.Code)
	assert.Equal(t, http.StatusOK, result2.Code)
}
