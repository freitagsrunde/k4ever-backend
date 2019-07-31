package k4ever

import (
	"encoding/json"
	"errors"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(params models.DefaultParams, config Config) (users []models.User, err error) {
	q := `
		{
			users(func: has(user)) {
				uid
				expand(_all_)
			}
		}
	`

	txn := config.DB().NewReadOnlyTxn().BestEffort()
	resp, err := txn.Query(config.Context(), q)
	if err != nil {
		return
	}

	var decode struct {
		Users []models.User
	}

	if err = json.Unmarshal(resp.GetJson(), &decode); err != nil {
		return
	}

	return decode.Users, nil
}

func GetUser(name string, config Config) (user models.User, err error) {
	q := `
		query withvar($name: string) {
			user(func: has(user)) @filter(eq(name, $name)) {
				uid
				expand(_all_)
			}
		}
	`

	txn := config.DB().NewReadOnlyTxn().BestEffort()
	resp, err := txn.QueryWithVars(config.Context(), q, map[string]string{"$name": name})
	if err != nil {
		return
	}

	var decode struct {
		User []models.User
	}

	if err = json.Unmarshal(resp.GetJson(), &decode); err != nil {
		return
	}

	if len(decode.User) < 1 {
		err = errors.New("User not found")
		return
	}

	return decode.User[0], nil
}

func CreateUser(user *models.User, config Config) (err error) {
	password, err := bcrypt.GenerateFromPassword([]byte((*user).Password), 8)
	if err != nil {
		return errors.New("Error while hashing password")
	}
	(*user).Password = string(password)

	checkingQuery := `
		query withvar($name: string) {
			user(func: has(user)) @filter(eq(name, $name)) {
				uid
			}
		}
	`

	txn := config.DB().NewTxn()
	defer txn.Discard(config.Context())
	resp, err := txn.QueryWithVars(config.Context(), checkingQuery, map[string]string{"$name": user.UserName})
	if err != nil {
		return err
	}

	if length := gjson.Get(string(resp.GetJson()), "user.#"); length.Num > 0 {
		return errors.New("user already exists")
	}

	ug := &models.UserDgraph{*user, true}
	mu := &api.Mutation{}

	ub, err := json.Marshal(ug)
	if err != nil {
		return err
	}
	mu.SetJson = ub

	_, err = txn.Mutate(config.Context(), mu)
	if err != nil {
		return err
	}
	txn.Commit(config.Context())

	return nil
}

func AddBalance(username string, amount float64, config Config) (balance models.History, err error) {
	q := `
		query withvar($name: string) {
			me(func: has(user)) @filter(eq(name, $name)) {
				user as uid
			}
		}
	`
	mutation := `uid(user) <balance> balance + amount .

				_:history  <from> uid(user)
				_:history <type> "history"
				_:history <historytype> "balance"
				_:history <amount> amount`

	txn := config.DB().NewTxn()
	defer txn.Discard(config.Context())
	mu := &api.Mutation{
		SetNquads: []byte(mutation),
		CommitNow: true,
	}
	mu.Query = q

	if _, err = txn.Mutate(config.Context(), mu); err != nil {
		return
	}
	/*tx := config.DB().Begin()

	var user models.User
	if user, err = GetUser(username, config); err != nil {
		return models.History{}, err
	}
	user.Balance = user.Balance + amount
	if err = tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return models.History{}, err
	}
	balance = models.History{Total: amount, Type: models.BalanceHistory}
	if err := tx.Create(&balance).Error; err != nil {
		tx.Rollback()
		return models.History{}, err
	}
	tx.Commit()*/
	return balance, err
}

func TransferToUser(from string, to string, amount float64, config Config) (transfer models.History, err error) {
	/*tx := config.DB().Begin()

	// Fetch both users from the database
	var fromUser models.User
	var toUser models.User
	if err := tx.Where("user_name = ?", from).First(&fromUser).Error; err != nil {
		return models.History{}, err
	}
	if err := tx.Where("user_name = ?", to).First(&toUser).Error; err != nil {
		return models.History{}, err
	}

	// Check if the amount is positive
	if amount <= 0 {
		return models.History{}, errors.New("amount must be positive")
	}

	// Update both accounts
	fromUser.Balance = fromUser.Balance - amount
	toUser.Balance = toUser.Balance + amount
	if err := tx.Save(&fromUser).Error; err != nil {
		tx.Rollback()
		return models.History{}, err
	}
	if err := tx.Save(&toUser).Error; err != nil {
		tx.Rollback()
		return models.History{}, err
	}
	transfer = models.History{Total: amount, Type: models.TransferHistory, Recipient: toUser.UserName}
	if err := tx.Create(&transfer).Error; err != nil {
		tx.Rollback()
		return models.History{}, err
	}
	tx.Commit()*/
	return transfer, nil
}
