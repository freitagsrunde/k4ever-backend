package k4ever

import (
	"github.com/freitagsrunde/k4ever-backend/internal/models"
)

func GetProducts(username string, order string, config Config) (products []models.Product, err error) {
	rows, err := config.DB().Table("products p").Select("*, COALESCE((?), 0) as most_bought, COALESCE((?), 0) as most_bought_total", config.DB().Table("purchase_items").Select("sum(amount)").Group("product_id").Where("purchase_items.product_id = p.id").QueryExpr(), config.DB().Table("purchase_items").Select("sum(purchase_items.amount)").Joins("join purchases on purchases.id = purchase_items.purchase_id").Joins("join users on users.id = purchases.user_id").Where("users.user_name = ? AND purchase_items.product_id = p.id", username).Group("purchase_items.product_id").QueryExpr()).Group("id").Order("order").Rows()
	if err != nil {
		return []models.Product{}, err
	}
	for rows.Next() {
		var p models.Product
		if errSql := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.Name, &p.Price, &p.Description, &p.Deposit, &p.Barcode, &p.Image, &p.TimesBoughtTotal, &p.TimesBought); errSql != nil {
			return []models.Product{}, errSql
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return []models.Product{}, err
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
