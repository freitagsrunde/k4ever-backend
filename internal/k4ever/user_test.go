package k4ever

import (
	"testing"

	"github.com/freitagsrunde/k4ever-backend/internal/models"
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

	params := DefaultParamsTest()
	params.SortBy = "user_name"
	users, err := GetUsers(params, conf)

	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(users))
}

func TestGetUsers(t *testing.T) {
	conf := NewK4everTest()

	user := UserTest()
	err := CreateUser(&user, conf)

	assert.Equal(t, nil, err)

	params := DefaultParamsTest()
	params.SortBy = "user_name"
	users, err := GetUsers(params, conf)

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(users))
}

func TestGetUserEmpty(t *testing.T) {
	conf := NewK4everTest()

	user, err := GetUser("name", conf)

	assert.Equal(t, models.User{}, user)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "record not found", err.Error())
}

func TestGetUser(t *testing.T) {
	conf := NewK4everTest()

	testUser := UserTest()
	err := CreateUser(&testUser, conf)

	assert.Equal(t, nil, err)

	user, err := GetUser(testUser.UserName, conf)

	assert.Equal(t, nil, err)
	assert.Equal(t, uint(1), user.ID)
}
