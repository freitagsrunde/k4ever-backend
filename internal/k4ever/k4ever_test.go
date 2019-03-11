package k4ever

import (
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/freitagsrunde/k4ever-backend/internal/test"
)

func NewK4everTest() (conf Config) {
	conf = test.NewConfig()
	conf.MigrateDB()

	return conf
}

func ProductTest() (testProduct models.Product) {
	testProduct.Name = "product"
	testProduct.Price = 1.0
	testProduct.Description = "A description"
	testProduct.Deposit = 0.0
	testProduct.Barcode = "12345678"

	return testProduct
}

func UserTest() (testUser models.User) {
	testUser.UserName = "user"
	testUser.Password = "password"
	testUser.DisplayName = "displayname"

	return testUser
}

func DefaultParamsTest() (params models.DefaultParams) {
	params.Order = "asc"

	return params
}
