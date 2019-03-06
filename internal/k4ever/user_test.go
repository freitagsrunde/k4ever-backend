package k4ever

import (
	"testing"

	"github.com/freitagsrunde/k4ever-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	conf := NewK4everTest()

	user := models.User{}
	user.UserName = "user"
	user.Password = "password"
	user.DisplayName = "displayname"

	err := CreateUser(&user, conf)

	assert.Equal(t, nil, err)
	assert.Equal(t, uint(1), user.ID)
}
