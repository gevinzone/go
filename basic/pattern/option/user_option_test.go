package option

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
	user := NewUser(target.Username,
		WithAddressOption(target.Address),
		WithNickyNameOption(target.NickyName),
		WithPhoneNumberOption(target.PhoneNumber))
	t.Log(user)
	target.ID = user.ID
	assert.Equal(t, target, *user)
}
