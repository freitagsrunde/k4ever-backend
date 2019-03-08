package k4ever

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	conf := NewK4everTest()

	user := UserTest()
	err := CreateUser(&user, conf)

	assert.Equal(t, nil, err)
	assert.Equal(t, uint(1), user.ID)
}

func TestGetUsersEmpty(t *testing.T) {
	conf := NewK4everTest()

	users, err := GetUsers("user_name", "asc", conf)

	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(users))
}

func TestGetUsers(t *testing.T) {
	conf := NewK4everTest()

	user := UserTest()
	err := CreateUser(&user, conf)

	assert.Equal(t, nil, err)

	users, err := GetUsers("user_name", "asc", conf)

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(users))
}
