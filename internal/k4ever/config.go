package k4ever

import "github.com/jinzhu/gorm"

type Config interface {
	AppVersion() string
	DB() *gorm.DB
	HttpServerPort() int
}
