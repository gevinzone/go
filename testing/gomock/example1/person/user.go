package person

type User struct {
	Person Person
}

func NewUser(p Person) *User {
	return &User{Person: p}
}

func (u *User) GetPersonInfo(id int64) (*Person, error) {
	p, err := u.Person.Get(id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
