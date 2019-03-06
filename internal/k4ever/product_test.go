package k4ever

import (
	"testing"

	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	conf := NewK4everTest()

	product := models.Product{}
	product.Name = "product"
	product.Price = 1.0
	product.Description = "A description"
	product.Deposit = 0.0
	product.Barcode = "12345678"

	err := CreateProduct(&product, conf)

	assert.Equal(t, nil, err)
	assert.Equal(t, uint(1), product.ID)
}

func TestBuyProduct(t *testing.T) {
	conf := NewK4everTest()

	product := models.Product{}
	product.Name = "product"
	product.Price = 1.0
	product.Description = "A description"
	product.Deposit = 0.0
	product.Barcode = "12345678"

	err := CreateProduct(&product, conf)

	assert.Equal(t, nil, err)

	user := models.User{}
	user.UserName = "user"
	user.Password = "password"
	user.DisplayName = "displayname"

	err2 := CreateUser(&user, conf)

	assert.Equal(t, nil, err2)

	purchase, err3 := BuyProduct("1", "user", conf)

	assert.Equal(t, nil, err3)
	assert.Equal(t, uint(1), purchase.ID)
}

func TestGetProductsEmpty(t *testing.T) {
	conf := NewK4everTest()

	products, err := GetProducts("name", conf)

	assert.Equal(t, 0, len(products))
	assert.Equal(t, nil, err)
}

func TestGetProducts(t *testing.T) {
	conf := NewK4everTest()

	product := models.Product{}
	product.Name = "product"
	product.Price = 1.0
	product.Description = "A description"
	product.Deposit = 0.0
	product.Barcode = "12345678"

	err := CreateProduct(&product, conf)

	assert.Equal(t, nil, err)

	user := models.User{}
	user.UserName = "user"
	user.Password = "password"
	user.DisplayName = "displayname"

	err2 := CreateUser(&user, conf)

	assert.Equal(t, nil, err2)

	_, err3 := BuyProduct("1", "user", conf)

	assert.Equal(t, nil, err3)

	products, err4 := GetProducts("name", conf)

	assert.Equal(t, nil, err4)
	assert.Equal(t, 1, len(products))
	assert.Equal(t, 1, products[0].TimesBoughtTotal)
}
