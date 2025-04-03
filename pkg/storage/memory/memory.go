package memory

import (
	"gamebooks/pkg/storage"
	"strings"
)

func New() storage.Storage {
	return &InMemory{
		data: make(map[string]interface{}),
	}
}

type InMemory struct {
	data map[string]interface{}
}

func (i *InMemory) Get(key string) (interface{}, error) {
	return i.data[key], nil
}

func (i *InMemory) Set(key string, value interface{}) error {
	i.data[key] = value
	return nil
}

func (i *InMemory) Remove(key string) error {
	delete(i.data, key)
	return nil
}

func (i *InMemory) Clear(keyPrefix string) error {
	for key := range i.data {
		if strings.HasPrefix(key, keyPrefix) {
			delete(i.data, key)
		}
	}

	return nil
}
