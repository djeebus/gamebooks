package storage

func NewInMemory() Storage {
	return &InMemory{
		data: make(map[string]interface{}),
	}
}

type InMemory struct {
	data map[string]interface{}
}

func (i *InMemory) Get(key string) interface{} {
	return i.data[key]
}

func (i *InMemory) Set(key string, value interface{}) {
	i.data[key] = value
}

func (i *InMemory) Remove(key string) {
	delete(i.data, key)
}
