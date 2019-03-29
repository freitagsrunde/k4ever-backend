package k4ever

import (
	"errors"
	"strings"

	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(params models.DefaultParams, config Config) (users []models.User, err error) {
	tx := config.DB()
	if params.Offset != 0 {
		tx = tx.Offset(params.Offset)
	}
	if params.Limit != 0 {
		tx = tx.Limit(params.Limit)
	}
	if err = tx.Find(&users).Order(params.SortBy + " " + params.Order).Error; err != nil {
		return []models.User{}, err
	}
	return users, err
}

func GetUser(name string, config Config) (user models.User, err error) {
	if err = config.DB().Where("user_name = ?", name).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func CreateUser(user *models.User, config Config) error {
	password, err := bcrypt.GenerateFromPassword([]byte((*user).Password), 8)
	if err != nil {
		return errors.New("Error while hashing password")
	}
	(*user).Password = string(password)
	if err = config.DB().Create(user).Error; err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed:") {
			return errors.New("Username already taken")
		}
		return errors.New("Error while creating user")
	}
	return nil
}

func TransferToUser(from string, to string, amount float64, config Config) error {
	tx := config.DB().Begin()

	// Fetch both users from the database
	var fromUser models.User
	var toUser models.User
	if err := tx.Where("user_name = ?", from).First(&fromUser).Error; err != nil {
		return err
	}
	if err := tx.Where("user_name = ?", to).First(&toUser).Error; err != nil {
		return err
	}

	// Check if the amount is positive
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Update both accounts
	fromUser.Balance = fromUser.Balance - amount
	toUser.Balance = toUser.Balance + amount
	if err := tx.Save(&fromUser).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Save(&toUser).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
