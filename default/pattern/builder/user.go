package builder

import "time"

type User struct {
	ID       int64
	Username string

	NickyName   string
	Address     string
	PhoneNumber string
}

type UserBuilderV1 struct {
	username    string
	nickyName   string
	address     string
	phoneNumber string
}

func (ub *UserBuilderV1) Username(username string) *UserBuilderV1 {
	ub.username = username
	return ub
}

func (ub *UserBuilderV1) NickyName(nickyName string) *UserBuilderV1 {
	ub.nickyName = nickyName
	return ub
}

func (ub *UserBuilderV1) Address(address string) *UserBuilderV1 {
	ub.address = address
	return ub
}

func (ub *UserBuilderV1) PhoneNumber(phoneNumber string) *UserBuilderV1 {
	ub.phoneNumber = phoneNumber
	return ub
}

func (ub *UserBuilderV1) Build() *User {
	return &User{
		ID:          time.Now().UnixMilli(),
		Username:    ub.username,
		NickyName:   ub.nickyName,
		Address:     ub.address,
		PhoneNumber: ub.phoneNumber,
	}
}

type UserBuilderV2 struct {
	user User
}

func (ub *UserBuilderV2) Username(username string) *UserBuilderV2 {
	ub.user.Username = username
	return ub
}

func (ub *UserBuilderV2) NickyName(nickyName string) *UserBuilderV2 {
	ub.user.NickyName = nickyName
	return ub
}

func (ub *UserBuilderV2) Address(address string) *UserBuilderV2 {
	ub.user.Address = address
	return ub
}

func (ub *UserBuilderV2) PhoneNumber(phoneNumber string) *UserBuilderV2 {
	ub.user.PhoneNumber = phoneNumber
	return ub
}

func (ub *UserBuilderV2) Build() *User {
	return &ub.user
}
