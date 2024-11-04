package storage

type Storage interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Remove(key string) error
	Clear(keyPrefix string) error
}
