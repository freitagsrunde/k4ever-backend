package context

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"google.golang.org/grpc"

	"github.com/freitagsrunde/k4ever-backend/internal/k4ever"
	"github.com/spf13/viper"
)

type Config struct {
	version        string
	context        context.Context
	db             *dgo.Dgraph
	dbHost         string
	dbPort         int
	dbName         string
	dbPass         string
	dbSSLMode      string
	filesPath      string
	ldapHost       string
	ldapBind       string
	ldapPassword   string
	ldapBaseDN     string
	ldapFilterDN   string
	httpServerHost string
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
	c.context = context.Background()
	// c.AppVersion = viper.GetString("version")
	fmt.Println(version)
	c.version = version
	c.httpServerPort = viper.GetInt("port")
	c.gitCommit = GitCommit
	c.gitBranch = GitBranch
	c.buildTime = BuildTime
	c.filesPath = k4ever.GetEnv("K4EVER_FILESPATH", ".")
	c.ldapHost = k4ever.GetEnv("K4EVER_LDAPHOST", "localhost")
	c.ldapBind = k4ever.GetEnv("K4EVER_LDAPBIND", "admin")
	c.ldapPassword = k4ever.GetEnv("K4EVER_LDAPPASSWORD", "admin")
	c.ldapBaseDN = k4ever.GetEnv("K4EVER_LDAPBASEDN", "CN=Users,DC=example,DC=com")
	c.ldapFilterDN = k4ever.GetEnv("K4EVER_LDAPFILTERDN", "(&(objectClass=person)(uid={username}))")
	c.httpServerHost = k4ever.GetEnv("K4EVER_DOMAIN", "localhost")

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

func (c *Config) FilesPath() string {
	return c.filesPath
}

func (c *Config) BuildTime() string {
	return c.buildTime
}

func (c *Config) DB() *dgo.Dgraph {
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

func (c *Config) HttpServerHost() string {
	return c.httpServerHost
}

func (c *Config) SetHttpServerPort(port int) {
	c.httpServerPort = port
}

func (c *Config) HttpServerPort() int {
	return c.httpServerPort
}

func (c *Config) Context() context.Context {
	return c.context
}

func (c *Config) connectToDatabase() error {
	host := k4ever.GetEnv("K4EVER_DBHOST", "localhost")
	portS := k4ever.GetEnv("K4EVER_DBPORT", "9080")
	_, err := strconv.Atoi(portS)
	if err != nil {
		return err
	}

	d, err := grpc.Dial(host+":"+portS, grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.db = dgo.NewDgraphClient(api.NewDgraphClient(d))
	/*user := k4ever.GetEnv("K4EVER_DBUSER", "postgres")
	dbname := k4ever.GetEnv("K4EVER_DBNAME", "postgres")
	password := k4ever.GetEnv("K4EVER_DBPASS", "postgres")
	sslmode := k4ever.GetEnv("K4EVER_DBSSL", "disable")
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode))
	c.db = db
	*/

	return err
}

func (c *Config) MigrateDB() error {
	db := c.DB()
	op := &api.Operation{}
	op.Schema = `
		name: string @index(exact) @upsert .
	`
	err := db.Alter(c.context, op)
	if err != nil {
		return err
	}
	return nil
}
