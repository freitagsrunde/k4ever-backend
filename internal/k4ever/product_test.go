package k4ever

import (
	"strconv"
	"testing"

	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	conf := NewK4everTest()

	testProduct := ProductTest()
	err := CreateProduct(&testProduct, conf)

	assert.Equal(t, nil, err)
	assert.Equal(t, uint(1), testProduct.ID)
}

func TestBuyProduct(t *testing.T) {
	conf := NewK4everTest()

	testProduct := ProductTest()
	err := CreateProduct(&testProduct, conf)

	assert.Equal(t, nil, err)

	testUser := UserTest()
	err2 := CreateUser(&testUser, conf)

	assert.Equal(t, nil, err2)

	purchase, err3 := BuyProduct(strconv.Itoa(int(testProduct.ID)), testUser.UserName, conf)

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

	testProduct := ProductTest()
	err := CreateProduct(&testProduct, conf)

	assert.Equal(t, nil, err)

	testUser := UserTest()
	err2 := CreateUser(&testUser, conf)

	assert.Equal(t, nil, err2)

	_, err3 := BuyProduct(strconv.Itoa(int(testProduct.ID)), testUser.UserName, conf)

	assert.Equal(t, nil, err3)

	products, err4 := GetProducts(testUser.UserName, conf)

	assert.Equal(t, nil, err4)
	assert.Equal(t, 1, len(products))
	assert.Equal(t, 1, products[0].TimesBoughtTotal)
	assert.Equal(t, 1, products[0].TimesBought)
}

func TestGetProductEmpty(t *testing.T) {
	conf := NewK4everTest()

	product, err := GetProduct("1", "name", conf)

	assert.Equal(t, models.Product{}, product)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "record not found", err.Error())
}

func TestGetProduct(t *testing.T) {
	conf := NewK4everTest()

	testProduct := ProductTest()
	err := CreateProduct(&testProduct, conf)

	assert.Equal(t, nil, err)

	testUser := UserTest()
	err2 := CreateUser(&testUser, conf)

	assert.Equal(t, nil, err2)

	_, err3 := BuyProduct(strconv.Itoa(int(testProduct.ID)), testUser.UserName, conf)

	assert.Equal(t, nil, err3)

	product, err4 := GetProduct(strconv.Itoa(int(testProduct.ID)), testUser.UserName, conf)

	assert.Equal(t, nil, err4)
	assert.Equal(t, 1, product.TimesBoughtTotal)
	assert.Equal(t, 1, product.TimesBought)
}
