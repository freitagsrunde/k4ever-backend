package k4ever

import (
	"github.com/freitagsrunde/k4ever-backend/internal/test"
)

func NewK4everTest() (conf Config) {
	conf = test.NewConfig()
	conf.MigrateDB()

	return conf
}
