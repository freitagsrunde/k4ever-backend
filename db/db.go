package db

import (
	"github.com/freitagsrunde/k4ever-backend/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Init() *gorm.DB {
	db, err := gorm.Open("sqlite3", "test.db")

	db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}

	return db
}
