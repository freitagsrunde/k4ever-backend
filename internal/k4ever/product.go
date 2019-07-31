package k4ever

import (
	"encoding/json"
	"errors"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/tidwall/gjson"
)

func GetProducts(username string, params models.DefaultParams, config Config) (products []models.Product, err error) {
	q := `
		{
			products(func: has(product)) {
				uid
				expand(_all_)	
			}
		}
	`
	txn := config.DB().NewReadOnlyTxn().BestEffort()
	resp, err := txn.Query(config.Context(), q)
	if err != nil {
		return []models.Product{}, err
	}

	var decode struct {
		Products []models.Product
	}

	if err = json.Unmarshal(resp.GetJson(), &decode); err != nil {
		return []models.Product{}, err
	}

	return decode.Products, nil

	// Subquery to get the sum of histories for each product
	//sumProductsTotal := config.DB().Table("purchase_items").Select("sum(amount)").Group("product_id").Where("purchase_items.product_id = p.id").QueryExpr()

	// Subquery to get the sum of histories by the logged in user for each product
	//sumProductsUser := config.DB().Table("purchase_items").Select("sum(purchase_items.amount)").Joins("join histories on histories.id = purchase_items.history_id").Joins("join users on users.id = histories.user_id").Where("users.user_name = ? AND purchase_items.product_id = p.id", username).Group("purchase_items.product_id").QueryExpr()

	// Subquery to get the time when the current user last bought the item
	//lastBoughtByUser := config.DB().Table("purchase_items").Select("purchase_items.updated_at").Joins("join histories on histories.id = purchase_items.history_id").Joins("join users on users.id = histories.user_id").Where("users.user_name = ? AND purchase_items.product_id = p.id", username).Order("purchase_items.updated_at desc").Limit(1).QueryExpr()

	// Query to get all product information
	//tx := config.DB().Table("products p").Select("*, COALESCE((?), 0) as times_bought_total, COALESCE((?), 0) as times_bought, (?) as last_bought", sumProductsTotal, sumProductsUser, lastBoughtByUser).Group("id").Order(params.SortBy + " " + params.Order)
}

func GetProduct(productID string, username string, config Config) (product models.Product, err error) {
	q := `
		query withvar($productID: string){
			product(func: uid($productID)) {
				uid
				expand(_all_)
			}
		}
	`

	txn := config.DB().NewReadOnlyTxn().BestEffort()
	resp, err := txn.QueryWithVars(config.Context(), q, map[string]string{"$productID": productID})
	if err != nil {
		return models.Product{}, err
	}

	var decode struct {
		Product []models.Product
	}

	if err = json.Unmarshal(resp.GetJson(), &decode); err != nil {
		return models.Product{}, err
	}

	return decode.Product[0], nil
	/*if err := config.DB().First(&product).Where("id = ?", productID).Error; err != nil {
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
	}*/
}

func CreateProduct(product *models.Product, config Config) (err error) {
	checkingQuery := `
		query withvar($name: string){
			product(func: eq(name, $name)) @filter(has(product)){
				uid
			}
		}
	`
	txn := config.DB().NewTxn()
	defer txn.Discard(config.Context())
	resp, err := txn.QueryWithVars(config.Context(), checkingQuery, map[string]string{"$name": product.Name})
	if err != nil {
		return err
	}

	if length := gjson.Get(string(resp.GetJson()), "product.#"); length.Num > 0 {
		return errors.New("product already exists")
	}

	pg := &models.ProductDgraph{*product, true}
	mu := &api.Mutation{}

	pb, err := json.Marshal(pg)
	if err != nil {
		return err
	}
	mu.SetJson = pb

	_, err = txn.Mutate(config.Context(), mu)
	if err != nil {
		return err
	}
	txn.Commit(config.Context())

	return nil
}

func UpdateProduct(product *models.Product, config Config) (err error) {
	/*if err := config.DB().Update(product).Error; err != nil {
		return err
	}*/
	return nil
}

func BuyProduct(productID string, username string, config Config) (purchase models.History, err error) {
	/*var product models.Product

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
	tx.Commit()*/
	return purchase, nil
}
