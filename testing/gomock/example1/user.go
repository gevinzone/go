package example1

import "github.com/gevinzone/go/testing/gomock/example1/person"

type User struct {
	Person person.Person
}

func NewUser(p person.Person) *User {
	return &User{Person: p}
}

func (u *User) GetName(id int64) (string, error) {
	return u.Person.GetName(id)
}
