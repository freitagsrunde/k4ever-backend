package k4ever

import (
	"errors"

	"github.com/freitagsrunde/k4ever-backend/internal/models"
)

func GetProducts(username string, params models.DefaultParams, config Config) (products []models.Product, err error) {
	// Subquery to get the sum of purchases for each product
	sumProductsTotal := config.DB().Table("purchase_items").Select("sum(amount)").Group("product_id").Where("purchase_items.product_id = p.id").QueryExpr()

	// Subquery to get the sum of purchases by the logged in user for each product
	sumProductsUser := config.DB().Table("purchase_items").Select("sum(purchase_items.amount)").Joins("join purchases on purchases.id = purchase_items.purchase_id").Joins("join users on users.id = purchases.user_id").Where("users.user_name = ? AND purchase_items.product_id = p.id", username).Group("purchase_items.product_id").QueryExpr()

	// Query to get all product information
	tx := config.DB().Preload("Users").Table("products p").Select("*, COALESCE((?), 0) as times_bought_total, COALESCE((?), 0) as times_bought", sumProductsTotal, sumProductsUser).Group("id").Order(params.SortBy + " " + params.Order)
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
		if errSql := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.Name, &p.Price, &p.Description, &p.Deposit, &p.Barcode, &p.Image, &p.Disabled, &p.TimesBoughtTotal, &p.TimesBought); errSql != nil {
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
	if err := config.DB().Preload("Users").Where("id = ?", productID).First(&product).Error; err != nil {
		return models.Product{}, err
	}

	var count int
	if err := config.DB().Table("purchase_items").Select("product_id, count(product_id)").Group("product_id").Count(&count).Error; err != nil {
		return models.Product{}, err
	}
	product.TimesBoughtTotal = count

	// Check wether product is liked by current User
	var isLiked int
	if err := config.DB().Table("users").Joins("JOIN liked_by ON liked_by.user_id = users.id").Joins("JOIN products on products.id = liked_by.product_id").Count(&isLiked).Error; err != nil {
		return models.Product{}, err
	}
	// The number of rows will be at most 1
	if isLiked > 0 {
		product.IsLiked = true
	}

	if err = config.DB().Table("purchase_items").Select("purchase_items.product_id, count(purchase_items.product_id)").Joins("join purchases on purchases.id = purchase_items.purchase_id").Joins("join users on users.id = purchases.user_id").Where("users.user_name = ?", username).Group("purchase_items.product_id").Count(&count).Error; err != nil {
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

	if product.Disabled == true {
		return models.Purchase{}, errors.New("The product is disabled and cannot be bought")
	}

	purchase = models.Purchase{Total: product.Price}
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

func LikeProduct(productID string, username string, config Config) (product models.Product, err error) {
	tx := config.DB().Begin()
	// Get Product
	if err = tx.Where("id = ?", productID).First(&product).Error; err != nil {
		tx.Rollback()
		return models.Product{}, err
	}

	// Get User
	var user models.User
	if err = tx.Where("user_name = ?", username).First(&user).Error; err != nil {
		tx.Rollback()
		return models.Product{}, err
	}

	// Check if user is already in list
	for _, v := range product.Users {
		if v.UserName == user.UserName {
			return models.Product{}, errors.New("Already liked product")
		}
	}

	// Update LikedBy list
	product.Users = append(product.Users, user)
	if err = tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return models.Product{}, err
	}

	// Commit queries
	tx.Commit()
	product.IsLiked = true
	return product, nil
}
