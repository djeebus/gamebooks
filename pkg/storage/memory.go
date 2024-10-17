package storage

func NewInMemory() Storage {
	return new(InMemory)
}

type InMemory struct{}
