package k4ever

import "github.com/jinzhu/gorm"

type Config interface {
	Version() string
	GitCommit() string
	GitBranch() string
	BuildTime() string
	DB() *gorm.DB
	HttpServerPort() int
	MigrateDB()
}
