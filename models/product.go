package models

import (
	"net/http"

	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string
	Price       float64
	Deposit     float64
	Barcode     int `gorm:"unique"`
	//Type        []Type `gorm:"many2many:product_types;"`
	Archived bool
}

func (p *Product) Bind(r *http.Request) error {
	return nil
}

type Producter interface {
	ListProducts() ([]Product, error)
	GetProduct(id string) (Product, error)
	CreateProduct(product *Product) error
	UpdateProduct(product *Product) error
}

type ProductResource struct {
	DB *gorm.DB
}

func (pr ProductResource) ListProducts() ([]Product, error) {
	var products []Product
	if err := pr.DB.Find(&products).Error; err != nil {
		return []Product{}, err
	}
	return products, nil
}

func (pr ProductResource) GetProduct(id uint) (Product, error) {
	var product Product
	if err := pr.DB.First(&product, "id = ?", id).Error; err != nil {
		return Product{}, err
	}
	return product, nil
}

func (pr ProductResource) CreateProduct(product *Product) error {
	if err := pr.DB.Create(product).Error; err != nil {
		return err
	}
	return nil
}

func (pr ProductResource) UpdateProduct(product *Product) error {
	if err := pr.DB.Model(product).Updates(product).Error; err != nil {
		return err
	}
	if err := pr.DB.First(product).Error; err != nil {
		return err
	}
	return nil
}
