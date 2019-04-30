package k4ever

import (
	"os"

	"github.com/jinzhu/gorm"
)

type Config interface {
	Version() string
	GitCommit() string
	GitBranch() string
	BuildTime() string
	LdapHost() string
	LdapBind() string
	LdapPassword() string
	LdapBaseDN() string
	LdapFilterDN() string
	HttpServerHost() string
	DB() *gorm.DB
	SetHttpServerPort(port int)
	HttpServerPort() int
	MigrateDB()
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
