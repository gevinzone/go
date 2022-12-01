package builder

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	target := User{
		Username:    "Tom",
		NickyName:   "Tommy",
		Address:     "somewhere",
		PhoneNumber: "18867654328",
	}
	userBuilder := UserBuilderV1{}
	user := userBuilder.Username(target.Username).
		NickyName(target.NickyName).
		Address(target.Address).
		PhoneNumber(target.PhoneNumber).
		Build()
	target.ID = user.ID
	t.Log(user)
	assert.Equal(t, target, *user)
}

func TestNewUserV2(t *testing.T) {
	target := User{
		Username:    "Tom",
		NickyName:   "Tommy",
		Address:     "somewhere",
		PhoneNumber: "18867654328",
	}

	userBuilder := UserBuilderV2{}
	user := userBuilder.Username(target.Username).
		NickyName(target.NickyName).
		Address(target.Address).
		PhoneNumber(target.PhoneNumber).
		Build()
	target.ID = user.ID
	t.Log(user)
	assert.Equal(t, target, *user)
}
