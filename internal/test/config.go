package test

import "github.com/jinzhu/gorm"

type Config struct {
	db *gorm.DB
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) AppVersion() string {
	return "1.0"
}

func (c *Config) DB() *gorm.DB {
	return nil
}

func (c *Config) HttpServerPort() int {
	return 8080
}

func (c *Config) connectToDatabase() error {
	return nil
}

func (c *Config) MigrateDB() {

}
