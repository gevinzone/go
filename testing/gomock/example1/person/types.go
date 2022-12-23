package person

//go:generate mockgen -source=./types.go -destination=./person_mock.go -package=person
type Person interface {
	Get(id int64) (Person, error)
	GetByName(name string) (Person, error)
}
