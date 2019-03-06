package k4ever

import (
	"github.com/freitagsrunde/k4ever-backend/internal/models"
)

func GetProducts(username string, config Config) (products []models.Product, err error) {
	if err := config.DB().Find(&products).Error; err != nil {
		return []models.Product{}, err
	}

	if err := config.DB().Table("purchase_items").Select("product_id as id, count(product_id) as times_bought_total").Group("product_id").Error; err != nil {
		return []models.Product{}, err
	}

	if err := config.DB().Table("purchase_items").Select("product_id as id, count(product_id) as times_bought").Joins("join purchase on purchase.id = purchase_items.purchase_id").Joins("join user on user.id = purchase.user_id").Where("user.user_name = ?", username).Group("product_id").Error; err != nil {
		return []models.Product{}, err
	}

	return products, nil
}

func CreateProduct(product *models.Product, config Config) (err error) {
	if err := config.DB().Create(product).Error; err != nil {
		return err
	}
	return nil
}

func BuyProduct(productID string, username string, config Config) (purchase models.Purchase, err error) {
	var product models.Product

	tx := config.DB().Begin()
	// Get Product
	if err := tx.Where("id = ?", productID).First(&product).Error; err != nil {
		return models.Purchase{}, err
	}

	purchase = models.Purchase{Total: product.Price}
	item := models.PurchaseItem{Amount: 1}
	item.ProductID = product.ID
	item.Name = product.Name
	item.Price = product.Price
	item.Description = product.Description
	item.Deposit = product.Deposit
	item.Barcode = product.Barcode
	item.Image = product.Image

	// Create PurchaseItem
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		return models.Purchase{}, err
	}
	purchase.Items = append(purchase.Items, item)
	// Create Purchase
	if err := tx.Create(&purchase).Error; err != nil {
		tx.Rollback()
		return models.Purchase{}, err
	}
	// Update Balance
	var user models.User
	if err := tx.Where("user_name = ?", username).First(&user).Error; err != nil {
		tx.Rollback()
		return models.Purchase{}, err
	}
	user.Balance = user.Balance - product.Price
	user.Purchases = append(user.Purchases, purchase)
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return models.Purchase{}, err
	}
	tx.Commit()
	return purchase, nil
}
