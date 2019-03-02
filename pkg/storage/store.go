package storage

type Store interface {
	Set(string, interface{}) error
	Get(string) (interface{}, error)
	Delete(string) error
}
