package k4ever

import (
	"context"
	"os"

	"github.com/dgraph-io/dgo"
)

type Config interface {
	Version() string
	GitCommit() string
	GitBranch() string
	BuildTime() string
	FilesPath() string
	LdapHost() string
	LdapBind() string
	LdapPassword() string
	LdapBaseDN() string
	LdapFilterDN() string
	HttpServerHost() string
	DB() *dgo.Dgraph
	SetHttpServerPort(port int)
	HttpServerPort() int
	MigrateDB() error
	Context() context.Context
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
