package context

import (
	"fmt"
	"os"
	"strconv"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

type Config struct {
	version        string
	db             *gorm.DB
	dbHost         string
	dbPort         int
	dbName         string
	dbPass         string
	dbSSLMode      string
	ldapHost       string
	ldapBind       string
	ldapPassword   string
	ldapBaseDN     string
	ldapFilterDN   string
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
		if err := c.connectToDatabase(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	return c.db
}

func (c *Config) LdapHost() string {
	return c.ldapHost
}

func (c *Config) LdapBind() string {
	return c.ldapBind
}

func (c *Config) LdapPassword() string {
	return c.ldapPassword
}

func (c *Config) LdapBaseDN() string {
	return c.ldapBaseDN
}

func (c *Config) LdapFilterDN() string {
	return c.ldapFilterDN
}

func (c *Config) SetHttpServerPort(port int) {
	c.httpServerPort = port
}

func (c *Config) HttpServerPort() int {
	return c.httpServerPort
}

func (c *Config) connectToDatabase() error {
	host := k4ever.GetEnv("K4EVER_DBHOST", "postgres")
	portS := k4ever.GetEnv("K4EVER_DBPORT", "5432")
	port, err := strconv.Atoi(portS)
	if err != nil {
		return err
	}
	user := k4ever.GetEnv("K4EVER_DBUSER", "postgres")
	dbname := k4ever.GetEnv("K4EVER_DBNAME", "postgres")
	password := k4ever.GetEnv("K4EVER_DBPASS", "postgres")
	sslmode := k4ever.GetEnv("K4EVER_DBSSL", "disable")
	c.ldapHost = k4ever.GetEnv("K4EVER_LDAPHOST", "localhost")
	c.ldapBind = k4ever.GetEnv("K4EVER_LDAPBIND", "admin")
	c.ldapPassword = k4ever.GetEnv("K4EVER_LDAPPASSWORD", "admin")
	c.ldapBaseDN = k4ever.GetEnv("K4EVER_LDAPBASEDN", "CN=Users,DC=example,DC=com")
	c.ldapFilterDN = k4ever.GetEnv("K4EVER_LDAPFILTERDN", "(&(objectClass=person)(uid={username}))")
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode))
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
