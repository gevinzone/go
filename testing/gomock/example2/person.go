package example2

//go:generate mockgen -source=./person.go -destination=./mocks/person_mock.go -package=person_mock
type Person interface {
	GetName(id int64) (string, error)
}
