package storage

func Noop() Storage {
	return noop{}
}

type noop struct{}

func (n noop) Get(key string) interface{} {
	return nil
}

func (n noop) Set(key string, value interface{}) {
}
