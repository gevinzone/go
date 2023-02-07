package option

import "time"

type User struct {
	ID       int64
	Username string

	NickyName   string
	Address     string
	PhoneNumber string
}

type UserOption func(*User)

func NewUser(username string, options ...UserOption) *User {
	user := &User{
		ID:       time.Now().UnixMilli(),
		Username: username,
	}

	for _, opt := range options {
		opt(user)
	}

	return user
}

func WithNickyNameOption(nickyName string) UserOption {
	return func(user *User) {
		user.NickyName = nickyName
	}
}

func WithPhoneNumberOption(phoneNumber string) UserOption {
	return func(user *User) {
		user.PhoneNumber = phoneNumber
	}
}

func WithAddressOption(address string) UserOption {
	return func(user *User) {
		user.Address = address
	}
}