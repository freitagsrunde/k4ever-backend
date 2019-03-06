package k4ever

import (
	"github.com/freitagsrunde/k4ever-backend/internal/models"
)

func GetProducts(username string, config Config) (products []models.Product, err error) {
	if err := config.DB().Find(&products).Error; err != nil {
		return []models.Product{}, err
	}

	productMap := make(map[uint]models.Product)
	var i int
	for i = 0; i < len(products); i++ {
		productMap[products[i].ID] = products[i]
	}

	rows, err := config.DB().Table("purchase_items").Select("product_id, count(product_id)").Group("product_id").Rows()
	if err != nil {
		return []models.Product{}, err
	}

	for rows.Next() {
		var id uint
		var count int
		if errSql := rows.Scan(&id, &count); errSql != nil {
			return []models.Product{}, errSql
		}
		product := productMap[id]
		product.TimesBoughtTotal = count
		productMap[id] = product
	}

	if err := rows.Err(); err != nil {
		return []models.Product{}, err
	}

	rows2, err2 := config.DB().Table("purchase_items").Select("purchase_items.product_id, count(purchase_items.product_id)").Joins("join purchases on purchases.id = purchase_items.purchase_id").Joins("join users on users.id = purchases.user_id").Where("users.user_name = ?", username).Group("purchase_items.product_id").Rows()
	if err2 != nil {
		return []models.Product{}, err2
	}
	for rows2.Next() {
		var id uint
		var count int
		if errSql := rows2.Scan(&id, &count); errSql != nil {
			return []models.Product{}, errSql
		}
		product := productMap[id]
		product.TimesBought = count
		productMap[id] = product
	}

	if err := rows2.Err(); err != nil {
		return []models.Product{}, err
	}

	for i = 0; i < len(products); i++ {
		products[i] = productMap[products[i].ID]
	}

	return products, nil
}

func GetProduct(productID string, username string, config Config) (product models.Product, err error) {
	if err := config.DB().First(&product).Where("id = ?", productID).Error; err != nil {
		return models.Product{}, err
	}

	var count int
	if err := config.DB().Table("purchase_items").Select("product_id, count(product_id)").Group("product_id").Count(&count).Error; err != nil {
		return models.Product{}, err
	}
	product.TimesBoughtTotal = count

	if err := config.DB().Table("purchase_items").Select("purchase_items.product_id, count(purchase_items.product_id)").Joins("join purchases on purchases.id = purchase_items.purchase_id").Joins("join users on users.id = purchases.user_id").Where("users.user_name = ?", username).Group("purchase_items.product_id").Count(&count).Error; err != nil {
		return models.Product{}, err
	}
	product.TimesBought = count

	return product, nil
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
