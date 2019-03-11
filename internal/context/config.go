package context

import (
	"fmt"

	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
)

type Config struct {
	version        string
	db             *gorm.DB
	httpServerPort int
	gitCommit      string
	gitBranch      string
	buildTime      string
}

var GitCommit = "undefined"
var GitBranch = "undefined"
var BuildTime = "No Time provided"
var version string

func NewConfig() *Config {
	c := &Config{}
	// c.AppVersion = viper.GetString("version")
	fmt.Println(version)
	c.version = version
	c.httpServerPort = viper.GetInt("port")
	c.gitCommit = GitCommit
	c.gitBranch = GitBranch
	c.buildTime = BuildTime

	return c
}

func (c *Config) Version() string {
	return c.version
}

func (c *Config) GitCommit() string {
	return c.gitCommit
}

func (c *Config) GitBranch() string {
	return c.gitBranch
}

func (c *Config) BuildTime() string {
	return c.buildTime
}

func (c *Config) DB() *gorm.DB {
	if c.db == nil {
		c.connectToDatabase()
	}
	return c.db
}

func (c *Config) HttpServerPort() int {
	return c.httpServerPort
}

func (c *Config) connectToDatabase() error {
	db, err := gorm.Open("postgres", "host=postgres port=5432 user=postgres dbname=postgres password=postgres sslmode=disable")
	c.db = db

	return err
}

func (c *Config) MigrateDB() {
	db := c.DB()

	db.AutoMigrate(
		&models.Product{},
		&models.User{},
		&models.Permission{},
		&models.Purchase{},
		&models.PurchaseItem{},
	)
}
