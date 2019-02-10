package k4ever

import (
	"errors"
	"strings"

	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

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
