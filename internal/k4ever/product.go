package k4ever

import (
	"errors"

	"github.com/freitagsrunde/k4ever-backend/internal/models"
)

func GetProducts(username string, params models.DefaultParams, config Config) (products []models.Product, err error) {
	// Subquery to get the sum of histories for each product
	sumProductsTotal := config.DB().Table("purchase_items").Select("sum(amount)").Group("product_id").Where("purchase_items.product_id = p.id").QueryExpr()

	// Subquery to get the sum of histories by the logged in user for each product
	sumProductsUser := config.DB().Table("purchase_items").Select("sum(purchase_items.amount)").Joins("join histories on histories.id = purchase_items.history_id").Joins("join users on users.id = histories.user_id").Where("users.user_name = ? AND purchase_items.product_id = p.id", username).Group("purchase_items.product_id").QueryExpr()

	// Subquery to get the time when the current user last bought the item
	lastBoughtByUser := config.DB().Table("purchase_items").Select("purchase_items.updated_at").Joins("join histories on histories.id = purchase_items.history_id").Joins("join users on users.id = histories.user_id").Where("users.user_name = ? AND purchase_items.product_id = p.id", username).Order("purchase_items.updated_at desc").Limit(1).QueryExpr()

	// Query to get all product information
	tx := config.DB().Table("products p").Select("*, COALESCE((?), 0) as times_bought_total, COALESCE((?), 0) as times_bought, (?) as last_bought", sumProductsTotal, sumProductsUser, lastBoughtByUser).Group("id").Order(params.SortBy + " " + params.Order)
	if params.Offset != 0 {
		tx = tx.Offset(params.Offset)
	}
	if params.Limit != 0 {
		tx = tx.Limit(params.Limit)
	}

	rows, err := tx.Rows()
	if err != nil {
		return []models.Product{}, err
	}
	for rows.Next() {
		var p models.Product
		if errSql := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.Name, &p.Price, &p.Description, &p.Deposit, &p.Barcode, &p.Image, &p.Disabled, &p.TimesBoughtTotal, &p.TimesBought, &p.LastBought); errSql != nil {
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

	if err := config.DB().Table("purchase_items").Select("purchase_items.product_id, count(purchase_items.product_id)").Joins("join histories on histories.id = purchase_items.history_id").Joins("join users on users.id = histories.user_id").Where("users.user_name = ?", username).Group("purchase_items.product_id").Count(&count).Error; err != nil {
		return models.Product{}, err
	}
	product.TimesBought = count

	// Subquery to get the time when the current user last bought the item
	if err := config.DB().Table("purchase_items").Select("purchase_items.updated_at as last_bought").Joins("join histories on histories.id = purchase_items.history_id").Joins("join users on users.id = histories.user_id").Where("users.user_name = ? AND purchase_items.product_id = ?", username, productID).Order("purchase_items.updated_at desc").Limit(1).Scan(&product).Error; err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func CreateProduct(product *models.Product, config Config) (err error) {
	if err := config.DB().Create(product).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProduct(product *models.Product, config Config) (err error) {
	if err := config.DB().Update(product).Error; err != nil {
		return err
	}
	return nil
}

func BuyProduct(productID string, username string, config Config) (purchase models.History, err error) {
	var product models.Product

	tx := config.DB().Begin()
	// Get Product
	if err := tx.Where("id = ?", productID).First(&product).Error; err != nil {
		return models.History{}, err
	}

	if product.Disabled == true {
		return models.History{}, errors.New("The product is disabled and cannot be bought")
	}

	purchase = models.History{Total: product.Price, Type: models.PurchaseHistory}
	item := models.PurchaseItem{Amount: 1}
	item.ProductID = product.ID
	item.Name = product.Name
	item.Price = product.Price
	item.Description = product.Description
	item.Deposit = product.Deposit
	item.Barcode = product.Barcode.String
	item.Image = product.Image

	// Create PurchaseItem
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		return models.History{}, err
	}
	purchase.Items = append(purchase.Items, item)
	// Create Purchase
	if err := tx.Create(&purchase).Error; err != nil {
		tx.Rollback()
		return models.History{}, err
	}
	// Update Balance
	var user models.User
	if err := tx.Where("user_name = ?", username).First(&user).Error; err != nil {
		tx.Rollback()
		return models.History{}, err
	}
	user.Balance = user.Balance - product.Price
	user.Histories = append(user.Histories, purchase)
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return models.History{}, err
	}
	tx.Commit()
	return purchase, nil
}
